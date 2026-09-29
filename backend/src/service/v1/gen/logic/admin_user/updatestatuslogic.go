package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatusLogic {
	return &UpdateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStatusLogic) UpdateStatus(req *types.AdminUpdateStatusReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if srcUser.IsAdmin {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.AdminCanNotChangeStatus, "根管理员不能修改状态"),
		}, nil
	}

	status, ok := db.UserStatusStartMap[req.Status]
	if !ok {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadUserStatus, "错误的用户状态"),
		}, nil
	}

	oldStatus := srcUser.Status
	/*
		保存信息和保存信息状态之间可以互转
		不保存信息和不保存信息状态之间可以互转
		保存信息状态可以转到不保存信息状态
		其他不被允许
	*/
	if !db.IsKeepInfoStatus(oldStatus) && db.IsKeepInfoStatus(status) { // 之前是不保存信息的状态，现在是保存信息的状态
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadUserStatus, "用户状态转换不允许"),
		}, nil
	} else if db.IsKeepInfoStatus(oldStatus) && !db.IsKeepInfoStatus(status) { // 之前是保存info的status，现在不是
		err = utils.DeleteUser(srcUser, status, 1000)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		userModel := db.NewUserModel(mysql.MySQLConn)
		srcUser.Status = status
		err = userModel.Update(l.ctx, srcUser) // 专门更新status
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		if db.IsBannedStatus(status) {
			_ = jwt.DeleteAllUserToken(l.ctx, srcUser.Uid, "")
			_ = jwt.DeleteAllUserWebsiteToken(l.ctx, srcUser.Uid)
			_ = jwt.DeleteAllUserSonToken(l.ctx, srcUser.Uid)
			_ = jwt.DeleteAllFatherUserToken(l.ctx, srcUser.Uid)
		}
	}

	audit.NewAdminAudit(user.Id, "管理员更新用户（%s）状态成功", srcUser.Uid)

	if (oldStatus == db.UserStatus_Register || oldStatus == db.UserStatus_Normal) && (status == db.UserStatus_Delete || status == db.UserStatus_Banned) {
		sender.PhoneSendDelete(srcUser.Id)
		sender.EmailSendDelete(srcUser.Id)
		sender.WxrobotSendDelete(srcUser.Id)
		sender.FuwuhaoSendDeleteUser(srcUser.Id)
		audit.NewUserAudit(srcUser.Id, "用户被注销或封禁")
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
