package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"atopmall_web/user_web/forms"
	"atopmall_web/user_web/global"
	"atopmall_web/user_web/global/responselist"
	"atopmall_web/user_web/middlewares"
	"atopmall_web/user_web/models"
	"atopmall_web/user_web/proto"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}

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

func GetUserList(ctx *gin.Context) {
	//调用user_web的user.proto接口 GetUserList 也就是远程
	//跨越问题-- 后端解决 也可以前端解决 这里采用后端解决，跨域问题如何产生？详看有道云笔记
	ip := global.ServerConfig.UserSrvInfo.Host
	port := global.ServerConfig.UserSrvInfo.Port
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

	//从上下文获取用户的id
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("用户的id为:%d", currentUser.ID)

	pnInt, _ := strconv.Atoi(ctx.DefaultQuery("pn", "0"))
	pnSizeInt, _ := strconv.Atoi(ctx.DefaultQuery("psize", "10"))
	rsp, err := userSrcClient.GetUserList(context.Background(), &proto.PageInfo{
		PageNum:  uint32(pnInt),
		PageSize: uint32(pnSizeInt),
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

func PasswordLogin(ctx *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	//同样连接用户的grpc服务，和上面重复，后面再优化
	// ip := global.ServerConfig.UserSrvInfo.Host
	// port := global.ServerConfig.UserSrvInfo.Port
	userCoun, err := grpc.NewClient(global.ServerConfig.UserSrvInfo.Host+":"+strconv.Itoa(global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接【用户服务失败】",
			"msg", err.Error(),
		)

	}
	// 关闭连接
	defer userCoun.Close()
	//生成grpc的client并调用接口
	userSrcClient := proto.NewUserClient(userCoun)

	//登录的逻辑：查询用户是否存在，如果存在，判断密码是否正确
	rsp, err := userSrcClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"mobile": "登录失败",
				})

			}

		}
		return
	} else {
		//检查密码是否正确
		if passRep, passErr := userSrcClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"mobile": "登录失败",
			})
		} else {
			if passRep.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityID: uint(rsp.Rolo),
					RegisteredClaims: jwt.RegisteredClaims{
						NotBefore: jwt.NewNumericDate(time.Now()), // //签名的生效时间
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
						Issuer:    "atopmall",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				//返回token
				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": claims.ExpiresAt.Unix(), //过期时间给前端用 ,单位是秒
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"msg": "密码错误",
				})
			}

		}
	}

}
