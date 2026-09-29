package center

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"gitee.com/wuntsong-auth/backend/src/yundun"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateHomePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateHomePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHomePageLogic {
	return &UpdateHomePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateHomePageLogic) UpdateHomePage(req *types.UserUpdateHomePageReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.IsDelete {
		homePageModel := db.NewHomepageModel(mysql.MySQLConn)
		_, err = homePageModel.InsertWithDelete(l.ctx, &db.Homepage{
			UserId: user.Id,
			Close:  true,
		})
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		audit.NewUserAudit(user.Id, "用户删除主页信息成功")
		return &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		}, nil
	}

	manValid := false
	isMan := true
	if req.Sex == "男" {
		manValid = true
		isMan = true
	} else if req.Sex == "女" {
		manValid = true
		isMan = false
	} else {
		manValid = false
	}

	if len(req.Phone) != 0 && !utils.IsPhoneCall(req.Phone) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPhone, "电话错误"),
		}, nil
	}

	if len(req.Email) != 0 && !utils.IsEmailAddress(req.Email) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadEmail, "邮箱错误"),
		}, nil
	}

	if len(req.QQ) != 0 && !utils.IsQQ(req.QQ) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadQQ, "QQ号错误"),
		}, nil
	}

	if len(req.WeChat) != 0 && !utils.IsWeChat(req.WeChat) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadWeChat, "微信号错误"),
		}, nil
	}

	if len(req.Company) != 0 && !utils.IsValidChineseCompanyName(req.Company) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadWeChat, "公司名错误"),
		}, nil
	}

	testText := ""
	if len(req.Address) != 0 {
		testText += fmt.Sprintf("地址：%s\n", req.Address)
	}

	if len(req.Company) != 0 {
		testText += fmt.Sprintf("公司：%s\n", req.Company)
	}

	if len(req.Position) != 0 {
		testText += fmt.Sprintf("职位：%s\n", req.Position)
	}

	if len(req.Industry) != 0 {
		testText += fmt.Sprintf("行业：%s\n", req.Industry)
	}

	if len(req.Introduction) != 0 {
		testText += fmt.Sprintf("简介：%s\n", req.Introduction)
	}

	res, err := yundun.CheckText(testText)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadHomePage, errors.WarpQuick(err), "信息审核失败"),
		}, nil
	} else if !res {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadHomePage, "行业信息审核不通过"),
		}, nil
	}

	if len(req.Link) != 0 && !utils.IsHttpOrHttps(req.Link) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadLink, "错误的外部链接"),
		}, nil
	}

	homePageModel := db.NewHomepageModel(mysql.MySQLConn)
	_, err = homePageModel.InsertWithDelete(context.Background(), &db.Homepage{
		UserId: user.Id,
		Introduction: sql.NullString{
			Valid:  len(req.Introduction) != 0,
			String: req.Introduction,
		},
		Address: sql.NullString{
			Valid:  len(req.Address) != 0,
			String: req.Address,
		},
		Phone: sql.NullString{
			Valid:  len(req.Phone) != 0,
			String: req.Phone,
		},
		Email: sql.NullString{
			Valid:  len(req.Email) != 0,
			String: req.Email,
		},
		Wechat: sql.NullString{
			Valid:  len(req.WeChat) != 0,
			String: req.WeChat,
		},
		Qq: sql.NullString{
			Valid:  len(req.QQ) != 0,
			String: req.QQ,
		},
		Man: sql.NullBool{
			Valid: manValid,
			Bool:  isMan,
		},
		Link: sql.NullString{
			Valid:  len(req.Link) != 0,
			String: req.Link,
		},
		Industry: sql.NullString{
			Valid:  len(req.Industry) != 0,
			String: req.Industry,
		},
		Position: sql.NullString{
			Valid:  len(req.Position) != 0,
			String: req.Position,
		},
		Company: sql.NullString{
			Valid:  len(req.Company) != 0,
			String: req.Company,
		},
		Close: false,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户更新主页信息成功")
	sender.MessageSend(user.Id, "主页更新", "主页信息已更新！")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
