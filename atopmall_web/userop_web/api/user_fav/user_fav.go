package userfav

import (
	"net/http"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"atopmall_web/userop_web/api"
	"atopmall_web/userop_web/forms"
	"atopmall_web/userop_web/global"
	"atopmall_web/userop_web/proto"
)

func GetUserFavList(ctx *gin.Context) {
	// 获取收藏列表
	userId, _ := ctx.Get("userId")
	// 限流
	e, b := sentinel.Entry("fav-list", sentinel.WithTrafficType(base.Inbound))
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
	userFavRsp, err := global.UserFavSrvCli.GetFavList(ctx.Request.Context(), &proto.UserFavRequest{
		UserId: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("获取收藏列表失败")
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range userFavRsp.Data {
		ids = append(ids, item.GoodsId)
	}

	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	//请求商品服务
	goods, err := global.GoodsSrvCli.Goods.BatchGetGoods(ctx.Request.Context(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[Goods] 批量查询【商品列表】失败")
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	goodsList := make([]interface{}, 0)
	for _, item := range userFavRsp.Data {
		data := gin.H{
			"id": item.GoodsId,
		}

		for _, good := range goods.Data {
			if item.GoodsId == good.Id {
				data["name"] = good.Name
				data["shop_price"] = good.ShopPrice
			}
		}

		goodsList = append(goodsList, data)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total": userFavRsp.Total,
		"data":  goodsList,
	})
}

func NewUserFav(ctx *gin.Context) {
	// 新建收藏记录
	userFavForm := forms.UserFavForm{}
	if err := ctx.ShouldBindJSON(&userFavForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 限流
	e, b := sentinel.Entry("fav-create", sentinel.WithTrafficType(base.Inbound))
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

	// 校验商品是否存在
	_, err := global.GoodsSrvCli.Goods.GetGoodsDetail(ctx.Request.Context(), &proto.GoodInfoRequest{
		Id: userFavForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("查询商品详情失败")
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	// 校验用户是否存在
	userId, _ := ctx.Get("userId")
	_, err = global.UserFavSrvCli.AddUserFav(ctx.Request.Context(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: userFavForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("添加收藏记录失败")
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func DeleteUserFav(ctx *gin.Context) {
	// 删除收藏记录
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	userId, _ := ctx.Get("userId")
	// 限流
	e, b := sentinel.Entry("fav-delete", sentinel.WithTrafficType(base.Inbound))
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
	_, err = global.UserFavSrvCli.DeleteUserFav(ctx.Request.Context(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除收藏记录失败")
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func GetUserFavDetail(ctx *gin.Context) {
	// 查询收藏状态
	goodsId := ctx.Param("id")
	goodsIdInt, err := strconv.ParseInt(goodsId, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	userId, _ := ctx.Get("userId")
	// 限流
	e, b := sentinel.Entry("fav-detail", sentinel.WithTrafficType(base.Inbound))
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
	_, err = global.UserFavSrvCli.GetUserFavDetail(ctx.Request.Context(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(goodsIdInt),
	})
	if err != nil {
		zap.S().Errorw("查询收藏状态失败")
		sentinel.TraceError(e, err)
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
