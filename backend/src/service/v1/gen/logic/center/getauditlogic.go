package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAuditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAuditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuditLogic {
	return &GetAuditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAuditLogic) GetAudit(req *types.UserGetAuditReq) (resp *types.UserGetAuditResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	var auditList []db.Audit
	var count int64

	auditModel := db.NewAuditModel(mysql.MySQLConn)
	if web.ID == warp.UserCenterWebsite {
		auditList, err = auditModel.GetList(l.ctx, user.Id, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.FromID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = auditModel.GetCount(l.ctx, user.Id, req.Src, req.StartTime, req.EndTime, req.TimeType, req.FromID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		auditList, err = auditModel.GetList(l.ctx, user.Id, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = auditModel.GetCount(l.ctx, user.Id, req.Src, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	respList := make([]types.Audit, 0, len(auditList))
	for _, a := range auditList {
		respList = append(respList, types.Audit{
			Content:  a.Content,
			From:     a.From,
			CreateAt: a.CreateAt.Unix(),
		})
	}

	return &types.UserGetAuditResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.UserGetAuditData{
			Count: count,
			Audit: respList,
		},
	}, nil
}
