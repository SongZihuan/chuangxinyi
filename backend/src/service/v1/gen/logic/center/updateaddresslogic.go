package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAddressLogic) UpdateAddress(req *types.UserUpdateAddressReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if len(req.Name) == 0 {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NameTooShort, "收件人必须提供"),
		}, nil
	}

	if len(req.Phone) != 0 && !utils.IsPhoneNumber(req.Phone) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPhone, "收件手机号错误"),
		}, nil
	}

	if len(req.Email) != 0 && !utils.IsEmailAddress(req.Email) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadEmail, "收件邮箱错误"),
		}, nil
	}

	a := &db.Address{
		UserId: user.Id,
		Name: sql.NullString{
			Valid:  len(req.Name) != 0,
			String: req.Name,
		},
		Phone: sql.NullString{
			Valid:  len(req.Phone) != 0,
			String: req.Phone,
		},
		Email: sql.NullString{
			Valid:  len(req.Email) != 0,
			String: req.Email,
		},
		Province: sql.NullString{
			Valid:  len(req.Province) != 0,
			String: req.Province,
		},
		City: sql.NullString{
			Valid:  len(req.City) != 0,
			String: req.City,
		},
		District: sql.NullString{
			Valid:  len(req.District) != 0,
			String: req.District,
		},
		Address: sql.NullString{
			Valid:  len(req.Address) != 0,
			String: req.Address,
		},
	}

	err = a.SetAreaCode(req.Area...)
	if errors.Is(err, db.BadArea) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadArea, "错误的地域代码"),
		}, nil
	} else if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadArea, errors.WarpQuick(err), "错误的地域信息"),
		}, nil
	}

	addressModel := db.NewAddressModel(mysql.MySQLConn)
	_, err = addressModel.InsertWithDelete(context.Background(), a)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户更新地址信息成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
