package categorybrand

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

func GetCategoryBrandList(ctx *gin.Context) {
	//根据品牌分类id获取列表
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 限流
	e, b := sentinel.Entry("category-brand-list", sentinel.WithTrafficType(base.Inbound))
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

	rsp, err := global.GoodsSrvCli.CategoryBrand.GetCategoryBrandList(ctx.Request.Context(), &proto.CategoryInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	ctx.JSON(http.StatusOK, result)
}

func CategoryBrandList(ctx *gin.Context) {
	//获取所有品牌分类列表
	// 限流
	e, b := sentinel.Entry("category-brand-all", sentinel.WithTrafficType(base.Inbound))
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
	//所有的list返回的数据结构
	/*
		{
			"total": 100,
			"data":[{},{}]
		}
	*/
	rsp, err := global.GoodsSrvCli.CategoryBrand.CategoryBrandList(ctx.Request.Context(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["category"] = map[string]interface{}{
			"id":   value.Category.Id,
			"name": value.Category.Name,
		}
		reMap["brand"] = map[string]interface{}{
			"id":   value.Brand.Id,
			"name": value.Brand.Name,
			"logo": value.Brand.Logo,
		}

		result = append(result, reMap)
	}

	reMap["data"] = result
	ctx.JSON(http.StatusOK, reMap)
}

func NewCategoryBrand(ctx *gin.Context) {
	//创建品牌分类
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 限流
	e, b := sentinel.Entry("create-category-brand", sentinel.WithTrafficType(base.Inbound))
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

	rsp, err := global.GoodsSrvCli.CategoryBrand.CreateCategoryBrand(ctx.Request.Context(), &proto.CategoryBrandRequest{
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id

	ctx.JSON(http.StatusOK, response)
}

func UpdateCategoryBrand(ctx *gin.Context) {
	//更新品牌分类
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
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
	e, b := sentinel.Entry("update-category-brand", sentinel.WithTrafficType(base.Inbound))
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

	_, err = global.GoodsSrvCli.CategoryBrand.UpdateCategoryBrand(ctx.Request.Context(), &proto.CategoryBrandRequest{
		Id:         int32(i),
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func DeleteCategoryBrand(ctx *gin.Context) {
	//删除品牌分类
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	// 限流
	e, b := sentinel.Entry("delete-category-brand", sentinel.WithTrafficType(base.Inbound))
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
	_, err = global.GoodsSrvCli.CategoryBrand.DeleteCategoryBrand(ctx.Request.Context(), &proto.CategoryBrandRequest{Id: int32(i)})
	if err != nil {
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "删除成功")
}
