package oss

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listBucketCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出 Bucket",
		Run:   listBucket,
	}

	return cmd
}

func listBucket(cmd *cobra.Command, args []string) {
	// 获取所有 bucket
	buckets, err := client.ListBuckets()
	if err != nil {
		fmt.Printf("列出 Bucket 失败: %v\n", err)
		return
	}

	// 遍历并输出所有 bucket 的名称
	fmt.Println("Bucket 列表:")
	for _, bucket := range buckets.Buckets {
		fmt.Printf("- %s\n", bucket.Name)
	}
}
