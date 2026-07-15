package pay

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"atopmall_web/order_web/global"
	"atopmall_web/order_web/proto"
)

type OrderInfo struct {
	OrderSn string
	Total   float32
}

// 实例化支付宝
func AlipayClient(ctx *gin.Context) (client *alipay.Client) {
	alipayInfo := global.ServerConfig.AlipayInfo
	client, err := alipay.New(alipayInfo.AppID, alipayInfo.PrivateKey, alipayInfo.IsProduction)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(alipayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	return client
}

// 获取支付宝支付链接
func AlipayUrl(ctx *gin.Context, orderInfo OrderInfo) string {
	alipayInfo := global.ServerConfig.AlipayInfo
	client := AlipayClient(ctx)
	var p = alipay.TradePagePay{}
	p.NotifyURL = alipayInfo.NotifyURL
	p.ReturnURL = alipayInfo.ReturnURL
	p.Subject = "atopmall订单-" + orderInfo.OrderSn
	p.OutTradeNo = orderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(orderInfo.Total), 'f', 2, 64)
	p.ProductCode = alipayInfo.ProductCode //网页支付

	alipay_url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付宝支付链接失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return ""
	}
	return alipay_url.String()
}

// 支付宝回调通知
func Notify(ctx *gin.Context) {
	client := AlipayClient(ctx)
	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	_, err = global.OrderSrvCli.Order.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.String(http.StatusOK, "success")

}
