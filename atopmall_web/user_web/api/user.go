package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
					"msg": s.Message(),
				})
			}
		}
	}

}

func GetUserList(ctx *gin.Context) {
	//从上下文获取用户的id
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("用户的id为:%d", currentUser.ID)

	pnInt, _ := strconv.Atoi(ctx.DefaultQuery("pn", "0"))
	pnSizeInt, _ := strconv.Atoi(ctx.DefaultQuery("psize", "10"))
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
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
			Email:    value.Email,
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
	ctx.JSON(http.StatusOK, result) //返回用户列表给前端
}

func PasswordLogin(ctx *gin.Context) {
	//表单验证,获取前端传递的手机号、密码、验证码、验证码id然后验证
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	//同样连接用户的grpc服务，和上面重复，后面再优化
	// ip := global.ServerConfig.UserSrvInfo.Host
	// port := global.ServerConfig.UserSrvInfo.Port

	//检查验证码是否正确
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	//登录的逻辑：查询用户是否存在，如果存在，判断密码是否正确
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
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
		if passRep, passErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
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
				expireDuration := time.Duration(global.ServerConfig.JWTInfo.LoginExpireHour) * time.Hour //nacos配置登录token过期时间
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityID: uint(rsp.Role),
					RegisteredClaims: jwt.RegisteredClaims{
						NotBefore: jwt.NewNumericDate(time.Now()),                     //签名的生效时间
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)), //过期时间
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
				//返回token ,注意敏感信息不要添加到token中
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

func Register(ctx *gin.Context) {
	//表单验证
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	//通过唯一的手机号查询用户是否存在
	if user, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: registerForm.Mobile,
	}); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": fmt.Sprintf("用户%s已存在,请直接登录", user.Mobile),
		})
		return
	} else {
		//用户不存在，继续注册,通过邮箱验证码注册用户
		//邮箱验证码检查
		value, err := global.RDB.Get(context.Background(), registerForm.Email).Result()
		if err == redis.Nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"email": "邮箱验证码错误",
			})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"data": "验证码获取失败",
			})
			return
		}
		if value != registerForm.Code {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"data": "验证码错误",
			})
			return
		}
		//验证码正确，注册用户
		user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: registerForm.Mobile,
			Password: registerForm.Password,
			Mobile:   registerForm.Mobile,
			Email:    registerForm.Email,
		})
		if err != nil {
			zap.S().Errorw("[Register] 注册用户失败")
			HandleGrpcErrorToHttpError(err, ctx)
			return
		}
		//生成token
		j := middlewares.NewJWT()
		expireDuration := time.Duration(global.ServerConfig.JWTInfo.LoginExpireHour) * time.Hour //nacos配置登录token过期时间
		claims := models.CustomClaims{
			ID:          uint(user.Id),
			NickName:    user.NickName,
			AuthorityID: uint(user.Role),
			RegisteredClaims: jwt.RegisteredClaims{
				NotBefore: jwt.NewNumericDate(time.Now()),                     //签名的生效时间
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)), //过期时间
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
		//返回token ,注意敏感信息不要添加到token中
		ctx.JSON(http.StatusOK, gin.H{
			"id":         user.Id,
			"nick_name":  user.NickName,
			"token":      token,
			"expired_at": claims.ExpiresAt.Unix(), //过期时间给前端用 ,单位是秒
		})
	}

}
