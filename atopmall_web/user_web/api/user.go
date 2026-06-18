package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"atopmall_web/user_web/global/responselist"
	"atopmall_web/user_web/proto"
)

func HandleGrpcErrorToHttpError(err error, c *gin.Context) {
	//将grpc的code转换为http的code(状态码)发送给前端
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{ //404
					"msg": s.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{ //500
					"msg": "服务器内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{ //400
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusServiceUnavailable, gin.H{ //503
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{ //500
					"msg": s.Code(),
				})
			}
		}
	}

}

func GerUserList(ctx *gin.Context) {
	//调用user_web的user.proto接口 GetUserList 也就是远程
	ip := "127.0.0.1"
	port := 50051
	userCoun, err := grpc.NewClient(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接【用户服务失败】",
			"msg", err.Error(),
		)

	}
	// 关闭连接
	defer userCoun.Close()
	//生成grpc的client并调用接口
	userSrcClient := proto.NewUserClient(userCoun)
	rsp, err := userSrcClient.GetUserList(context.Background(), &proto.PageInfo{
		PageNum:  1,
		PageSize: 10,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】 失败")
		HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {

		user := responselist.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			BirthDay: time.Unix(int64(value.BirthDay), 0).Format("2006-01-02"),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		result = append(result, user)
		// data := make(map[string]interface{})
		// data["id"] = value.Id
		// data["name"] = value.NickName
		// data["birthday"] = value.BirthDay
		// data["gender"] = value.Gender
		// data["mobile"] = value.Mobile
		// result = append(result, data)
	}
	ctx.JSON(http.StatusOK, result)
}
