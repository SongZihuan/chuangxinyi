package peers

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/peername"
	websocket2 "gitee.com/wuntsong-auth/backend/src/global/websocket"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/peers/wstype"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	"gitee.com/wuntsong-auth/backend/src/signalexit"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/gorilla/websocket"
	"github.com/wuntsong-org/wterrors"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Peer struct {
	Name    string   `json:"name"`
	Url     string   `json:"url"`
	IP      []string `json:"IP"`
	Connect []string `json:"connect"`
}

func ConnectPeers() errors.WTError {
	if peername.PeerName == peername.Single { // 不启动peer链接
		return nil
	}

	if len(peername.Eth0IP) == 0 || len(peername.SelfIP) == 0 {
		return errors.Errorf("bad eth ip")
	}

	var PeerUrl string
	if strings.Contains(peername.Eth0IP, ":") {
		PeerUrl = fmt.Sprintf("ws://[%s]:%d/api/v1/peers/ws", peername.Eth0IP, config.BackendConfig.User.Port)
	} else {
		PeerUrl = fmt.Sprintf("ws://%s:%d/api/v1/peers/ws", peername.Eth0IP, config.BackendConfig.User.Port)
	}

	key := fmt.Sprintf("service:peers:%s:%s", config.BackendConfig.User.Group, peername.PeerName)

	var connectList = make([]string, 0)

	signalexit.AddExitByFunc(func(ctx context.Context, _ os.Signal) context.Context {
		_ = redis.Del(context.Background(), key)
		return context.WithValue(ctx, "PeerDelete", true)
	})

	go func() {
		defer func() {
			_ = redis.Del(context.Background(), key)
		}()

		for {
			func() {
				defer utils.Recover(logger.Logger, nil, "check peers connect")

				peerDataByte, err := utils.JsonMarshal(Peer{
					Name:    peername.PeerName,
					Url:     PeerUrl,
					IP:      peername.SelfIP,
					Connect: connectList, // 更新在PeersConnMapMutex上锁之后
				})
				if err == nil {
					peerData := string(peerDataByte)
					_ = redis.Set(context.Background(), key, peerData, time.Second*30)
				}

				peers, redisErr := redis.Keys(context.Background(), fmt.Sprintf("service:peers:%s:*", config.BackendConfig.User.Group)).Result()
				if redisErr != nil {
					return
				}

				websocket2.PeersConnMapMutex.Lock()
				defer websocket2.PeersConnMapMutex.Unlock()

				for _, p := range peers {
					if p == key {
						continue
					}

					keySpilt := strings.Split(p, ":")
					if len(keySpilt) != 4 {
						continue
					}

					dataString, err := redis.Get(context.Background(), p).Result()
					if err != nil {
						continue
					}

					var data Peer
					err = utils.JsonUnmarshal([]byte(dataString), &data)
					if err != nil {
						continue
					}

					if data.Name == peername.PeerName || keySpilt[3] != data.Name {
						continue
					}

					if keySpilt[3] != data.Name {
						continue
					}

					if len(data.Url) == 0 || len(data.IP) == 0 {
						continue
					}

					_, ok := websocket2.PeersConnMap[data.Name]
					if !ok {
						go func() { // 防止死锁
							_ = Connect(data.Url, data.Name)
						}()
					}
				}

				connectList = make([]string, 0, len(websocket2.PeersConnMap))
				for c := range websocket2.PeersConnMap {
					connectList = append(connectList, c)
				}
			}()

			time.Sleep(5 * time.Second)
		}

	}()

	return nil
}

func SendPeersPing() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)

			func() {
				defer utils.Recover(logger.Logger, nil, "send peers ping")

				websocket2.PeersConnMapMutex.Lock()
				defer websocket2.PeersConnMapMutex.Unlock()

				for _, p := range websocket2.PeersConnMap {
					data := struct {
						Host string `json:"host"`
					}{Host: peername.PeerName} // ping时发送本机peerName

					dataByte, err := utils.JsonMarshal(data)
					if err != nil {
						continue
					}

					websocket2.WriteMessage(p, websocket2.WSPeersMessage{
						Code: wstype.PeersPing,
						Data: string(dataByte),
					})
				}

			}()

		}
	}()
}

