package message

import (
	"net/http"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"atopmall_web/userop_web/api"
	"atopmall_web/userop_web/forms"
	"atopmall_web/userop_web/global"
	"atopmall_web/userop_web/models"
	"atopmall_web/userop_web/proto"
)

func GetMessageList(ctx *gin.Context) {
	// 获取留言列表
	request := &proto.MessageRequest{}

	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityID == 1 {
		request.UserId = int32(userId.(uint))
	}

	// 限流
	e, b := sentinel.Entry("message-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求频率过快,请稍后重试",
		})
		return
	}
	defer func() {
		if e != nil {
			e.Exit()
		}
	}()

	rsp, err := global.MessageSrvCli.MessageList(ctx.Request.Context(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["user_id"] = value.UserId
		reMap["type"] = value.MessageType
		reMap["subject"] = value.Subject
		reMap["message"] = value.Message
		reMap["file"] = value.File

		result = append(result, reMap)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  result,
	})
}

func NewMessage(ctx *gin.Context) {
	// 新建留言
	userId, _ := ctx.Get("userId")

	messageForm := forms.MessageForm{}
	if err := ctx.ShouldBindJSON(&messageForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 限流
	e, b := sentinel.Entry("message-create", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求频率过快,请稍后重试",
		})
		return
	}
	defer func() {
		if e != nil {
			e.Exit()
		}
	}()

	rsp, err := global.MessageSrvCli.CreateMessage(ctx.Request.Context(), &proto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
