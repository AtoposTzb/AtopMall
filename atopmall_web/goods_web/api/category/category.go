package category

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	empty "github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"

	"atopmall_web/goods_web/api"
	"atopmall_web/goods_web/forms"
	"atopmall_web/goods_web/global"
	"atopmall_web/goods_web/proto"
)

func GetCategoryList(ctx *gin.Context) {
	//获取所有分类列表
	r, err := global.GoodsSrvCli.Category.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	data := make([]interface{}, 0)
	err = json.Unmarshal([]byte(r.JsonData), &data)
	if err != nil {
		zap.S().Errorw("[List] 查询 【分类列表】失败： ", err.Error())
	}

	ctx.JSON(http.StatusOK, data)
}

func GetCategoryDetail(ctx *gin.Context) {
	//获取分类详情
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	reMap := make(map[string]interface{})
	subCategorys := make([]interface{}, 0)
	if r, err := global.GoodsSrvCli.Category.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: int32(i),
	}); err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	} else {
		//写文档 特别是数据多的时候很慢， 先开发后写文档
		for _, value := range r.SubCategorys {
			subCategorys = append(subCategorys, map[string]interface{}{
				"id":              value.Id,
				"name":            value.Name,
				"level":           value.Level,
				"parent_category": value.ParentCategory,
				"is_tab":          value.IsTab,
			})
		}
		reMap["id"] = r.Info.Id
		reMap["name"] = r.Info.Name
		reMap["level"] = r.Info.Level
		reMap["parent_category"] = r.Info.ParentCategory
		reMap["is_tab"] = r.Info.IsTab
		reMap["sub_categorys"] = subCategorys

		ctx.JSON(http.StatusOK, reMap)
	}
}

func NewCategory(ctx *gin.Context) {
	//创建分类
	categoryForm := forms.CategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvCli.Category.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,
		IsTab:          *categoryForm.IsTab,
		Level:          categoryForm.Level,
		ParentCategory: categoryForm.ParentCategory,
	})
	if err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["parent"] = rsp.ParentCategory
	request["level"] = rsp.Level
	request["is_tab"] = rsp.IsTab

	ctx.JSON(http.StatusOK, request)
}

func DeleteCategory(ctx *gin.Context) {
	//删除分类
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	//1. 先查询出该分类写的所有子分类
	//2. 将所有的分类全部逻辑删除
	//3. 将该分类下的所有的商品逻辑删除
	_, err = global.GoodsSrvCli.Category.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func UpdateCategory(ctx *gin.Context) {
	//更新分类
	categoryForm := forms.UpdateCategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   int32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}
	_, err = global.GoodsSrvCli.Category.UpdateCategory(context.Background(), request)
	if err != nil {
		api.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
