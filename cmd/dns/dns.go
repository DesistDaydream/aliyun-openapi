// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"encoding/json"
	"os"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("alidns.cn-beijing.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

// 认证信息
type AuthInfo struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

func NewAuthInfo() (auth *AuthInfo) {
	// 读取认证信息
	fileByte, err := os.ReadFile("auth.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileByte, &auth)
	if err != nil {
		panic(err)
	}
	return auth
}

func _main(args []*string) (_err error) {
	auth := NewAuthInfo()
	logrus.Info(auth)

	// 初始化账号Client
	client, _err := CreateClient(tea.String(auth.AccessKeyId), tea.String(auth.AccessKeySecret))
	if _err != nil {
		return _err
	}

	describeDomainsRequest := &alidns20150109.DescribeDomainsRequest{}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = client.DescribeDomainsWithOptions(describeDomainsRequest, runtime)
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		util.AssertAsString(error.Message)
	}
	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
