package brands

import (
	"net/http"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"

	"atopmall_web/goods_web/api"
	"atopmall_web/goods_web/forms"
	"atopmall_web/goods_web/global"
	"atopmall_web/goods_web/proto"
)

func BrandList(ctx *gin.Context) {
	//获取所有品牌列表
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	// 限流
	e, b := sentinel.Entry("brand-list", sentinel.WithTrafficType(base.Inbound))
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

	rsp, err := global.GoodsSrvCli.Brand.BrandList(ctx.Request.Context(), &proto.BrandFilterRequest{
		Pages:       int32(pnInt),
		PagePerNums: int32(pSizeInt),
	})

	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	reMap := make(map[string]interface{})
	reMap["total"] = rsp.Total
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func NewBrand(ctx *gin.Context) {
	//创建品牌
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 限流
	e, b := sentinel.Entry("create-brand", sentinel.WithTrafficType(base.Inbound))
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

	rsp, err := global.GoodsSrvCli.Brand.CreateBrand(ctx.Request.Context(), &proto.BrandRequest{
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["logo"] = rsp.Logo

	ctx.JSON(http.StatusOK, request)
}

func DeleteBrand(ctx *gin.Context) {
	//删除品牌
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	// 限流
	e, b := sentinel.Entry("delete-brand", sentinel.WithTrafficType(base.Inbound))
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
	_, err = global.GoodsSrvCli.Brand.DeleteBrand(ctx.Request.Context(), &proto.BrandRequest{Id: int32(i)})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func UpdateBrand(ctx *gin.Context) {
	//更新品牌
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
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
	e, b := sentinel.Entry("update-brand", sentinel.WithTrafficType(base.Inbound))
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

	_, err = global.GoodsSrvCli.Brand.UpdateBrand(ctx.Request.Context(), &proto.BrandRequest{
		Id:   int32(i),
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
