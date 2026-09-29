package admin_website

import (
	"context"
	"encoding/base64"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/google/uuid"
	errors "github.com/wuntsong-org/wterrors"
	"math/big"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AddWebsiteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddWebsiteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddWebsiteLogic {
	return &AddWebsiteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddWebsiteLogic) AddWebsite(req *types.AddWebsiteReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if !db.IsWebsiteStatus(req.Status) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadWebsiteStatus, "错误的外站状态"),
		}, nil
	}

	websiteUIDByte, success := redis.GenerateUUIDMore(l.ctx, "website", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		websiteModel := db.NewWebsiteModel(mysql.MySQLConn)
		_, err := websiteModel.FindOneByUIDWithoutDelete(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CreateWebsiteFail, "无法生成站点uuid"),
		}, nil
	}

	websiteModel := db.NewWebsiteModel(mysql.MySQLConn)

	count, err := websiteModel.GetCount(l.ctx)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if count > config.BackendConfig.MySQL.SystemResourceLimit {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "超出限额"),
		}, nil
	}

	pubkey, err := utils.DecodeReqBase64(req.PubKey, false)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadWebsitePubKey, errors.WarpQuick(err), "从base64解析公钥错误"),
		}, nil
	}

	_, err = utils.ReadRsaPublicKey(pubkey)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadWebsitePubKey, errors.WarpQuick(err), "读取rsa公钥错误"),
		}, nil
	}

	keyMap := make(map[string]string, len(req.KeyMap))
	for _, v := range req.KeyMap {
		keyMap[v.Label] = v.Value
	}

	keyMapJson, err := utils.JsonMarshal(keyMap)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadKeyMap, errors.WarpQuick(err), "编码KeyMap错误"),
		}, nil
	}

	var p big.Int
	for _, ps := range req.Policy {
		np, ok := (model.WebsitePermissionsSign())[ps]
		if ok && np.Status != db.WebsitePolicyStatusBanned {
			p = permission.AddPermission(p, np.Permission)
		}
	}

	_, err = websiteModel.Insert(l.ctx, &db.Website{
		Uid:        websiteUIDByte.String(),
		Name:       req.Name,
		Describe:   req.Describe,
		Keymap:     string(keyMapJson),
		Pubkey:     base64.StdEncoding.EncodeToString(pubkey),
		Agreement:  req.Agreement,
		Permission: p.Text(16),
		Status:     req.Status,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	cron.WebsiteUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员添加站点（%s）成功", req.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
