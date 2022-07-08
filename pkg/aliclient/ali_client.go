package aliclient

import "github.com/spf13/pflag"

type ClientInfo struct {
	AK     string
	SK     string
	Region string
}

// 用于在根命令的 PersistentPreRun 中初始化账号Client
// 以便在子命令中直接引用该变量的值
var Info *ClientInfo

// 命令行标志
type AlidnsFlags struct {
	AuthFile string
	UserName string
	Region   string
}

// 设置命令行标志
func (flags *AlidnsFlags) AddFlags() {
	pflag.StringVarP(&flags.AuthFile, "auth-file", "F", "auth.yaml", "认证信息文件")
	pflag.StringVarP(&flags.UserName, "user-name", "u", "", "用户名")
	pflag.StringVarP(&flags.Region, "region", "r", "alidns.cn-beijing.aliyuncs.com", "区域")
}
