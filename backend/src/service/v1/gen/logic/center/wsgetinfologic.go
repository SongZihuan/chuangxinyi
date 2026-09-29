package center

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	websocket2 "gitee.com/wuntsong-auth/backend/src/global/websocket"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/center/userwstype"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	"gitee.com/wuntsong-auth/backend/src/signalexit"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/gorilla/websocket"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
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

func DeleteCh[T string | int64](m map[T][]chan websocket2.WSMessage, k T, index int, mutex *sync.Mutex) int {
	defer func() {
		recover()
	}()

	mutex.Lock()
	defer mutex.Unlock()

	lst := m[k]

	lst[index] = nil
	lst = utils.RemoveTrailingNil(lst, func(c chan websocket2.WSMessage) bool { return c == nil })
	if lst == nil {
		delete(m, k)
	} else {
		m[k] = lst
	}

	return -1
}

func AppendCh[T string | int64](m map[T][]chan websocket2.WSMessage, k T, c chan websocket2.WSMessage, mutex *sync.Mutex) int {
	defer func() {
		recover()
	}()

	mutex.Lock()
	defer mutex.Unlock()

	for i, lc := range m[k] {
		if lc == nil {
			m[k][i] = c
			return i
		}
	}

	m[k] = append(m[k], c)
	return len(m[k]) - 1
}

func DeleteChFromList(m []chan websocket2.WSMessage, index int, mutex *sync.Mutex) int {
	defer func() {
		recover()
	}()

	mutex.Lock()
	defer mutex.Unlock()

	m[index] = nil
	return -1
}

func AppendChToList(m *[]chan websocket2.WSMessage, c chan websocket2.WSMessage, mutex *sync.Mutex) int {
	defer func() {
		recover()
	}()

	mutex.Lock()
	defer mutex.Unlock()

	for i, lc := range *m {
		if lc == nil {
			(*m)[i] = c
			return i
		}
	}

	*m = append(*m, c)
	return len(*m) - 1
}

func (l *WSGetInfoLogic) WSGetInfo(w http.ResponseWriter, r *http.Request) {
	if !websocket.IsWebSocketUpgrade(r) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("connect ok"))
		return
	}

	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if config.BackendConfig.GetModeFromRequests(r) != config.RunModeDevelop {
				origin := r.Header.Get("Origin")

				if len(origin) == 0 { // 同源请求
					return true
				}

				for _, o := range config.BackendConfig.User.Origin {
					if origin == o {
						return true
					}
				}

				u, err := url.Parse(origin)
				if err != nil {
					return false
				}

				hostname := u.Hostname()
				for _, w := range model.WebsiteList() {
					if w.Status == db.WebsiteStatusBanned {
						continue
					}
					for _, d := range w.Domain {
						if d.Domain == hostname {
							return true
						}
					}
				}

				return false
			}
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	Connect(conn)
}

