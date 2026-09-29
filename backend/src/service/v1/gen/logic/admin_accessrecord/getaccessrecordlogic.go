package admin_accessrecord

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"strconv"
	"strings"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAccessRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccessRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccessRecordLogic {
	return &GetAccessRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccessRecordLogic) GetAccessRecord(req *types.GetAccessRecord) (resp *types.GetAccessRecordResp, err error) {
	accessModel := db.NewAccessRecordModel(mysql.MySQLConn)

	idSplit := strings.Split(req.RequestID, "-")
	if len(idSplit) != 3 {
		return &types.GetAccessRecordResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "解析id错误，非三位横杠连接"),
			Data: types.GetAccessRecordData{
				Find: false,
			},
		}, nil
	}

	id, err := strconv.ParseInt(idSplit[2], 10, 64)
	if err != nil {
		return &types.GetAccessRecordResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "解析id错误，第三位非十进制数字"),
			Data: types.GetAccessRecordData{
				Find: false,
			},
		}, nil
	}

	access, err := accessModel.FindOne(l.ctx, id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.GetAccessRecordResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.GetAccessRecordData{
				Find: false,
			},
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if access.RequestIdPrefix != fmt.Sprintf("%s-%s", idSplit[0], idSplit[1]) {
		return &types.GetAccessRecordResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Success, "请求ID不匹配"),
			Data: types.GetAccessRecordData{
				Find: false,
			},
		}, nil
	}

	startAt := int64(0)
	if access.StartAt.Valid {
		startAt = access.StartAt.Time.Unix()
	}

	endAt := int64(0)
	if access.EndAt.Valid {
		endAt = access.EndAt.Time.Unix()
	}

	return &types.GetAccessRecordResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetAccessRecordData{
			Find: true,
			Record: types.AccessRecord{
				Id:                access.Id,
				RequestIdPrefix:   access.RequestIdPrefix,
				RequestId:         fmt.Sprintf("%s-%d", access.RequestIdPrefix, access.Id),
				ServerName:        access.ServerName,
				UserId:            access.UserId.Int64,
				UserUid:           access.UserUid.String,
				UserToken:         access.UserToken.String,
				RoleId:            access.RoleId.Int64,
				RoleName:          access.RoleName.String,
				RoleSign:          access.RoleSign.String,
				WebId:             access.WebId.Int64,
				WebName:           access.WebName.String,
				RequestsWebId:     access.RequestsWebId.Int64,
				RequestsWebName:   access.RequestsWebName.String,
				Ip:                access.Ip,
				GeoCode:           access.GeoCode,
				Geo:               access.Geo,
				Scheme:            access.Scheme,
				Method:            access.Method,
				Host:              access.Host,
				Path:              access.Path,
				Query:             access.Query,
				RequestsBody:      access.RequestsBody,
				ResponseBody:      access.ResponseBody.String,
				ResponseBodyError: access.ResponseBodyError.String,
				RequestsHeader:    access.RequestsHeader,
				ResponseHeader:    access.ResponseHeader.String,
				StatusCode:        access.StatusCode.Int64,
				PanicError:        access.PanicError.String,
				Message:           access.Message.String,
				UseTime:           access.UseTime.Int64,
				CreateAt:          access.CreateAt.Unix(),
				StartAt:           startAt,
				EndAt:             endAt,
			},
		},
	}, nil
}
