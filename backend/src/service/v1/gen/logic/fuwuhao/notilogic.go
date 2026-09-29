package fuwuhao

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/fuwuhao"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/fastwego/offiaccount/type/type_event"
	"github.com/fastwego/offiaccount/type/type_message"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type NotiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotiLogic {
	return &NotiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotiLogic) Noti(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	log.Println(string(body))

	message, err := fuwuhao.OffiAccount.Server.ParseXML(body)
	if err != nil {
		log.Println(err)
	}

	var output interface{}
	switch msg := message.(type) {
	case type_message.MessageText: // 文本 消息
		if msg.MsgType == type_message.MsgTypeText {
			switch msg.Content {
			case fmt.Sprintf("绑定%s账号", config.BackendConfig.User.ReadableName):
				res := fuwuhao.Bind(r.Context(), msg.FromUserName)
				content := ""
				if res {
					content = "绑定成功"
				} else {
					content = fmt.Sprintf("绑定失败，请检查%s账号是否已经绑定该微信号。", config.BackendConfig.User.ReadableName)
				}

				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA(content),
				}
			case "Ping", "ping":
				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA("Pong"),
				}
			case "客服", "在线客服", "人工客服":
				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA(fmt.Sprintf(`欢迎咨询人工客服，请点击客服链接：<a  href="%s">客服</a>`, config.BackendConfig.FuWuHao.Kefu.HuanChuang)),
				}
			case "创思域变客服":
				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA(fmt.Sprintf(`欢迎咨询人工客服，请点击客服链接：<a  href="%s">创思域变客服</a>`, config.BackendConfig.FuWuHao.Kefu.Vxwk)),
				}
			case "Ping业务", "ping业务":
				go func() {
					err = fuwuhao.SendPing(msg.FromUserName)
					if err != nil {
						logger.Logger.Error(fmt.Sprintf("ping error: %s", err.Error()))
					}
				}()

				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA("正在为你处理ping业务"),
				}
			default:
				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeTransferCustomerService,
					},
				}
			}
		} else {
			output = type_message.ReplyMessageText{
				ReplyMessage: type_message.ReplyMessage{
					ToUserName:   type_message.CDATA(msg.FromUserName),
					FromUserName: type_message.CDATA(msg.ToUserName),
					CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
					MsgType:      type_message.ReplyMsgTypeTransferCustomerService,
				},
			}
		}
	case type_message.MessageImage:
		output = type_message.ReplyMessageText{
			ReplyMessage: type_message.ReplyMessage{
				ToUserName:   type_message.CDATA(msg.FromUserName),
				FromUserName: type_message.CDATA(msg.ToUserName),
				CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
				MsgType:      type_message.ReplyMsgTypeTransferCustomerService,
			},
		}
	case type_message.MessageVideo:
		output = type_message.ReplyMessageText{
			ReplyMessage: type_message.ReplyMessage{
				ToUserName:   type_message.CDATA(msg.FromUserName),
				FromUserName: type_message.CDATA(msg.ToUserName),
				CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
				MsgType:      type_message.ReplyMsgTypeTransferCustomerService,
			},
		}
	case type_message.MessageVoice:
		output = type_message.ReplyMessageText{
			ReplyMessage: type_message.ReplyMessage{
				ToUserName:   type_message.CDATA(msg.FromUserName),
				FromUserName: type_message.CDATA(msg.ToUserName),
				CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
				MsgType:      type_message.ReplyMsgTypeTransferCustomerService,
			},
		}
	case type_event.EventMenuClick:
		if msg.Event.Event == type_event.EventTypeMenuClick {
			switch msg.EventKey {
			case fuwuhao.BindAuthKey:
				res := fuwuhao.Bind(r.Context(), msg.FromUserName)
				content := ""
				if res {
					content = "绑定成功"
				} else {
					content = fmt.Sprintf("绑定失败，请检查%s账号是否已经绑定该微信号。", config.BackendConfig.User.ReadableName)
				}

				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA(content),
				}
			case fuwuhao.AboutUsContactKey:
				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA(config.BackendConfig.FuWuHao.Menu.AboutUsContact),
				}
			case fuwuhao.AboutUsKefuKey:
				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA(config.BackendConfig.FuWuHao.Menu.AboutUsKefu),
				}
			}
		} else {
			logger.Logger.Error("bad event type")
		}
	case type_event.EventSubscribe:
		if msg.Event.Event == type_event.EventTypeSubscribe {
			res := fuwuhao.Bind(r.Context(), msg.FromUserName)
			if res {
				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA(fmt.Sprintf("欢迎关注%s公众号，已经为了自动绑定%s账号。", config.BackendConfig.User.ReadableName, config.BackendConfig.User.ReadableName)),
				}
			} else {
				output = type_message.ReplyMessageText{
					ReplyMessage: type_message.ReplyMessage{
						ToUserName:   type_message.CDATA(msg.FromUserName),
						FromUserName: type_message.CDATA(msg.ToUserName),
						CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
						MsgType:      type_message.ReplyMsgTypeText,
					},
					Content: type_message.CDATA(fmt.Sprintf("欢迎关注%s公众号，注册%s账号尽情体验吧！", config.BackendConfig.User.ReadableName, config.BackendConfig.User.ReadableName)),
				}
			}
		} else {
			output = type_message.ReplyMessageText{
				ReplyMessage: type_message.ReplyMessage{
					ToUserName:   type_message.CDATA(msg.FromUserName),
					FromUserName: type_message.CDATA(msg.ToUserName),
					CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
					MsgType:      type_message.ReplyMsgTypeTransferCustomerService,
				},
			}
		}
	}

	if output != nil {
		_ = fuwuhao.OffiAccount.Server.Response(w, r, output)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
