package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetOrderFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderFileLogic {
	return &GetOrderFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderFileLogic) GetOrderFile(req *types.GetOrderFileReq, w http.ResponseWriter, r *http.Request) error {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return respmsg.BadContextError.New("X-Token-Website")
	}

	workOrderModel := db.NewWorkOrderModel(mysql.MySQLConn)
	workOrderCommunicateModel := db.NewWorkOrderCommunicateModel(mysql.MySQLConn)

	communicate, err := workOrderCommunicateModel.FindOneWithoutDelete(l.ctx, req.CommunicateID)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("work order not found")
	} else if err != nil {
		return respmsg.MySQLSystemError.WarpQuick(err)
	}

	workOrder, err := workOrderModel.FindOneWithoutDelete(l.ctx, communicate.OrderId)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("work order not found")
	} else if err != nil {
		return respmsg.MySQLSystemError.WarpQuick(err)
	} else if web.ID != warp.UserCenterWebsite && web.ID != workOrder.Id {
		return errors.Errorf("website not allow")
	} else if workOrder.UserId != user.Id {
		return errors.Errorf("work order not found")
	}

	url, err := oss.GetWorkOrderFile(workOrder, communicate, req.Fid, req.Download)
	if err != nil {
		return respmsg.OSSError.WarpQuick(err)
	}

	http.Redirect(w, r, url, http.StatusFound)
	return nil
}