func Connect(u string, otherSidePeerName string) errors.WTError {
	uu, err := url.Parse(u)

	if err != nil {
		logger.Logger.Error("connect ws error: %s", err.Error())
		return errors.WarpQuick(err)
	}

	if uu.Scheme != "ws" && uu.Scheme != "wss" {
		logger.Logger.Error("connect ws error: not websocket scheme")
		return errors.Errorf("bad url")
	}

	q := url.Values{}
	q.Add("peername", peername.PeerName)
	q.Add("you", otherSidePeerName)

	connUrl := fmt.Sprintf("%s?%s", u, q.Encode())

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(connUrl, nil) // 不用defer close
	if err != nil {
		return errors.WarpQuick(err)
	}

	NewPeerConn(conn, otherSidePeerName)
	return nil
}

func NewPeerConn(conn *websocket.Conn, otherSidePeerName string) {
	ch := make(chan websocket2.WSPeersMessage, 20)

	if otherSidePeerName == peername.PeerName {
		_ = conn.Close()
		return
	}

	peersCount := 0
	if func() bool {
		websocket2.PeersConnMapMutex.Lock()
		defer websocket2.PeersConnMapMutex.Unlock()

		_, ok := websocket2.PeersConnMap[otherSidePeerName]
		if ok {
			_ = conn.Close()
			return true
		}

		websocket2.PeersConnMap[otherSidePeerName] = ch
		peersCount = len(websocket2.PeersConnMap)
		return false
	}() {
		return
	}

	websocket2.PeersCounter.Set(int64(peersCount)) // 不要放在上述函数中，避免同时两把锁
	logger.Logger.WXInfo("%s connect %s success", otherSidePeerName, peername.PeerName)

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
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersClose,
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
			websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
				Code: wstype.PeersClose,
			})
			// 不需要 _ = conn.Close()
			closer()
			deadlineNow()
		})()

		for {
			var msg websocket2.WSPeersMessage
			err := conn.ReadJSON(&msg)
			if utils.IsNetClose(err) || errors.Is(err, io.ErrUnexpectedEOF) || websocket.IsUnexpectedCloseError(err) {
				return
			} else if err != nil {
				return
			}

			key := fmt.Sprintf("peers:%s:%s", otherSidePeerName, msg.ID)
			res, err := redis.Exists(context.Background(), key).Result()
			if err != nil {
				logger.Logger.Error("redis error: %s", err.Error())
				continue
			} else if res == 1 {
				continue // 消息已经处理过了
			}

			_ = redis.Set(context.Background(), key, "1", time.Minute*5)

			idSplit := strings.Split(msg.ID, "@")
			if len(idSplit) != 2 {
				continue
			}

			if idSplit[0] == peername.PeerName { // 自己发送的消息
				continue
			}

			msgMicro, err := strconv.ParseInt(idSplit[1], 10, 64)
			if err != nil {
				continue
			}

			msgSendTime := time.UnixMicro(msgMicro)
			if msgSendTime.Add(time.Minute * 3).Before(time.Now()) { // 过时的消息
				continue
			}

			if NeedForward(msg.Code) {
				go func() {
					websocket2.PeersConnMapMutex.Lock()
					defer websocket2.PeersConnMapMutex.Unlock()

					for _, p := range websocket2.PeersConnMap {
						if p != ch {
							websocket2.WriteMessage(p, msg)
						}
					}
				}()
			}

			switch msg.Code {
			default:
				websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
					Code: wstype.PeersBadCode,
				})
			case wstype.PeersPing:
				func() {
					defer func() {
						recover()
					}()

					data := struct {
						Host string `json:"host"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					if data.Host != otherSidePeerName {
						logger.Logger.Error("peers ping want %s but get %s", peername.PeerName, data.Host)
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersClose,
						})
						return
					}

					deadlineMutex.Lock()
					defer deadlineMutex.Unlock()

					deadline = time.Now().Add(time.Minute * 5)
				}()
			case wstype.PeersLogoutToken:
				func() {
					data := struct {
						Token string `json:"token"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					wsMsg := websocket2.WSMessage{}
					err = utils.JsonUnmarshal([]byte(msg.Message), &wsMsg)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					jwt.LogoutTokenByMsg(data.Token, wsMsg)
				}()
			case wstype.PeersUpdateUserInfo:
				func() {
					data := struct {
						UserID int64 `json:"userID"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					wsMsg := websocket2.WSMessage{}
					err = utils.JsonUnmarshal([]byte(msg.Message), &wsMsg)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					db.UpdateUserByMsg(data.UserID, wsMsg)
				}()
			case wstype.PeersUpdateRoleInfo:
				func() {
					data := struct {
						RoleID int64 `json:"roleID"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					wsMsg := websocket2.WSMessage{}
					err = utils.JsonUnmarshal([]byte(msg.Message), &wsMsg)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					urlpath.UpdateRoleByMsg(data.RoleID, wsMsg)
				}()
			case wstype.PeersUpdateWalletInfo:
				func() {
					data := struct {
						WalletID int64 `json:"walletID"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					wsMsg := websocket2.WSMessage{}
					err = utils.JsonUnmarshal([]byte(msg.Message), &wsMsg)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					balance.UpdateWalletByMsg(data.WalletID, wsMsg)
				}()
			case wstype.PeersUpdateMessage:
				func() {
					data := struct {
						UserID int64 `json:"userID"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					wsMsg := websocket2.WSMessage{}
					err = utils.JsonUnmarshal([]byte(msg.Message), &wsMsg)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					db.UpdateMessageByMsg(data.UserID, wsMsg)
				}()
			case wstype.PeersUpdateOrder:
				func() {
					data := struct {
						OrderID string `json:"orderID"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					wsMsg := websocket2.WSMessage{}
					err = utils.JsonUnmarshal([]byte(msg.Message), &wsMsg)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					db.UpdateWorkOrderByMsg(data.OrderID, wsMsg)
				}()
			case wstype.PeersUpdateAnnouncement:
				func() {
					wsMsg := websocket2.WSMessage{}
					err = utils.JsonUnmarshal([]byte(msg.Message), &wsMsg)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					db.UpdateAnnouncementByMsg(wsMsg)
				}()
			case wstype.PeersRoleChange:
				func() {
					data := struct {
						UserID int64 `json:"userID"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					wsMsg := websocket2.WSMessage{}
					err = utils.JsonUnmarshal([]byte(msg.Message), &wsMsg)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					urlpath.ChangeRoleByMsg(data.UserID, wsMsg)
				}()
			case wstype.PeersRoleChangeByDelete:
				func() {
					data := struct {
						RoleID int64 `json:"roleID"`
					}{}

					err := utils.JsonUnmarshal([]byte(msg.Data), &data)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					wsMsg := websocket2.WSMessage{}
					err = utils.JsonUnmarshal([]byte(msg.Message), &wsMsg)
					if err != nil {
						websocket2.WriteMessage(ch, websocket2.WSPeersMessage{
							Code: wstype.PeersBadData,
							Data: websocket2.WriteJson(struct {
								Err string `json:"err"`
							}{Err: err.Error()}),
						})
						return
					}

					urlpath.ChangeRoleByDeleteByMsg(data.RoleID, wsMsg)
				}()
			case wstype.PeersRoleDBUpdate:
				cron.RoleUpdateHandler(false)
			case wstype.PeersMenuDBUpdate:
				cron.MenuUpdateHandler(false)
			case wstype.PeersPermissionDBUpdate:
				cron.PermissionUpdateHandler(false)
			case wstype.PeersUrlPathDBUpdate:
				cron.PathUpdateHandler(false)
			case wstype.PeersWebsiteUrlPathDBUpdate:
				cron.WebsitePathUpdateHandler(false)
			case wstype.PeersWebsitePermissionDBUpdate:
				cron.WebsitePermissionUpdateHandler(false)
			case wstype.PeersWebsiteDBUpdate:
				cron.WebsiteUpdateHandler(false)
			case wstype.PeersApplicationDBUpdate:
				cron.ApplicationUpdateHandler(false)
			}

		}
	}()

	go func() {
		defer signalexit.AddExitFuncAsDefer(func() {
			defer utils.Recover(logger.Logger, nil, "ws writer")
			_ = conn.Close()
			closer()
			peersCount := 0
			func() {
				websocket2.PeersConnMapMutex.Lock()
				defer websocket2.PeersConnMapMutex.Unlock()
				delete(websocket2.PeersConnMap, otherSidePeerName)
				peersCount = len(websocket2.PeersConnMap)
			}()
			websocket2.PeersCounter.Set(int64(peersCount)) // 不要放在上述函数中，避免同时两把锁
			deadlineNow()
			close(ch)
		})()

		defer utils.Recover(logger.Logger, nil, "ws writer")

		for {
			msg := <-ch

			switch msg.Code {
			case wstype.PeersClose:
				return
			}

			if len(msg.ID) == 0 {
				sendTime := time.Now()
				msg.ID = fmt.Sprintf("%s@%d", peername.PeerName, sendTime.UnixMicro())
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

func NeedForward(code string) bool {
	for _, c := range wstype.PeersForward {
		if c == code {
			return true
		}
	}

	return false
}
