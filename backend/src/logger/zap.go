package logger

import (
	"bytes"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/peername"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"runtime"
	"time"
)

type BackendLogger struct {
	*zap.Logger
	ServiceName string `json:"serviceName"`
}

var Logger BackendLogger

func (l BackendLogger) Error(msg string, args ...any) {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)

	data := fmt.Sprintf(msg, args...)
	dataMsg := fmt.Sprintf("%s\nError stack:\n%s", data, string(buf[:n]))

	l.Logger.Error(dataMsg)
	_ = LogMsg(true, fmt.Sprintf("[%s-%s] %s", l.ServiceName, peername.PeerName, dataMsg))
}

func (l BackendLogger) Tag(name string, args ...any) {
	data := fmt.Sprintln(args...)
	dataMsg := fmt.Sprintf("TAG %s - %s\n%s", name, time.Now().Format("15:04:05.000"), data) // data包含回车
	dataMsg = dataMsg[0 : len(dataMsg)-1]                                                    // 删除回车

	l.Logger.Info(dataMsg)
	_ = LogMsg(true, fmt.Sprintf("[%s-%s] %s", l.ServiceName, peername.PeerName, dataMsg))
}

func (l BackendLogger) Info(msg string, args ...any) {
	data := fmt.Sprintf(msg, args...)
	l.Logger.Info(data)
}

func (l BackendLogger) WXInfo(msg string, args ...any) {
	data := fmt.Sprintf(msg, args...)
	l.Logger.Info(data)
	_ = LogMsg(false, fmt.Sprintf("[%s-%s] %s", l.ServiceName, peername.PeerName, data))
}

func InitLogger(serviceName string) (err errors.WTError) {
	if config.BackendConfig.GetMode() == config.RunModeDevelop {
		c := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
			Development: true,
			Encoding:    "console",
			EncoderConfig: zapcore.EncoderConfig{
				// Keys can be anything except the empty string.
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}

		l, err := c.Build(zap.AddCaller())
		if err != nil {
			return errors.WarpQuick(err)
		}

		Logger = BackendLogger{
			Logger:      l,
			ServiceName: serviceName,
		}
	} else {
		c := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding: "console",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.EpochTimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}

		l, err := c.Build(zap.AddCaller())
		if err != nil {
			return errors.WarpQuick(err)
		}

		Logger = BackendLogger{
			Logger:      l,
			ServiceName: serviceName,
		}
	}

	return nil
}

type Msg struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
}

type Text struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list"`
}

func LogMsg(atall bool, text string, args ...any) errors.WTError {
	return WxRobotSendNotRecord(config.BackendConfig.WXRobot.Log, fmt.Sprintf(text, args...), atall)
}

func WxRobotSendNotRecord(webhook string, text string, atAll bool) errors.WTError {
	if len(webhook) == 0 {
		return nil
	}

	t := Text{
		Content: text,
	}

	if atAll {
		t.MentionedList = append(t.MentionedList, "@all")
	}

	data := Msg{
		MsgType: "text",
		Text:    t,
	}
	dataByte, jsonErr := utils.JsonMarshal(data)
	if jsonErr != nil {
		return jsonErr
	}

	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewBuffer(dataByte))
	if err != nil {
		return errors.WarpQuick(err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return errors.WarpQuick(err)
	}

	if resp.StatusCode != 200 {
		return errors.Errorf("get bad status code")
	}

	return nil
}
