package respmsg

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/record"
	errors "github.com/wuntsong-org/wterrors"
	"runtime"
)

func GetRespSuccess(ctx context.Context) types.Resp {
	recordData := record.GetRecord(ctx)
	recordData.Msg = "成功"

	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	recordData.Stack = string(buf[:n])

	return types.Resp{
		RequestsID: recordData.RequestsID,
		Code:       SuccessCode.Code,
		NumCode:    SuccessCode.NumCode,
		SubCode:    Success.Code,
		SubNumCode: Success.NumCode,
		Msg:        "成功",
	}
}

func GetRespSuccessWithDebug(ctx context.Context, m ...any) types.Resp {
	recordData := record.GetRecord(ctx)
	var debugMsg string
	if len(m) == 0 {
		debugMsg = "成功"
	} else {
		var mm string
		m1, ok := m[0].(string)
		if ok {
			mm = fmt.Sprintf(m1, m[1:]...)
		} else {
			mm = fmt.Sprintln(m...)
			mm = mm[:len(mm)-1] // 删除回车
		}
		debugMsg = fmt.Sprintf("成功：%s", mm)
	}
	recordData.Msg = debugMsg

	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	recordData.Stack = string(buf[:n])

	if config.BackendConfig.GetMode() == config.RunModeDevelop {
		return types.Resp{
			RequestsID: recordData.RequestsID,
			Code:       SuccessCode.Code,
			NumCode:    SuccessCode.NumCode,
			SubCode:    Success.Code,
			SubNumCode: Success.NumCode,
			Msg:        "成功",
			DebugMsg:   debugMsg,
		}
	}

	return types.Resp{
		RequestsID: recordData.RequestsID,
		Code:       SuccessCode.Code,
		NumCode:    SuccessCode.NumCode,
		SubCode:    Success.Code,
		SubNumCode: Success.NumCode,
		Msg:        "成功",
	}
}

func GetRespByMsg(ctx context.Context, code SubCode, m string, args ...interface{}) types.Resp {
	msg := fmt.Sprintf(m, args...)

	recordData := record.GetRecord(ctx)
	recordData.Msg = msg

	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	recordData.Stack = string(buf[:n])

	return types.Resp{
		RequestsID: recordData.RequestsID,
		Code:       LogicCode.Code,
		NumCode:    LogicCode.NumCode,
		SubCode:    code.Code,
		SubNumCode: code.NumCode,
		Msg:        msg,
	}
}

func GetRespByMsgWithDebug(ctx context.Context, code SubCode, debugMsg string, m string, args ...interface{}) types.Resp {
	msg := fmt.Sprintf(m, args...)

	recordData := record.GetRecord(ctx)
	recordData.Msg = fmt.Sprintf("%s（%s）", debugMsg, msg)

	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	recordData.Stack = string(buf[:n])

	if config.BackendConfig.GetMode() == config.RunModeDevelop {
		return types.Resp{
			RequestsID: recordData.RequestsID,
			Code:       LogicCode.Code,
			NumCode:    LogicCode.NumCode,
			SubCode:    code.Code,
			SubNumCode: code.NumCode,
			Msg:        msg,
			DebugMsg:   debugMsg,
		}
	}

	return types.Resp{
		RequestsID: recordData.RequestsID,
		Code:       LogicCode.Code,
		NumCode:    LogicCode.NumCode,
		SubCode:    code.Code,
		SubNumCode: code.NumCode,
		Msg:        msg,
	}
}

