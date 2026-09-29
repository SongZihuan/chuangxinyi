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

type HeaderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHeaderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeaderLogic {
	return &HeaderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HeaderLogic) Header(req *types.HeaderReq, w http.ResponseWriter, r *http.Request) error {
	if len(req.Header) != 0 { // 判断header是否被删除
		headerModel := db.NewHeaderModel(mysql.MySQLConn)
		_, err := headerModel.FindByHeaderWithoutDelete(l.ctx, req.Header)
		if errors.Is(err, db.ErrNotFound) {
			req.Header = ""
		} else if err != nil {
			return respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		req.Header = oss.DefaultHeaderName
	}

	url, err := oss.GetHeader(req.Header, true)
	if err != nil {
		return respmsg.OSSError.WarpQuick(err)
	}

	http.Redirect(w, r, url, http.StatusFound)
	return nil
}