func Connect(conn *websocket.Conn) {
	userID := "" // len(userID) != 0, 必须要设置下面的值
	token := ""
	tokenType := 0
	webID := int64(0)

	ch := make(chan websocket2.WSMessage, 20)

	var closerFuncMutex = new(sync.Mutex)
	closerFunc := make(map[string]func(), 10)

	getCloser := func(name string) (func(), bool) {
		closerFuncMutex.Lock()
		defer closerFuncMutex.Unlock()

		f, ok := closerFunc[name]
		if !ok {
			return nil, false
		}

		return f, true
	}

	deleteCloser := func(name string) {
		closerFuncMutex.Lock()
		defer closerFuncMutex.Unlock()

		delete(closerFunc, name)
	}

	setCloser := func(name string, f func()) {
		closerFuncMutex.Lock()
		defer closerFuncMutex.Unlock()

		closerFunc[name] = f
	}

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
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.Close,
						})
					}()
					return true
				}

				return false
			}() {
				return // 退出go
			}

			if func() bool {
				defer utils.Recover(logger.Logger, nil, "ws check token")

				_, _, err := jwt.ParserUserToken(context.Background(), token)
				if err != nil {
					jwt.LogoutToken(token)
					return true
				}

				return false
			}() {
				return // 退出go
			}

			if len(token) != 0 {
				jwt.SetLogin(context.Background(), token, jwt.LoginWS)
			}
		}
	}()

	go func() {
		defer signalexit.AddExitFuncAsDefer(func() {
			defer utils.Recover(logger.Logger, nil, "ws reader")
			websocket2.WriteMessage(ch, websocket2.WSMessage{
				Code: userwstype.Close,
			})
		})()

		for {
			var msg websocket2.WSInMessage
			err := conn.ReadJSON(&msg)
			if utils.IsNetClose(err) || errors.Is(err, io.ErrUnexpectedEOF) || websocket.IsUnexpectedCloseError(err) {
				return
			} else if err != nil {
				return
			}

			var user *db.User
			userModel := db.NewUserModel(mysql.MySQLConn)

			if len(userID) == 0 {
				user = nil
			} else {
				user, err = userModel.FindOneByUidWithoutDelete(context.Background(), userID)
				if errors.Is(err, db.ErrNotFound) || err != nil {
					return
				}
			}

			switch msg.Code {
			default:
				websocket2.WriteMessage(ch, websocket2.WSMessage{
					Code: userwstype.BadCode,
				})
			case userwstype.Bye:
				websocket2.WriteMessage(ch, websocket2.WSMessage{
					Code: userwstype.Close,
				})
			case userwstype.Token:
				func() {
					data := struct {
						Token string `json:"token"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.BadData,
							Data: err.Error(),
						})
						return
					}

					if len(userID) != 0 {
						if token != data.Token {
							websocket2.WriteMessage(ch, websocket2.WSMessage{
								Code: userwstype.BadData,
								Data: "double set",
							})
						}
						return
					}

					token = data.Token
					userData, _, err := jwt.ParserUserToken(context.Background(), token)
					if err != nil {
						token = ""

						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.BadToken,
							Data: err.Error(),
						})
						return
					} else {
						userID = userData.UserID
						tokenType = userData.SubType
						webID = userData.WebsiteID
					}

					_, mysqlErr := userModel.FindOneByUidWithoutDelete(context.Background(), userID)
					if errors.Is(mysqlErr, db.ErrNotFound) {
						userID = ""
						token = ""
						tokenType = 0
						webID = 0

						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.BadToken,
							Data: "user not found",
						})
						return
					} else if mysqlErr != nil {
						userID = ""
						token = ""
						tokenType = 0
						webID = 0

						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.BadToken,
							Data: mysqlErr.Error(),
						})
						return
					}

					f, ok := getCloser("TokenInfo")
					if ok {
						func() {
							defer utils.Recover(logger.Logger, nil, "ws get token info error")
							f()
						}()

						deleteCloser("TokenInfo")
					}

					func() {
						index := AppendCh(websocket2.TokenConnMap, token, ch, websocket2.TokenConnMapMutex)
						setCloser("TokenInfo", func() {
							DeleteCh(websocket2.TokenConnMap, token, index, websocket2.TokenConnMapMutex)
						})
					}()

					jwt.SetLogin(context.Background(), token, jwt.LoginWS)
					websocket2.WriteMessage(ch, websocket2.WSMessage{
						Code: userwstype.Pong,
					})
				}()
			case userwstype.Ping:
				websocket2.WriteMessage(ch, websocket2.WSMessage{
					Code: userwstype.Pong,
				})

				func() {
					defer func() {
						recover()
					}()

					deadlineMutex.Lock()
					defer deadlineMutex.Unlock()

					deadline = time.Now().Add(time.Minute * 5)
				}()
			case userwstype.GetUserInfo:
				func() {
					if user == nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WithoutToken,
						})
						return
					}

					if webID != 0 {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WebsiteNotAllow,
						})
						return
					}

					f1, ok := getCloser("UserInfo")
					if ok {
						func() {
							defer utils.Recover(logger.Logger, nil, "ws get user info error")
							f1()
						}()

						deleteCloser("UserInfo")
					}

					f2, ok := getCloser("RoleInfo")
					if ok {
						func() {
							defer utils.Recover(logger.Logger, nil, "ws get role info error")
							f2()
						}()

						deleteCloser("RoleInfo")
					}

					index1 := AppendCh(websocket2.UserConnMap, user.Id, ch, websocket2.UserConnMapMutex)
					setCloser("UserInfo", func() {
						index1 = DeleteCh(websocket2.UserConnMap, user.Id, index1, websocket2.UserConnMapMutex)
					})

					index2 := AppendCh(websocket2.RoleConnMap, user.RoleId, ch, websocket2.RoleConnMapMutex)
					setCloser("RoleInfo", func() {
						index2 = DeleteCh(websocket2.RoleConnMap, user.RoleId, index2, websocket2.RoleConnMapMutex)
					})

					role := action.GetRole(user.RoleId, user.IsAdmin)
					db.UpdateUser(user.Id, mysql.MySQLConn, ch) // 推送一次用户信息
					urlpath.UpdateRole(role, ch)                // 推送一次角色消息
				}()
			case userwstype.GetWalletInfo:
				func() {
					if user == nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WithoutToken,
						})
						return
					}

					if webID != 0 {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WebsiteNotAllow,
						})
						return
					}

					f, ok := getCloser("WalletInfo")
					if ok {
						func() {
							defer utils.Recover(logger.Logger, nil, "ws get wallet info error")
							f()
						}()

						deleteCloser("WalletInfo")
					}

					index := AppendCh(websocket2.WalletConnMap, user.WalletId, ch, websocket2.WalletConnMapMutex)
					setCloser("WalletInfo", func() {
						index = DeleteCh(websocket2.WalletConnMap, user.WalletId, index, websocket2.WalletConnMapMutex)
					})

					wallet, err := balance.QueryBalance(context.Background(), user.Id)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.UpdateWalletInfo,
							Data: balance.UserBalance{
								Balance:      0,
								WaitBalance:  0,
								Cny:          0,
								NotBilled:    0,
								Billed:       0,
								HasBilled:    0,
								Withdraw:     0,
								WaitWithdraw: 0,
								NotWithdraw:  0,
								HasWithdraw:  0,
								WalletID:     0,
							},
						})
					} else {
						balance.UpdateWallet(wallet, ch) // 推送一次用户信息
					}
				}()
			case userwstype.GetAnnouncement:
				func() {
					if user == nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WithoutToken,
						})
						return
					}

					if webID != 0 {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WebsiteNotAllow,
						})
						return
					}

					f, ok := getCloser("AnnouncementInfo")
					if ok {
						func() {
							defer utils.Recover(logger.Logger, nil, "ws get announcement info error")
							f()
						}()

						deleteCloser("AnnouncementInfo")
					}

					index := AppendChToList(&websocket2.AnnouncementConnList, ch, websocket2.AnnouncementConnListMutex)
					setCloser("AnnouncementInfo", func() {
						index = DeleteChFromList(websocket2.AnnouncementConnList, index, websocket2.AnnouncementConnListMutex)
					})

					websocket2.WriteMessage(ch, websocket2.WSMessage{
						Code: userwstype.Pong,
					})
				}()
			case userwstype.GetMessage:
				func() {
					if user == nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WithoutToken,
						})
						return
					}

					f, ok := getCloser("MessageInfo")
					if ok {
						func() {
							defer utils.Recover(logger.Logger, nil, "ws get message info error")
							f()
						}()

						deleteCloser("MessageInfo")
					}

					index := AppendCh(websocket2.MessageConnMap, user.Id, ch, websocket2.MessageConnMapMutex)
					setCloser("MessageInfo", func() {
						index = DeleteCh(websocket2.MessageConnMap, user.Id, index, websocket2.MessageConnMapMutex)
					})

					websocket2.WriteMessage(ch, websocket2.WSMessage{
						Code: userwstype.Pong,
					})
				}()
			case userwstype.GetOrder:
				func() { // 封装到函数里面，可以使用return终止执行
					if user == nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WithoutToken,
						})
						return
					}

					data := struct {
						OrderID string `json:"orderID"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.BadData,
							Data: err.Error(),
						})
						return
					}

					orderModel := db.NewWorkOrderModel(mysql.MySQLConn)
					order, mysqlErr := orderModel.FindOneByUidWithoutDelete(context.Background(), data.OrderID)
					if errors.Is(mysqlErr, db.ErrNotFound) {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.BadData,
							Data: "order not found",
						})
						return
					} else if mysqlErr != nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.BadData,
							Data: mysqlErr.Error(),
						})
						return
					}

					web := action.GetWebsite(webID)
					matcher, ok := urlpath.CheckWebsiteUrlPath("/api/v1/website/msg/order", http.MethodPost)
					if !ok || !urlpath.CheckWebsiteMatherPermission(matcher, web) {
						if order.UserId != user.Id {
							websocket2.WriteMessage(ch, websocket2.WSMessage{
								Code: userwstype.BadData,
								Data: "order not found",
							})
							return
						}
					}

					if webID != 0 && order.FromId != webID {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WebsiteNotAllow,
						})
						return
					}

					funcKey := fmt.Sprintf("OrderReplyInfo-%s", data.OrderID)
					f, ok := getCloser(funcKey)
					if ok {
						func() {
							defer utils.Recover(logger.Logger, nil, "ws get order reply info info error")
							f()
						}()

						deleteCloser(funcKey)
					}

					index := AppendCh(websocket2.OrderConnMap, data.OrderID, ch, websocket2.OrderConnMapMutex)
					setCloser(funcKey, func() {
						index = DeleteCh(websocket2.OrderConnMap, data.OrderID, index, websocket2.OrderConnMapMutex)
					})

					websocket2.WriteMessage(ch, websocket2.WSMessage{
						Code: userwstype.Pong,
					})
				}()
			case userwstype.CloseOrder:
				func() { // 封装到函数里面，可以使用return终止执行
					if user == nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.WithoutToken,
						})
						return
					}

					data := struct {
						OrderID string `json:"orderID"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSMessage{
							Code: userwstype.BadData,
							Data: err.Error(),
						})
						return
					}

					funcKey := fmt.Sprintf("OrderReplyInfo-%s", data.OrderID)
					f, ok := getCloser(funcKey)
					if ok {
						func() {
							defer utils.Recover(logger.Logger, nil, "ws get order reply info info error")
							f()
						}()

						deleteCloser(funcKey)
					}

					websocket2.WriteMessage(ch, websocket2.WSMessage{
						Code: userwstype.Pong,
					})
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

			if len(token) != 0 {
				jwt.DelLogin(context.Background(), token, jwt.LoginWS)
			}
		})()

		defer utils.Recover(logger.Logger, nil, "ws writer")

		for {
			msg := <-ch

			if msg.Code == userwstype.LogoutToken ||
				msg.Code == userwstype.UpdateUserInfo ||
				msg.Code == userwstype.UpdateWalletInfo ||
				msg.Code == userwstype.UpdateRoleInfo ||
				msg.Code == userwstype.NewAnnouncement ||
				msg.Code == userwstype.UpdateAnnouncement ||
				msg.Code == userwstype.DeleteAnnouncement ||
				msg.Code == userwstype.UpdateMessage ||
				msg.Code == userwstype.ReadMessage ||
				msg.Code == userwstype.NewOrderReply ||
				msg.Code == userwstype.UpdateOrder ||
				msg.Code == userwstype.RoleChange {
				if len(userID) == 0 {
					continue
				}
			}

			if msg.Code == userwstype.UpdateUserInfo ||
				msg.Code == userwstype.UpdateWalletInfo ||
				msg.Code == userwstype.UpdateRoleInfo ||
				msg.Code == userwstype.NewAnnouncement ||
				msg.Code == userwstype.UpdateAnnouncement ||
				msg.Code == userwstype.DeleteAnnouncement ||
				msg.Code == userwstype.RoleChange {
				if webID != 0 {
					continue
				}
			}

			switch msg.Code {
			case userwstype.Close:
				return
			case userwstype.LogoutToken:
				websocket2.WriteMessage(ch, websocket2.WSMessage{
					Code: userwstype.Close,
				})
			case userwstype.RoleChange:
				userModel := db.NewUserModel(mysql.MySQLConn)
				user, err := userModel.FindOneByUidWithoutDelete(context.Background(), userID)
				if errors.Is(err, db.ErrNotFound) || err != nil {
					continue
				}

				f, ok := getCloser("RoleInfo")
				if ok {
					func() {
						defer utils.Recover(logger.Logger, nil, "ws get role info error")
						f()
					}()

					deleteCloser("RoleInfo")
				}

				func() {
					index := AppendCh(websocket2.RoleConnMap, user.RoleId, ch, websocket2.RoleConnMapMutex)
					setCloser("RoleInfo", func() {
						index = DeleteCh(websocket2.RoleConnMap, user.RoleId, index, websocket2.RoleConnMapMutex)
					})
				}()
			case userwstype.UpdateMessage:
				if webID != 0 && webID != msg.WebID {
					continue
				}
			case userwstype.ReadMessage:
				if webID != 0 && webID != msg.WebID {
					continue
				}
			case userwstype.NewOrderReply:
				if webID != 0 && webID != msg.WebID {
					continue
				}
			case userwstype.UpdateOrder:
				if webID != 0 && webID != msg.WebID {
					continue
				}
			}

			msg.ClearData(tokenType)

			err := conn.WriteJSON(msg)
			if utils.IsNetClose(err) || websocket.IsUnexpectedCloseError(err) {
				return
			} else if err != nil {
				return
			}
		}
	}()

}
