package ws

import (
	"context"
	websocket2 "gitee.com/wuntsong-auth/backend/src/global/websocket"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/ws/webwstype"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/signalexit"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/gorilla/websocket"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"
	"sync"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type WSGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWSGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WSGetInfoLogic {
	return &WSGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WSGetInfoLogic) WSGetInfo(w http.ResponseWriter, r *http.Request) {
	web, ok := r.Context().Value("X-Src-Website").(warp.Website)
	if !ok {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusInternalServerError, &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.UnknownError, respmsg.BadContextError.New("X-Src-Website")),
		})
		return
	}

	if !websocket.IsWebSocketUpgrade(r) {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusBadRequest, &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		})
		return
	}

	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// IPCheck之后不需要检测跨域
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusBadRequest, &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.SystemError, errors.WarpQuick(err), "websocket建立链接失败"),
		})
		return
	}

	Connect(conn, web)
}

func Connect(conn *websocket.Conn, web warp.Website) {
	ch := make(chan websocket2.WSWebsiteMessage, 20)

	var closerFuncMutex = new(sync.Mutex)
	closerFunc := make(map[string]func(), 10)

	closer := func() {
		closerFuncMutex.Lock()
		defer closerFuncMutex.Unlock()

		for _, f := range closerFunc {
			if f != nil {
				func() {
					defer func() {
						recover()
					}()

					f()
				}()
			}
		}

		closerFunc = make(map[string]func(), 10)
	}

	var deadlineMutex = new(sync.Mutex)
	deadline := time.Now().Add(time.Minute * 5)

	deadlineNow := func() {
		deadlineMutex.Lock()
		defer deadlineMutex.Unlock()
		deadline = time.Now()
	}

	go func() {
		for {
			time.Sleep(30 * time.Second)

			if func() bool {
				defer utils.Recover(logger.Logger, nil, "ws check deadline")

				deadlineMutex.Lock()
				defer deadlineMutex.Unlock()

				if time.Now().After(deadline) {
					func() {
						websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
							Code: webwstype.WebsiteClose,
						})
						// 不需要 _ = conn.Close()
						closer()
					}()
					return true
				}

				return false
			}() {
				return // 退出go
			}
		}
	}()

	go func() {
		defer signalexit.AddExitFuncAsDefer(func() {
			defer utils.Recover(logger.Logger, nil, "ws reader")
			websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
				Code: webwstype.WebsiteClose,
			})
			// 不需要 _ = conn.Close()
			closer()
			deadlineNow()
		})()

		for {
			var msg websocket2.WSWebsiteMessage
			err := conn.ReadJSON(&msg)
			if utils.IsNetClose(err) || errors.Is(err, io.ErrUnexpectedEOF) || websocket.IsUnexpectedCloseError(err) {
				return
			} else if err != nil {
				return
			}

			switch msg.Code {
			default:
				websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
					Code: webwstype.WebsiteBadCode,
				})
			case webwstype.WebsitePing:
				websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
					Code: webwstype.WebsitePong,
				})

				func() {
					defer func() {
						recover()
					}()

					deadlineMutex.Lock()
					defer deadlineMutex.Unlock()

					deadline = time.Now().Add(time.Minute * 5)
				}()
			case webwstype.WebsiteGetUserInfo:
				func() {
					type RespData struct {
						Success bool               `json:"success"`
						Error   bool               `json:"error"`
						Token   string             `json:"token"`
						Msg     string             `json:"msg"`
						User    types.UserEasy     `json:"user"`
						Data    types.UserData     `json:"data"`
						Info    types.UserInfoEsay `json:"info"`
					}

					data := struct {
						Token string `json:"token"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
							Code: webwstype.WebsiteBadToken,
							Data: websocket2.WriteJson(RespData{
								Token:   data.Token,
								Error:   false,
								Msg:     err.Error(),
								Success: false,
							}),
						})
						return
					}

					loginData, err := jwt.ParserLoginToken(context.Background(), data.Token)
					if err != nil || loginData.WebID != web.ID {
						websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
							Code: webwstype.WebsiteBadToken,
							Data: websocket2.WriteJson(RespData{
								Token:   data.Token,
								Error:   false,
								Msg:     err.Error(),
								Success: false,
							}),
						})
					}

					userModel := db.NewUserModel(mysql.MySQLConn)

					user, mysqlErr := userModel.FindOneByUidWithoutDelete(context.Background(), loginData.UserID)
					if mysqlErr != nil {
						if errors.Is(mysqlErr, db.ErrNotFound) {
							websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
								Code: webwstype.WebsiteBadToken,
								Data: websocket2.WriteJson(RespData{
									Token:   data.Token,
									Error:   false,
									Msg:     "user not found",
									Success: false,
								}),
							})
							return
						}

						websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
							Code: webwstype.WebsiteBadToken,
							Data: websocket2.WriteJson(RespData{
								Token: data.Token,
								Msg: websocket2.WriteJson(RespData{
									Token:   data.Token,
									Error:   true,
									Msg:     err.Error(),
									Success: false,
								}),
							}),
						})
						return
					}

					if db.IsBanned(user) {
						websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
							Code: webwstype.WebsiteBadToken,
							Data: websocket2.WriteJson(RespData{
								Token:   data.Token,
								Error:   false,
								Msg:     "user not found",
								Success: false,
							}),
						})
						return
					}

					userData, mysqlErr := utils2.GetUserInfoEasy(context.Background(), user)
					if mysqlErr != nil {
						websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
							Code: webwstype.WebsiteBadToken,
							Data: websocket2.WriteJson(RespData{
								Token: data.Token,
								Msg: websocket2.WriteJson(RespData{
									Token:   data.Token,
									Error:   true,
									Msg:     mysqlErr.Error(),
									Success: false,
								}),
							}),
						})
						return
					}

					websocket2.WriteMessage(ch, websocket2.WSWebsiteMessage{
						Code: webwstype.WebsiteUserInfo,
						Data: websocket2.WriteJson(RespData{
							Token:   data.Token,
							Msg:     "success",
							Success: true,
							Error:   false,
							User:    userData.User,
							Info:    userData.InfoEasy,
							Data:    userData.Data,
						}),
					})
					return
				}()
			}

		}
	}()

	go func() {
		defer signalexit.AddExitFuncAsDefer(func() {
			defer utils.Recover(logger.Logger, nil, "ws writer")
			_ = conn.Close()
			closer()
			deadlineNow()
			close(ch)
		})()

		defer utils.Recover(logger.Logger, nil, "ws writer")

		for {
			msg := <-ch

			switch msg.Code {
			case webwstype.WebsiteClose:
				return
			}

			err := conn.WriteJSON(msg)
			if utils.IsNetClose(err) || websocket.IsUnexpectedCloseError(err) {
				return
			} else if err != nil {
				return
			}
		}
	}()
}
