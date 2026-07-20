package shoppingcart

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"atopmall_web/order_web/api"
	"atopmall_web/order_web/forms"
	"atopmall_web/order_web/global"
	"atopmall_web/order_web/proto"
)

func ShoppingCartList(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	//调用订单服务的购物车接口列表
	resp, err := global.OrderSrvCli.ShoppingCart.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("[ShoppingCartList] 调用【购物车列表查询】接口失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	//返回的是 total 和 data[]
	ids := make([]int32, 0)
	for _, item := range resp.Data {
		ids = append(ids, item.GoodsId)
	}
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}
	//调用商品服务的商品列表接口
	goodsResp, err := global.GoodsSrvCli.Goods.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[ShoppingCartList] 调用【商品列表查询】接口失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	//融合购物车列表和商品列表
	goodsList := make([]interface{}, 0)
	for _, item := range resp.Data {
		for _, goods := range goodsResp.Data {
			if item.GoodsId == goods.Id {
				tmpMap := map[string]interface{}{}
				tmpMap["id"] = item.Id
				tmpMap["goods_id"] = item.GoodsId
				tmpMap["goods_name"] = goods.Name
				tmpMap["goods_price"] = goods.ShopPrice
				tmpMap["goods_image"] = goods.GoodsFrontImage
				tmpMap["nums"] = item.Nums
				tmpMap["checked"] = item.Checked

				goodsList = append(goodsList, tmpMap)
			}
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  goodsList,
	})

}

func NewShoppingCart(ctx *gin.Context) {
	//添加商品到购物车
	itemForm := forms.ShoppingCartIntemForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	//先判断商品是否存在
	_, err := global.GoodsSrvCli.Goods.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[NewShoppingCart] 调用【商品详情查询】接口失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	//检查商品库存是否足够
	invResp, err := global.InventorySrvCli.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[NewShoppingCart] 调用【商品库存详情查询】接口失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	if invResp.Num < itemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "商品库存不足",
		})
		return
	}
	userId, _ := ctx.Get("userId")
	orederRsp, err := global.OrderSrvCli.ShoppingCart.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: itemForm.GoodsId,
		Nums:    itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("[NewShoppingCart] 调用【购物车添加】接口失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": orederRsp.Id,
	})

}

func DeleteShoppingCart(ctx *gin.Context) {
	//删除购物车
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "购物车id格式错误",
		})
		return
	}
	userId, _ := ctx.Get("userId")
	_, err = global.OrderSrvCli.ShoppingCart.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("[DeleteShoppingCart] 调用【购物车删除】接口失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

func UpdateShoppingCart(ctx *gin.Context) {
	//更新购物车
	itemForm := forms.ShoppingCartUpdateForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	userId, _ := ctx.Get("userId")
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "购物车id格式错误",
		})
		return
	}
	request := &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
		Checked: false,
	}
	if itemForm.Check != nil {
		request.Checked = *itemForm.Check
	}

	_, err = global.OrderSrvCli.ShoppingCart.UpdateCartItem(context.Background(), request)
	if err != nil {
		zap.S().Errorw("[UpdateShoppingCart] 调用【购物车更新】接口失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}
