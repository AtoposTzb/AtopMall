package banners

import (
	"net/http"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"

	"atopmall_web/goods_web/api"
	"atopmall_web/goods_web/forms"
	"atopmall_web/goods_web/global"
	"atopmall_web/goods_web/proto"
)

func GetBannerList(ctx *gin.Context) {
	// 获取轮播图列表
	// 限流
	e, b := sentinel.Entry("banner-list", sentinel.WithTrafficType(base.Inbound))
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
	rsp, err := global.GoodsSrvCli.Banner.BannerList(ctx.Request.Context(), &empty.Empty{})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["index"] = value.Index
		reMap["image"] = value.Image
		reMap["url"] = value.Url

		result = append(result, reMap)
	}

	ctx.JSON(http.StatusOK, result)
}

func NewBanner(ctx *gin.Context) {
	// 新建轮播图
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 限流
	e, b := sentinel.Entry("create-banner", sentinel.WithTrafficType(base.Inbound))
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

	rsp, err := global.GoodsSrvCli.Banner.CreateBanner(ctx.Request.Context(), &proto.BannerRequest{
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
		Image: bannerForm.Image,
	})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id
	response["index"] = rsp.Index
	response["url"] = rsp.Url
	response["image"] = rsp.Image

	ctx.JSON(http.StatusOK, response)
}

func UpdateBanner(ctx *gin.Context) {
	// 修改轮播图信息
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 限流
	e, b := sentinel.Entry("update-banner", sentinel.WithTrafficType(base.Inbound))
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

	_, err = global.GoodsSrvCli.Banner.UpdateBanner(ctx.Request.Context(), &proto.BannerRequest{
		Id:    int32(i),
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
	})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func DeleteBanner(ctx *gin.Context) {
	// 删除轮播图
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	// 限流
	e, b := sentinel.Entry("delete-banner", sentinel.WithTrafficType(base.Inbound))
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
	_, err = global.GoodsSrvCli.Banner.DeleteBanner(ctx.Request.Context(), &proto.BannerRequest{Id: int32(i)})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}