func GetRespByError(ctx context.Context, code SubCode, err errors.WTError, m ...any) types.Resp {
	recordData := record.GetRecord(ctx)

	if err == nil {
		err = errors.Errorf("<empty> unknown error")
	}

	msg := ""
	if len(m) != 0 {
		m1, ok := m[0].(string)
		if ok {
			msg = fmt.Sprintf(m1, m[1:]...)
		} else {
			msg = fmt.Sprintln(m...)
			msg = msg[:len(msg)-1] // 删除回车
		}
	}

	recordData.Msg = msg
	recordData.Err = err

	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	recordData.Stack = string(buf[:n])

	if config.BackendConfig.GetMode() == config.RunModeDevelop {
		return types.Resp{
			RequestsID: recordData.RequestsID,
			Code:       LogicCode.Code,
			NumCode:    LogicCode.NumCode,
			SubCode:    code.Code,
			SubNumCode: code.NumCode,
			Msg:        msg,
			DebugMsg:   err.Message(),
		}
	}

	return types.Resp{
		RequestsID: recordData.RequestsID,
		Code:       LogicCode.Code,
		NumCode:    LogicCode.NumCode,
		SubCode:    code.Code,
		SubNumCode: code.NumCode,
		DebugMsg:   err.Message(),
	}
}

func GetRespByErrorWithDebug(ctx context.Context, code SubCode, err errors.WTError, debugMsg string, m ...any) types.Resp {
	recordData := record.GetRecord(ctx)

	if err == nil {
		err = errors.Errorf("<empty> unknown error")
	}

	msg := ""
	if len(m) != 0 {
		m1, ok := m[0].(string)
		if ok {
			msg = fmt.Sprintf(m1, m[1:]...)
		} else {
			msg = fmt.Sprintln(m...)
			msg = msg[:len(msg)-1] // 删除回车
		}
	}

	recordData.Msg = fmt.Sprintf("%s（%s）", debugMsg, msg)
	recordData.Err = err

	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	recordData.Stack = string(buf[:n])

	if config.BackendConfig.GetMode() == config.RunModeDevelop {
		return types.Resp{
			RequestsID: recordData.RequestsID,
			Code:       LogicCode.Code,
			NumCode:    LogicCode.NumCode,
			SubCode:    code.Code,
			SubNumCode: code.NumCode,
			Msg:        msg,
			DebugMsg:   fmt.Sprintf("%s（%s）", debugMsg, err.Error()),
		}
	}

	return types.Resp{
		RequestsID: recordData.RequestsID,
		Code:       LogicCode.Code,
		NumCode:    LogicCode.NumCode,
		SubCode:    code.Code,
		SubNumCode: code.NumCode,
		DebugMsg:   err.Message(),
	}
}

func GetRespByMsgWithCode(ctx context.Context, code Code, subCode SubCode, m string, args ...interface{}) types.Resp {
	msg := fmt.Sprintf(m, args...)

	recordData := record.GetRecord(ctx)
	recordData.Msg = msg

	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	recordData.Stack = string(buf[:n])

	return types.Resp{
		RequestsID: recordData.RequestsID,
		Code:       code.Code,
		NumCode:    code.NumCode,
		SubCode:    subCode.Code,
		SubNumCode: subCode.NumCode,
		Msg:        msg,
	}
}

func GetRespByErrorWithCode(ctx context.Context, code Code, subCode SubCode, err errors.WTError, m ...any) types.Resp {
	recordData := record.GetRecord(ctx)

	if err == nil {
		err = errors.Errorf("<empty> unknown error")
	}

	msg := ""
	if len(m) != 0 {
		m1, ok := m[0].(string)
		if ok {
			msg = fmt.Sprintf(m1, m[1:]...)
		} else {
			msg = fmt.Sprintln(m...)
			msg = msg[:len(msg)-1] // 删除回车
		}
	}

	recordData.Msg = msg
	recordData.Err = err

	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	recordData.Stack = string(buf[:n])

	if config.BackendConfig.GetMode() == config.RunModeDevelop {
		return types.Resp{
			RequestsID: recordData.RequestsID,
			Code:       code.Code,
			NumCode:    code.NumCode,
			SubCode:    subCode.Code,
			SubNumCode: subCode.NumCode,
			Msg:        msg,
			DebugMsg:   err.Message(),
		}
	}

	return types.Resp{
		RequestsID: recordData.RequestsID,
		Code:       code.Code,
		NumCode:    code.NumCode,
		SubCode:    subCode.Code,
		SubNumCode: subCode.NumCode,
		DebugMsg:   err.Message(),
	}
}
