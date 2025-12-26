package router

import (
	"net/http"

	"cloud-disk/core/internal/handler"
	"cloud-disk/core/svc"

	"github.com/gin-gonic/gin"
)

// Register wires all HTTP routes.
func Register(r *gin.Engine, serviceName string, svcCtx *svc.ServiceContext) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"status":  "ok",
		})
	})

	r.POST("/user/login", handler.UserLoginHandler(svcCtx))
	r.POST("/user/logout", handler.UserLogoutHandler(svcCtx))
	r.POST("/user/detail", handler.UserDetailHandler(svcCtx))
	r.POST("/mail/code/send/register", handler.MailCodeSendRegisterHandler(svcCtx))
	r.POST("/user/register", handler.UserRegisterHandler(svcCtx))
	r.GET("/share/basic/detail", handler.ShareBasicDetailHandler(svcCtx))

	auth := r.Group("/")
	auth.Use(svcCtx.Auth)
	{
		auth.POST("/file/upload", handler.FileUploadHandler(svcCtx))
		auth.GET("/file/download", handler.FileDownloadHandler(svcCtx))
		auth.POST("/user/repository/save", handler.UserRepositorySaveHandler(svcCtx))
		auth.POST("/user/file/list", handler.UserFileListHandler(svcCtx))
		auth.POST("/user/folder/list", handler.UserFolderListHandler(svcCtx))
		auth.POST("/user/file/name/update", handler.UserFileNameUpdateHandler(svcCtx))
		auth.POST("/user/folder/create", handler.UserFolderCreateHandler(svcCtx))
		auth.DELETE("/user/file/delete", handler.UserFileDeleteHandler(svcCtx))
		auth.PUT("/user/file/move", handler.UserFileMoveHandler(svcCtx))
		auth.POST("/share/basic/create", handler.ShareBasicCreateHandler(svcCtx))
		auth.POST("/share/basic/save", handler.ShareBasicSaveHandler(svcCtx))
		auth.POST("/refresh/authorization", handler.RefreshAuthorizationHandler(svcCtx))
		auth.POST("/file/upload/prepare", handler.FileUploadPrepareHandler(svcCtx))
		auth.POST("/file/upload/chunk", handler.FileUploadChunkHandler(svcCtx))
		auth.POST("/file/upload/chunk/complete", handler.FileUploadChunkCompleteHandler(svcCtx))
		
		// Friend system endpoints
		auth.POST("/friend/request/send", handler.FriendRequestSendHandler(svcCtx))
		auth.POST("/friend/request/list", handler.FriendRequestListHandler(svcCtx))
		auth.POST("/friend/request/respond", handler.FriendRequestRespondHandler(svcCtx))
		auth.POST("/friend/list", handler.FriendListHandler(svcCtx))
		auth.POST("/friend/share/create", handler.FriendShareCreateHandler(svcCtx))
		auth.POST("/friend/share/list", handler.FriendShareListHandler(svcCtx))
		auth.POST("/friend/share/mark-read", handler.FriendShareMarkReadHandler(svcCtx))
	}
}
