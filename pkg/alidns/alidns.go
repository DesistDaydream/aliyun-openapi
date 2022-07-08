package alidns

import (
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/spf13/pflag"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func createClient(accessKeyId string, accessKeySecret string, region string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: &accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: &accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(region)
	// _result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

// 命令行标志
type AlidnsFlags struct {
	DomainName string
	RRFile     string
}

// 设置命令行标志
func (flags *AlidnsFlags) AddFlags() {
	pflag.StringVarP(&flags.DomainName, "domain", "d", "", "域名")
	pflag.StringVarP(&flags.RRFile, "rr-file", "f", "", "存有域名资源记录的文件")
}

type AlidnsHandler struct {
	Client     *alidns20150109.Client
	DomainName string
	Runtime    *util.RuntimeOptions
}

// 实例化云解析处理器
func NewAlidnsHandler(auth *config.AuthConfig, userName, domainName, region string) *AlidnsHandler {
	client, err := createClient(auth.AuthList[userName].AccessKeyID, auth.AuthList[userName].AccessKeySecret, region)
	if err != nil {
		panic(err)
	}

	return &AlidnsHandler{
		DomainName: domainName,
		Runtime:    &util.RuntimeOptions{},
		Client:     client,
	}
}
