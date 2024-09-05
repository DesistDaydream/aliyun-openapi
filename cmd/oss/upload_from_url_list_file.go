package oss

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type UploadFromFileFlags struct {
	file string
}

var uploadFromFileFlags UploadFromFileFlags

func uploadFromFileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-from-url-list-file",
		Short: "从 URL 列表中将 URL 指向的文件上传到阿里云 OSS",
		Run:   uploadFromUrlListFile,
		PreRun: func(cmd *cobra.Command, args []string) {
			getBucket()
		},
	}

	cmd.Flags().StringVar(&uploadFromFileFlags.file, "file", "test/files/url_list.txt", "URL 列表文件路径")

	return cmd
}

func uploadFromUrlListFile(cmd *cobra.Command, args []string) {
	file, err := os.Open(uploadFromFileFlags.file)
	if err != nil {
		logrus.Fatalf("无法打开文件：%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()

		// 下载图片
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("下载图片失败:", err)
			continue
		}
		defer resp.Body.Close()

		// 从 URL 中提取文件名
		fileName := filepath.Base(url)

		// 上传到 OSS
		err = bucket.PutObject(fmt.Sprintf("%s/%s", aliossFlags.UploadPath, fileName), resp.Body)
		if err != nil {
			fmt.Println("上传图片失败:", err)
			continue
		}

		fmt.Printf("图片 %s 上传成功\n", fileName)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件失败:", err)
	}
}
