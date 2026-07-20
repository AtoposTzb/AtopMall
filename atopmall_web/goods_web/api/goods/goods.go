package goods

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"atopmall_web/goods_web/api"
	"atopmall_web/goods_web/forms"
	"atopmall_web/goods_web/global"
	"atopmall_web/goods_web/proto"
)

func GetGoodsList(ctx *gin.Context) {
	//查询商品列表（各种过滤）
	req := &proto.GoodsFilterRequest{}
	priceMinInt, _ := strconv.Atoi(ctx.DefaultQuery("pmin", "0"))
	req.PriceMin = int32(priceMinInt)

	priceMaxInt, _ := strconv.Atoi(ctx.DefaultQuery("pmax", "0"))
	req.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ishot", "0")
	if isHot == "1" {
		req.IsHot = true
	}
	isNew := ctx.DefaultQuery("isnew", "0")
	if isNew == "1" {
		req.IsNew = true
	}
	isTab := ctx.DefaultQuery("istab", "0")
	if isTab == "1" {
		req.IsTab = true
	}
	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	req.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	req.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	req.PagePerNums = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	req.KeyWords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	req.Brand = int32(brandIdInt)

	//grpc 请求商品的service服务
	r, err := global.GoodsSrvCli.Goods.GoodsList(context.WithValue(context.Background(), "ginContext", ctx), req)
	if err != nil {
		zap.S().Errorw("[GetGoodsList] 查询 【商品列表】失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": r.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"ctegory": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList

	ctx.JSON(http.StatusOK, reMap)
}

func NewGoods(ctx *gin.Context) {
	//新增商品
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	//通过grpc调用商品的service服务
	rsp, err := global.GoodsSrvCli.Goods.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, rsp)

}

func GetGoodsDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	idINt, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	rsp, err := global.GoodsSrvCli.Goods.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: int32(idINt),
	})
	if err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	rsps := map[string]interface{}{
		"id":          rsp.Id,
		"name":        rsp.Name,
		"goods_brief": rsp.GoodsBrief,
		"desc":        rsp.GoodsDesc,
		"ship_free":   rsp.ShipFree,
		"images":      rsp.Images,
		"desc_images": rsp.DescImages,
		"front_image": rsp.GoodsFrontImage,
		"shop_price":  rsp.ShopPrice,
		"ctegory": map[string]interface{}{
			"id":   rsp.Category.Id,
			"name": rsp.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   rsp.Brand.Id,
			"name": rsp.Brand.Name,
			"logo": rsp.Brand.Logo,
		},
		"is_hot":  rsp.IsHot,
		"is_new":  rsp.IsNew,
		"on_sale": rsp.OnSale,
	}
	ctx.JSON(http.StatusOK, rsps)
}

func DeleteGoods(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvCli.Goods.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Stocks(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	resp, err := global.InventorySrvCli.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: int32(idInt),
	})
	if err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"goodsId": resp.GoodsId,
		"num":     resp.Num,
	})
}

// 更新部分状态 GoodsStatusForm字段
func UpdateGoodsStatus(ctx *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvCli.Goods.UpdateGoodsStatus(context.Background(), &proto.GoodsStatusRequest{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

// 更新商品
func UpdateGoods(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvCli.Goods.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
