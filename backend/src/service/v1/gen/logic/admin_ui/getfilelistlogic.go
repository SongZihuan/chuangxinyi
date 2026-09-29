package admin_ui

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileListLogic {
	return &GetFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileListLogic) GetFileList(req *types.GetUIFileList) (resp *types.FileListResp, err error) {
	ossModel := db.NewOssFileModel(mysql.MySQLConn)
	fileList, err := ossModel.GetList(l.ctx, req.Name, req.MediaType, req.Page, req.PageSize)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := ossModel.GetCount(l.ctx, req.Name, req.MediaType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.File, 0, len(fileList))
	for _, d := range fileList {
		respList = append(respList, types.File{
			Fid:       d.Fid,
			MediaType: d.MediaType,
		})
	}

	return &types.FileListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.FileListData{
			Count: count,
			File:  respList,
		},
	}, nil
}
