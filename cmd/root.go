package cmd

import (
	"os"

	"github.com/DesistDaydream/aliyun-openapi/cmd/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/DesistDaydream/aliyun-openapi/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Execute() {
	app := newApp()
	err := app.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func newApp() *cobra.Command {
	long := `对 alicloud-openapi 工具的长描述，包含用例等，比如:
	https://next.api.aliyun.com/home
	API 在线调试：https://next.api.aliyun.com/api`

	var RootCmd = &cobra.Command{
		Use:              "alicloud-openapi",
		Short:            "通过阿里云 OpenAPI 管理资源的工具",
		Long:             long,
		PersistentPreRun: rootPersistentPreRun,
	}

	RootCmd.PersistentFlags().StringP("log-level", "", "info", "日志级别:[debug, info, warn, error, fatal]")
	RootCmd.PersistentFlags().StringP("log-output", "", "", "日志输出位置，不填默认标准输出 stdout")
	RootCmd.PersistentFlags().StringP("log-format", "", "text", "日志输出格式: [text, json]")
	RootCmd.PersistentFlags().BoolP("log-caller", "", false, "是否输出函数名、文件名、行号")

	RootCmd.PersistentFlags().StringP("auth-file", "F", "owner.yaml", "认证信息文件")
	RootCmd.PersistentFlags().StringP("username", "u", "", "用户名")
	RootCmd.PersistentFlags().StringP("region", "r", "alidns.cn-beijing.aliyuncs.com", "区域")

	// 添加子命令
	RootCmd.AddCommand(
		alidns.CreateCommand(),
	// elb.CreateCommand(),
	)

	return RootCmd
}

// 执行每个 root 下的子命令时，都需要执行的函数
func rootPersistentPreRun(cmd *cobra.Command, args []string) {
	// 初始化日志
	logLevel, _ := cmd.Flags().GetString("log-level")
	logOutput, _ := cmd.Flags().GetString("log-output")
	logFormat, _ := cmd.Flags().GetString("log-format")
	logCaller, _ := cmd.Flags().GetBool("log-caller")
	if err := logging.LogInit(logLevel, logOutput, logFormat, logCaller); err != nil {
		logrus.Fatal("初始化日志失败", err)
	}

	// 认证信息文件处理的相关逻辑
	authFile, _ := cmd.Flags().GetString("auth-file")
	userName, err := cmd.Flags().GetString("username")
	if err != nil {
		logrus.Fatalln("请指定用户名")
	}
	region, _ := cmd.Flags().GetString("region")

	// 检查 clientFlags.AuthFile 文件是否存在
	if _, err := os.Stat(authFile); os.IsNotExist(err) {
		logrus.Fatalf("打开【%v】文件失败: %v", authFile, err)
	}
	// 获取认证信息
	auth := config.NewAuthInfo(authFile)

	// 判断传入的用户是否存在在认证信息中
	if !auth.IsUserExist(userName) {
		logrus.Fatalf("认证信息中不存在 %v 用户, 请检查认证信息文件或命令行参数的值", userName)
	}

	aliclient.Info = &aliclient.ClientInfo{
		AK:     auth.AuthList[userName].AccessKeyID,
		SK:     auth.AuthList[userName].AccessKeySecret,
		Region: region,
	}
}
