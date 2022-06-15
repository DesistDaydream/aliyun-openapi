package alidns

import (
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/spf13/pflag"
)

// 命令行标志
type AlidnsFlags struct {
	DomainName string
	AuthFile   string
	RRFile     string
}

// 设置命令行标志
func (flags *AlidnsFlags) AddFlags() {
	pflag.StringVarP(&flags.DomainName, "domain", "d", "", "域名")
	pflag.StringVarP(&flags.AuthFile, "auth-file", "F", "auth.yaml", "认证信息文件")
	pflag.StringVarP(&flags.RRFile, "rr-file", "f", "", "存有域名资源记录的文件")
}

type AlidnsHandler struct {
	DomainName string
	Runtime    *util.RuntimeOptions
	Client     *alidns20150109.Client
}

// 实例化云解析处理器
func NewAlidnsHandler(flags *AlidnsFlags, client *alidns20150109.Client) *AlidnsHandler {
	return &AlidnsHandler{
		DomainName: flags.DomainName,
		Runtime:    &util.RuntimeOptions{},
		Client:     client,
	}
}
