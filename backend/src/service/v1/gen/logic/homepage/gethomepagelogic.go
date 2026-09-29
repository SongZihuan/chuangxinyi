package homepage

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetHomePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHomePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHomePageLogic {
	return &GetHomePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHomePageLogic) GetHomePage(req *types.GetHomePageReq) (resp *types.GetHomePageResp, err error) {
	userModel := db.NewUserModel(mysql.MySQLConn)
	homepageModel := db.NewHomepageModel(mysql.MySQLConn)

	user, err := userModel.FindOneByUidWithoutDelete(l.ctx, req.UserID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.GetHomePageResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotHomePage, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	homepage, err := homepageModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.GetHomePageResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotHomePage, "用户没有主页"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if homepage.Close {
		return &types.GetHomePageResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotHomePage, "用户主页关闭"),
		}, nil
	}

	sex := "保密"
	if homepage.Man.Valid {
		if homepage.Man.Bool {
			sex = "男"
		} else {
			sex = "女"
		}
	} else {
		sex = "保密"
	}

	return &types.GetHomePageResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetHomePageData{
			HomePage: types.HomePage{
				CompanyName:  homepage.Company.String,
				Introduction: homepage.Introduction.String,
				Address:      homepage.Address.String,
				Phone:        homepage.Phone.String,
				Email:        homepage.Email.String,
				WeChat:       homepage.Wechat.String,
				QQ:           homepage.Qq.String,
				Sex:          sex,
				Link:         homepage.Link.String,
				Industry:     homepage.Industry.String,
				Position:     homepage.Position.String,
			},
		},
	}, nil
}
