package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetCommunicateListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommunicateListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommunicateListLogic {
	return &GetCommunicateListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommunicateListLogic) GetCommunicateList(req *types.GetCommunicateListReq) (resp *types.GetOrderCommunicateListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	workOrderModel := db.NewWorkOrderModel(mysql.MySQLConn)
	workOrderCommunicateModel := db.NewWorkOrderCommunicateModel(mysql.MySQLConn)
	workOrderCommunicateFileModel := db.NewWorkOrderCommunicateFileModel(mysql.MySQLConn)

	order, err := workOrderModel.FindOneByUidWithoutDelete(l.ctx, req.OrderID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.GetOrderCommunicateListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WorkOrderNotFound, "工单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if order.UserId != user.Id {
		return &types.GetOrderCommunicateListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WorkOrderNotFound, "工单不属于用户"),
		}, nil
	} else if web.ID != warp.UserCenterWebsite && web.ID != order.FromId {
		return &types.GetOrderCommunicateListResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WorkOrderNotFound, "工单不属于外站"),
		}, nil
	}

	communicateList, err := workOrderCommunicateModel.GetList(l.ctx, order.Id, req.Page, req.PageSize, 0, 0, db.TimeCreateAt)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := workOrderCommunicateModel.GetCount(l.ctx, order.Id, 0, 0, db.TimeCreateAt)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.WorkOrderCommunicate, 0, len(communicateList))
	for _, c := range communicateList {
		fileList, err := workOrderCommunicateFileModel.GetList(l.ctx, c.Id)
		if err != nil {
			logger.Logger.Error("mysql errors: %s", err.Error())
			fileList = []db.WorkOrderCommunicateFile{}
		}

		respFileList := make([]types.WorkOrderCommunicateFile, 0, len(fileList))
		for _, f := range fileList {
			respFileList = append(respFileList, types.WorkOrderCommunicateFile{
				Fid: f.Fid,
			})
		}

		respList = append(respList, types.WorkOrderCommunicate{
			ID:       c.Id,
			Content:  c.Content,
			From:     c.From,
			CreateAt: c.CreateAt.Unix(),
			File:     respFileList,
		})
	}

	return &types.GetOrderCommunicateListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetOrderCommunicateListData{
			Count:       count,
			Communicate: respList,
		},
	}, nil
}
