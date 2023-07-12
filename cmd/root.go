package cmd

import (
	"os"

	"github.com/DesistDaydream/aliyun-openapi/cmd/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	logging "github.com/DesistDaydream/logging/pkg/logrus_init"
)

type Flags struct {
	AuthFile string
	Username string
	Region   string
}

func AddFlags(f *Flags) {

}

var (
	flags    Flags
	logFlags logging.LogrusFlags
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
		Use:   "alicloud-openapi",
		Short: "通过阿里云 OpenAPI 管理资源的工具",
		Long:  long,
		// PersistentPreRun: rootPersistentPreRun,
	}

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&flags.AuthFile, "auth-file", "F", "pkg/config/my_auth.yaml", "认证信息文件")
	RootCmd.PersistentFlags().StringVarP(&flags.Username, "username", "u", "", "用户名")
	RootCmd.PersistentFlags().StringVarP(&flags.Region, "region", "r", "alidns.cn-beijing.aliyuncs.com", "区域")

	logging.AddFlags(&logFlags)

	// 添加子命令
	RootCmd.AddCommand(
		alidns.CreateCommand(),
		// elb.CreateCommand(),
	)

	return RootCmd
}

// 执行每个 root 下的子命令时，都需要执行的函数
func initConfig() {
	// 初始化日志
	if err := logging.LogrusInit(&logFlags); err != nil {
		logrus.Fatal("初始化日志失败", err)
	}

	// 检查 clientFlags.AuthFile 文件是否存在
	if _, err := os.Stat(flags.AuthFile); os.IsNotExist(err) {
		logrus.Fatalf("打开【%v】文件失败: %v", flags.AuthFile, err)
	}
	// 获取认证信息
	auth := config.NewAuthInfo(flags.AuthFile)

	// 判断传入的用户是否存在在认证信息中
	if !auth.IsUserExist(flags.Username) {
		logrus.Fatalf("认证信息中不存在 %v 用户, 请检查认证信息文件或命令行参数的值", flags.Username)
	}

	aliclient.Info = &aliclient.ClientInfo{
		AK:     auth.AuthList[flags.Username].AccessKeyID,
		SK:     auth.AuthList[flags.Username].AccessKeySecret,
		Region: flags.Region,
	}
}
