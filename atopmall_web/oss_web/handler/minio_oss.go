package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"atopmall_web/oss_web/global"
)

// 统一返回结构体，对齐原阿里云接口返回格式
type OssTokenResp struct {
	UploadUrl string `json:"upload_url"` // 前端PUT直传签名地址
	Url       string `json:"url"`        // 存入数据库的图片访问地址
}

// Token 对应路由 GET /oss/token
func Token(c *gin.Context) {
	conf := global.ServerConfig.MinIOInfo
	cli := global.MinioCli

	// 生成唯一文件路径，防止重名覆盖
	objPath := "goods/" + uuid.NewString() + ".jpg"
	expireTime := time.Duration(3600) * time.Second

	// 生成前端直传PUT预签名URL
	putUrlObj, err := cli.PresignedPutObject(context.Background(), conf.BucketName, objPath, expireTime)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "获取上传凭证失败",
			"err":  err.Error(),
		})
		return
	}
	// 拼接完整访问URL
	fullImgUrl := conf.PublicPrefix + objPath

	c.JSON(200, gin.H{
		"code": 200,
		"data": OssTokenResp{
			UploadUrl: putUrlObj.String(),
			Url:       fullImgUrl,
		},
	})
}

// HandlerRequest 阿里云专属回调接口，MinIO无服务端回调，保留空兼容（前端直传后主动提交url，此接口可注释删除）
func HandlerRequest(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "MinIO无OSS服务端回调，前端上传成功后主动携带url调用商品保存接口",
	})
}

// CleanupOrphanFiles 清理孤儿文件（上传后超过指定时间未被业务服务认领的文件）
// 前端上传文件到 MinIO 后，如果业务服务没有将 URL 存入数据库，该文件即为孤儿文件
// 此接口扫描 MinIO 中超过过期时间的文件并删除
func CleanupOrphanFiles(c *gin.Context) {
	conf := global.ServerConfig.MinIOInfo
	cli := global.MinioCli
	ctx := context.Background()

	// 默认清理 24 小时前的孤儿文件，可通过 query 参数 hours 自定义
	hours := 24
	if h := c.Query("hours"); h != "" {
		fmt.Sscanf(h, "%d", &hours)
		if hours < 1 || hours > 720 {
			c.JSON(400, gin.H{"code": 400, "msg": "hours 参数范围为 1-720"})
			return
		}
	}
	cutoff := time.Now().Add(-time.Duration(hours) * time.Hour)

	deletedCount := 0
	var deletedFiles []string

	// 列出桶中所有对象
	objectCh := cli.ListObjects(ctx, conf.BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for obj := range objectCh {
		if obj.Err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "列出对象失败", "err": obj.Err.Error()})
			return
		}
		// 超过过期时间的文件视为孤儿文件，执行删除
		if obj.LastModified.Before(cutoff) {
			err := cli.RemoveObject(ctx, conf.BucketName, obj.Key, minio.RemoveObjectOptions{})
			if err != nil {
				continue // 删除失败跳过，继续处理下一个
			}
			deletedCount++
			deletedFiles = append(deletedFiles, obj.Key)
		}
	}

	c.JSON(200, gin.H{
		"code":    200,
		"msg":     fmt.Sprintf("清理完成，共删除 %d 个孤儿文件（超过 %d 小时未认领）", deletedCount, hours),
		"deleted": deletedFiles,
	})
}
