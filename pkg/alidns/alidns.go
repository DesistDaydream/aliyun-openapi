package alidns

import (
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"

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

type AlidnsHandler struct {
	Client     *alidns20150109.Client
	DomainName string
	Runtime    *util.RuntimeOptions
}

// 实例化云解析处理器
func NewAlidnsHandler(ak, sk, domainName, region string) *AlidnsHandler {
	client, err := createClient(ak, sk, region)
	if err != nil {
		panic(err)
	}

	return &AlidnsHandler{
		DomainName: domainName,
		Runtime:    &util.RuntimeOptions{},
		Client:     client,
	}
}
