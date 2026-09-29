package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"gitee.com/wuntsong-auth/backend/src/yundun"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UploadUserInfoByJsonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadUserInfoByJsonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadUserInfoByJsonLogic {
	return &UploadUserInfoByJsonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadUserInfoByJsonLogic) UploadUserInfoByJson(req *types.UploadUserInfoByJson) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	phoneModel := db.NewPhoneModel(mysql.MySQLConn)

	phone, err := phoneModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SystemError, "系统错误，注册后没有手机号"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if !utils.IsValidIDCard(req.UserIDCard) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadUserInfo, "错误的用户身份证号码"),
		}, nil
	}

	if !utils.IsValidChineseName(req.UserName) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadUserInfo, "错误的用户名称"),
		}, nil
	}

	phoneRes, err := func() (bool, error) {
		defer utils.Recover(logger.Logger, nil, "")

		res1, err := yundun.CheckIDCard(req.UserName, req.UserIDCard)
		if err != nil {
			return false, err
		}

		if !res1 {
			// 身份证识别错误
			return false, errors.Errorf("bad id info")
		}

		res2, err := yundun.CheckPhone(req.UserName, req.UserIDCard, phone.Phone)
		if err != nil {
			return false, err
		}

		return res2, nil
	}()
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadIDCard, errors.WarpQuick(err), "用户身份信息检验失败"),
		}, nil
	}

	newIDCard := &db.Idcard{
		UserId:     user.Id,
		UserName:   req.UserName,
		UserIdCard: req.UserIDCard,
	}

	if phoneRes {
		newIDCard.Phone = sql.NullString{
			Valid:  true,
			String: phone.Phone,
		}
	}

	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	_, err = idcardModel.InsertWithDelete(context.Background(), newIDCard)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户更新使用者实名成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
