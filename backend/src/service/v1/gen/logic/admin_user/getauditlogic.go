package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

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

func (l *GetAuditLogic) GetAudit(req *types.AdminGetAuditReq) (resp *types.AdminGetAuditResp, err error) {
	var auditList []db.Audit
	var count int64

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	auditModel := db.NewAuditModel(mysql.MySQLConn)
	if web.ID == warp.UserCenterWebsite {
		if req.ID == 0 && len(req.UID) == 0 {
			auditList, err = auditModel.GetList(l.ctx, 0, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.FromID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = auditModel.GetCount(l.ctx, 0, req.Src, req.StartTime, req.EndTime, req.TimeType, req.FromID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		} else {
			user, err := GetUser(l.ctx, req.ID, req.UID, true)
			if errors.Is(err, UserNotFound) {
				return &types.AdminGetAuditResp{
					Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
				}, nil
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			auditList, err = auditModel.GetList(l.ctx, user.Id, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.FromID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = auditModel.GetCount(l.ctx, user.Id, req.Src, req.StartTime, req.EndTime, req.TimeType, req.FromID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		}
	} else {
		if req.ID == 0 && len(req.UID) == 0 {
			auditList, err = auditModel.GetList(l.ctx, 0, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = auditModel.GetCount(l.ctx, 0, req.Src, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		} else {
			user, err := GetUser(l.ctx, req.ID, req.UID, true)
			if errors.Is(err, UserNotFound) {
				return &types.AdminGetAuditResp{
					Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
				}, nil
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			auditList, err = auditModel.GetList(l.ctx, user.Id, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = auditModel.GetCount(l.ctx, user.Id, req.Src, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		}
	}

	respList := make([]types.AdminAudit, 0, len(auditList))
	for _, a := range auditList {
		respList = append(respList, types.AdminAudit{
			UserID:   a.UserId,
			Content:  a.Content,
			From:     a.From,
			FromID:   a.FromId,
			CreateAt: a.CreateAt.Unix(),
		})
	}

	return &types.AdminGetAuditResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetAuditData{
			Count: count,
			Audit: respList,
		},
	}, nil
}
