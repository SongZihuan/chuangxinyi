package admin_user

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

type GetUserHomepageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserHomepageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserHomepageLogic {
	return &GetUserHomepageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserHomepageLogic) GetUserHomepage(req *types.AdminGetUserReq) (resp *types.AdminGetHomePageResp, err error) {
	user, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminGetHomePageResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	companyModel := db.NewCompanyModel(mysql.MySQLConn)
	homepageModel := db.NewHomepageModel(mysql.MySQLConn)

	homepage, err := homepageModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.AdminGetHomePageResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotRegister, "用户未开启主页"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if homepage.Close {
		return &types.AdminGetHomePageResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotRegister, "用户主页处于关闭状态"),
		}, nil
	}

	company, err := companyModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		company = &db.Company{
			CompanyName: "",
		}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sex := "未知"
	if homepage.Man.Valid {
		if homepage.Man.Bool {
			sex = "男"
		} else {
			sex = "女"
		}
	} else {
		sex = "未知"
	}

	return &types.AdminGetHomePageResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetHomePageData{
			CompanyName:  company.CompanyName,
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
	}, nil
}
