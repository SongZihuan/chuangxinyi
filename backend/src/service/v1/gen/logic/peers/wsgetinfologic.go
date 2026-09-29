package peers

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/peername"
	websocket2 "gitee.com/wuntsong-auth/backend/src/global/websocket"
	"gitee.com/wuntsong-auth/backend/src/peers"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/gorilla/websocket"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"
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
	ip, ok := r.Context().Value("X-Real-IP").(string)
	if !ok {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusInternalServerError, &types.RespEmpty{
			Resp: respmsg.GetRespByErrorWithCode(l.ctx, respmsg.PeerDenyCode, respmsg.SystemError, respmsg.BadContextError.New("X-Real-IP")),
		})
		return
	}

	otherSidePeerName := r.URL.Query().Get("peername")
	if len(otherSidePeerName) == 0 {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusBadRequest, &types.RespEmpty{
			Resp: respmsg.GetRespByMsgWithCode(l.ctx, respmsg.PeerDenyCode, respmsg.NotPeerName, "没有配置peername入参"),
		})
		return
	}

	selfSidePeerName := r.URL.Query().Get("you")
	if len(selfSidePeerName) == 0 {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusBadRequest, &types.RespEmpty{
			Resp: respmsg.GetRespByMsgWithCode(l.ctx, respmsg.PeerDenyCode, respmsg.NotPeerName, "没有配置you入参"),
		})
		return
	} else if selfSidePeerName != peername.PeerName {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusBadRequest, &types.RespEmpty{
			Resp: respmsg.GetRespByMsgWithCode(l.ctx, respmsg.PeerDenyCode, respmsg.PeerNameError, "入参you不对"),
		})
		return
	}

	if !IsRightPeerName(r.Context(), otherSidePeerName, ip) {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusBadRequest, &types.RespEmpty{
			Resp: respmsg.GetRespByMsgWithCode(l.ctx, respmsg.PeerDenyCode, respmsg.BadPeerIP, "错误的peer ip: %s", ip),
		})
		return
	}

	if !websocket.IsWebSocketUpgrade(r) {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusOK, &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		})
		return
	}

	if func() bool {
		websocket2.PeersConnMapMutex.Lock()
		defer websocket2.PeersConnMapMutex.Unlock()

		_, ok := websocket2.PeersConnMap[otherSidePeerName]
		if ok {
			httpx.WriteJsonCtx(context.Background(), w, http.StatusOK, &types.RespEmpty{
				Resp: respmsg.GetRespByMsgWithCode(l.ctx, respmsg.PeerDenyCode, respmsg.DoubleConnect, "双重连接"),
			})
			return true
		}
		return false
	}() {
		return
	}

	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		httpx.WriteJsonCtx(context.Background(), w, http.StatusInternalServerError, &types.RespEmpty{
			Resp: respmsg.GetRespByErrorWithCode(l.ctx, respmsg.PeerDenyCode, respmsg.SystemError, errors.WarpQuick(err), "websocket建立连接错误"),
		})
		return
	}

	peers.NewPeerConn(conn, otherSidePeerName)
}

func IsRightPeerName(ctx context.Context, otherSidePeerName string, ip string) bool {
	dataString, err := redis.Get(ctx, fmt.Sprintf("service:peers:%s:%s", config.BackendConfig.User.Group, otherSidePeerName)).Result()
	if err != nil {
		return false
	}

	var data peers.Peer
	err = utils.JsonUnmarshal([]byte(dataString), &data)
	if err != nil {
		return false
	}

	if data.Name != otherSidePeerName {
		return false
	}

	if utils.CheckIPInList(ip, data.IP, true, true, true) {
		return true
	}

	return false
}
