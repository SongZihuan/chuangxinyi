package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/wechat"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
)

type UpdateWeChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWeChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWeChatLogic {
	return &UpdateWeChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWeChatLogic) UpdateWeChat(req *types.UserUpdateWeChatReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.IsDelete || len(req.WeChatToken) == 0 {
		err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			wechatModel := db.NewWechatModelWithSession(session)
			_, err := wechatModel.InsertWithDelete(l.ctx, &db.Wechat{
				UserId: user.Id,
				OpenId: sql.NullString{
					Valid: false,
				},
			})
			return err
		})
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		wechatData, err := jwt.ParserWeChatToken(req.WeChatToken)
		if err != nil {
			return nil, respmsg.JWTError.WarpQuick(err)
		}

		if len(wechatData.UnionID) == 0 || len(wechatData.OpenID) == 0 {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.BadWechatInfo, errors.WarpQuick(err), "获取微信用户信息失败"),
			}, nil
		}

		userInfo, err := wechat.GetUserInfo(wechatData.AccessToken, wechatData.OpenID)
		if err != nil {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.BadWechatInfo, errors.WarpQuick(err), "获取微信用户信息失败"),
			}, nil
		} else {
			mysqlErr := mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
				wechatModel := db.NewWechatModelWithSession(session)

				wxOpenID := ""
				unionID := ""
				fuwuhaoOpenID := ""

				w, err := wechatModel.FindByUserID(l.ctx, user.Id)
				if errors.Is(err, db.ErrNotFound) {
					wxOpenID = ""
					unionID = wechatData.UnionID
					fuwuhaoOpenID = ""
				} else if err == nil {
					if w.UnionId.Valid && wechatData.UnionID == w.UnionId.String {
						wxOpenID = w.OpenId.String
						unionID = w.UnionId.String
						fuwuhaoOpenID = w.Fuwuhao.String
					} else {
						wxOpenID = ""
						unionID = wechatData.UnionID
						fuwuhaoOpenID = ""
					}
				} else {
					return err
				}

				if wechatData.IsFuwuhao {
					fuwuhaoOpenID = wechatData.OpenID
				} else {
					wxOpenID = wechatData.OpenID
				}

				_, err = wechatModel.InsertWithDelete(l.ctx, &db.Wechat{
					UserId: user.Id,
					OpenId: sql.NullString{
						Valid:  len(wxOpenID) != 0,
						String: wxOpenID,
					},
					UnionId: sql.NullString{
						Valid:  len(unionID) != 0,
						String: unionID,
					},
					Nickname: sql.NullString{
						Valid:  len(userInfo.Nickname) != 0,
						String: userInfo.Nickname,
					},
					Headimgurl: sql.NullString{
						Valid:  len(userInfo.Headimgurl) != 0,
						String: userInfo.Headimgurl,
					},
					Fuwuhao: sql.NullString{
						Valid:  len(fuwuhaoOpenID) != 0,
						String: fuwuhaoOpenID,
					},
				})
				return err
			})
			if mysqlErr != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(mysqlErr)
			}
		}
	}

	sender.MessageSendChange(user.Id, "微信")
	sender.WxrobotSendChange(user.Id, "微信")
	audit.NewUserAudit(user.Id, "用户更新微信绑定成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
