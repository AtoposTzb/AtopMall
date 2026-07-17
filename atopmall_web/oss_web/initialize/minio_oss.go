package initialize

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"atopmall_web/oss_web/global"
)

// 直接读取全局 ServerConfig 里的 MinIOInfo
func InitMinIO() error {
	conf := global.ServerConfig.MinIOInfo
	// 创建MinIO客户端
	cli, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("创建MinIO客户端失败: %w", err)
	}

	// 自动创建桶（不存在则新建）
	ctx := context.Background()
	hasBucket, err := cli.BucketExists(ctx, conf.BucketName)
	if err != nil {
		return fmt.Errorf("检测桶失败: %w", err)
	}
	if !hasBucket {
		err = cli.MakeBucket(ctx, conf.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("创建桶失败: %w", err)
		}
	}

	global.MinioCli = cli //将 MinIO客户端存储到全局变量中
	return nil
}
