package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UserHeaderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserHeaderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHeaderLogic {
	return &UserHeaderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserHeaderLogic) UserHeader(req *types.UserHeaderReq, w http.ResponseWriter, r *http.Request) error {
	userModel := db.NewUserModel(mysql.MySQLConn)
	headerModel := db.NewHeaderModel(mysql.MySQLConn)

	var url string
	user, err := userModel.FindOneByUidWithoutDelete(l.ctx, req.UserUID)
	if errors.Is(err, db.ErrNotFound) {
		url, err = oss.GetHeader(oss.DefaultHeaderName, true)
		if err != nil {
			return respmsg.OSSError.WarpQuick(err)
		}
	} else if err != nil {
		return respmsg.MySQLSystemError.WarpQuick(err)
	} else {
		header, err := headerModel.FindByUserID(l.ctx, user.Id)
		if errors.Is(err, db.ErrNotFound) {
			url, err = oss.GetHeader(oss.DefaultHeaderName, true)
			if err != nil {
				return respmsg.OSSError.WarpQuick(err)
			}
		} else if err != nil {
			return respmsg.MySQLSystemError.WarpQuick(err)
		} else if header.Header.Valid {
			url, err = oss.GetHeader(header.Header.String, true)
			if err != nil {
				return respmsg.OSSError.WarpQuick(err)
			}
		} else {
			url, err = oss.GetHeader(oss.DefaultHeaderName, true)
			if err != nil {
				return respmsg.OSSError.WarpQuick(err)
			}
		}
	}

	http.Redirect(w, r, url, http.StatusFound)
	return nil
}
