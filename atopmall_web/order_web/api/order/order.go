package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"atopmall_web/order_web/api"
	"atopmall_web/order_web/api/pay"
	"atopmall_web/order_web/forms"
	"atopmall_web/order_web/global"
	"atopmall_web/order_web/models"
	"atopmall_web/order_web/proto"
)

func OrderList(ctx *gin.Context) {
	//获取订单列表
	user_id, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	request := proto.OrderFilterRequest{}
	if currentUser.AuthorityID == 1 {
		request.UserId = int32(user_id.(uint))
	}
	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	orderRspList, err := global.OrderSrvCli.Order.OrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("[OrderList] 调用【订单列表查询】接口失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	orderList := make([]interface{}, 0)
	for _, item := range orderRspList.Data {
		tmpMap := map[string]interface{}{}

		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["id"] = item.Id
		tmpMap["add_time"] = item.AddTime

		orderList = append(orderList, tmpMap)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":  200,
		"total": orderRspList.Total,
		"data":  orderList,
	})
}

func OrderDetail(ctx *gin.Context) {
	//获取订单详情
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	//如果是普通用户，只返回自己的订单(简单的权限校验)
	request := proto.OrderRequest{
		Id: int32(i),
	}
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityID == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.OrderSrvCli.Order.OrderDetail(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("[OrderDetail] 调用【订单详情查询】接口失败")
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	reMap := gin.H{}
	reMap["id"] = rsp.OrderInfo.Id
	reMap["status"] = rsp.OrderInfo.Status
	reMap["user"] = rsp.OrderInfo.UserId
	reMap["post"] = rsp.OrderInfo.Post
	reMap["total"] = rsp.OrderInfo.Total
	reMap["address"] = rsp.OrderInfo.Address
	reMap["name"] = rsp.OrderInfo.Name
	reMap["mobile"] = rsp.OrderInfo.Mobile
	reMap["pay_type"] = rsp.OrderInfo.PayType
	reMap["order_sn"] = rsp.OrderInfo.OrderSn

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		tmpMap := gin.H{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		}

		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList
	//添加支付宝支付链接
	url := pay.AlipayUrl(ctx, pay.OrderInfo{
		OrderSn: rsp.OrderInfo.OrderSn,
		Total:   rsp.OrderInfo.Total,
	})
	reMap["alipay_url"] = url

	ctx.JSON(http.StatusOK, reMap)
}

func OrderCreate(ctx *gin.Context) {
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	user_id, _ := ctx.Get("userId")
	orderRsp, err := global.OrderSrvCli.Order.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  int32(user_id.(uint)),
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Address: orderForm.Address,
		Post:    orderForm.Post,
	})
	if err != nil {
		zap.S().Errorw("[OrderCreate] 调用【订单创建】接口失败", "err", err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	//生成支付宝支付链接
	url := pay.AlipayUrl(ctx, pay.OrderInfo{
		OrderSn: orderRsp.OrderSn,
		Total:   orderRsp.Total,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"id":         orderRsp.Id,
		"alipay_url": url,
	})
}
