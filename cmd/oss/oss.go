package oss

import (
	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type AliOSSFlags struct {
	Endpoint   string
	BucketName string
	UploadPath string
}

var aliossFlags AliOSSFlags

func CreateCommand() *cobra.Command {
	long := `TODO`
	cmd := &cobra.Command{
		Use:              "oss",
		Short:            "对象存储",
		Long:             long,
		PersistentPreRun: preRun,
	}

	cmd.PersistentFlags().StringVar(&aliossFlags.Endpoint, "endpoint", "oss-cn-beijing.aliyuncs.com", "阿里云OSS的Endpoint")
	cmd.PersistentFlags().StringVar(&aliossFlags.BucketName, "bucket", "", "阿里云OSS的BucketName")
	cmd.PersistentFlags().StringVar(&aliossFlags.UploadPath, "upload-path", "", "上传文件路径")

	cmd.AddCommand(
		uploadFromFileCommand(),
		listBucketCommand(),
	)

	return cmd
}

var (
	client *oss.Client
	bucket *oss.Bucket
)

func preRun(cmd *cobra.Command, args []string) {
	if aliclient.Info == nil {
		logrus.Fatalf("阿里云客户端信息为空，请先初始化阿里云客户端")
	}

	var err error

	// 实例化 OSS Client
	client, err = oss.New(aliossFlags.Endpoint, aliclient.Info.AK, aliclient.Info.SK)
	if err != nil {
		logrus.Fatalf("实例化 OSS 客户端失败: %v", err)
	}
}

func getBucket() {
	isExist, err := client.IsBucketExist(aliossFlags.BucketName)
	if err != nil || !isExist {
		logrus.Fatalf("Bucket [%v] 不存在: %v", aliossFlags.BucketName, err)
	}

	bucket, err = client.Bucket(aliossFlags.BucketName)
	if err != nil {
		logrus.Fatalf("获取 Bucket [%v] 信息失败: %v", aliossFlags.BucketName, err)
	}
}
