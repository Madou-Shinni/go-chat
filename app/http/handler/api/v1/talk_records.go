package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/service"
	"net/http"
)

type TalkRecords struct {
	service            *service.TalkRecordsService
	groupMemberService *service.GroupMemberService
	fileSystem         *filesystem.Filesystem
}

func NewTalkRecordsHandler(service *service.TalkRecordsService, groupMemberService *service.GroupMemberService, fileSystem *filesystem.Filesystem) *TalkRecords {
	return &TalkRecords{
		service:            service,
		groupMemberService: groupMemberService,
		fileSystem:         fileSystem,
	}
}

// GetRecords 获取会话记录
func (c *TalkRecords) GetRecords(ctx *gin.Context) {
	params := &request.TalkRecordsRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	records, err := c.service.GetTalkRecords(ctx, &service.QueryTalkRecordsOpts{
		TalkType:   params.TalkType,
		UserId:     auth.GetAuthUserID(ctx),
		ReceiverId: params.ReceiverId,
		RecordId:   params.RecordId,
		Limit:      params.Limit,
	})

	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	rid := 0
	if length := len(records); length > 0 {
		rid = records[length-1].Id
	}

	response.Success(ctx, gin.H{
		"limit":     params.Limit,
		"record_id": rid,
		"rows":      records,
	})
}

// SearchHistoryRecords 查询下会话记录
func (c *TalkRecords) SearchHistoryRecords(ctx *gin.Context) {
	c.GetRecords(ctx)
}

// GetForwardRecords 获取转发记录
func (c *TalkRecords) GetForwardRecords(ctx *gin.Context) {
	params := &request.TalkForwardRecordsRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	records, err := c.service.GetForwardRecords(ctx.Request.Context(), auth.GetAuthUserID(ctx), int64(params.RecordId))
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"rows": records,
	})
}

// Download 聊天文件下载
func (c *TalkRecords) Download(ctx *gin.Context) {
	params := &request.DownloadChatFileRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	resp, err := c.service.Dao().FindFileRecord(ctx.Request.Context(), params.RecordId)
	if err != nil {
		return
	}

	uid := auth.GetAuthUserID(ctx)
	if uid != resp.Record.UserId {
		if resp.Record.TalkType == entity.PrivateChat {
			if resp.Record.ReceiverId != uid {
				response.Unauthorized(ctx, "无访问权限！")
				return
			}
		} else {
			if !c.groupMemberService.Dao().IsMember(resp.Record.ReceiverId, uid, false) {
				response.Unauthorized(ctx, "无访问权限！")
				return
			}
		}
	}

	switch resp.FileInfo.Drive {
	case entity.FileDriveLocal:
		ctx.FileAttachment(c.fileSystem.Local.Path(resp.FileInfo.Path), resp.FileInfo.OriginalName)
	case entity.FileDriveCos:
		ctx.Redirect(http.StatusFound, c.fileSystem.Cos.PrivateUrl(resp.FileInfo.Path, 60))
	}
}
