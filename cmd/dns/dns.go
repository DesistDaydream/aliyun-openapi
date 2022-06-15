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
	"github.com/spf13/pflag"

	"github.com/DesistDaydream/aliyun-openapi/pkg/file"
)

type DomainHandler struct {
	DomainName string
	Runtime    *util.RuntimeOptions
}

// 实例化域名处理器
func NewDomainHandler(flags *Flags) *DomainHandler {
	return &DomainHandler{
		DomainName: flags.DomainName,
		Runtime:    &util.RuntimeOptions{},
	}
}

// var runtime = &util.RuntimeOptions{}

// 获取解析记录列表
func (d *DomainHandler) DomainRecordsList(client *alidns20150109.Client) {
	// 发起 DescribeDomainRecords 请求时需要携带的参数
	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(d.DomainName),
	}

	// 使用参数调用 DescribeDomainRecords 接口
	dd, err := client.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, d.Runtime)
	if err != nil {
		panic(err)
	}
	logrus.Infoln(dd)
}

// 逐一添加解析记录
func (d *DomainHandler) BatchAddDomainRecord(client *alidns20150109.Client, excelFile string) {
	// 处理 Excel 文件，读取 Excel 文件中的数据，并转换成 OperateBatchDomainRequestDomainRecordInfo 结构体
	excelInfos := file.HandleExcel(excelFile, d.DomainName)

	file.HandleExcel(excelFile, d.DomainName)

	for _, info := range excelInfos {
		logrus.Infoln(info.Type, info.Value, info.Host)

		// 发起 AddDomainRecord 请求时需要携带的参数
		addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
			DomainName: tea.String(d.DomainName),
			Type:       tea.String(info.Type),
			Value:      tea.String(info.Value),
			RR:         tea.String(info.Host),
		}
		dd, err := client.AddDomainRecordWithOptions(addDomainRecordRequest, d.Runtime)
		if err != nil {
			logrus.Error(err)
		}

		logrus.Infoln(dd)
	}

}

// 批量删除全部解析记录
func (d *DomainHandler) BatchDeleteAll(client *alidns20150109.Client) {
	domainRecordInfo0 := &alidns20150109.OperateBatchDomainRequestDomainRecordInfo{
		Domain: tea.String(d.DomainName),
	}
	operateBatchDomainRequest := &alidns20150109.OperateBatchDomainRequest{
		DomainRecordInfo: []*alidns20150109.OperateBatchDomainRequestDomainRecordInfo{domainRecordInfo0},
		Type:             tea.String("RR_DEL"),
	}
	runtime := &util.RuntimeOptions{}

	result, err := client.OperateBatchDomainWithOptions(operateBatchDomainRequest, runtime)
	if err != nil {
		panic(err)
	}
	logrus.Info(result)
}

// 批量操作
// operateType 可用值如下：
// DOMAIN_ADD：批量添加域名
// DOMAIN_DEL：批量删除域名
// RR_ADD：批量添加解析
// RR_DEL：批量删除解析（删除满足N.RR、N.VALUE、N.RR&amp;N.VALUE条件的解析记录。如果无N.RR&&N.VALUE则清空参数DomainRecordInfo.N.Domain下的解析记录）
func (d *DomainHandler) Batch(client *alidns20150109.Client, operateType string, excelFile string) {

	var domainRecordInfos []*alidns20150109.OperateBatchDomainRequestDomainRecordInfo
	var domainRecordInfo alidns20150109.OperateBatchDomainRequestDomainRecordInfo

	// 处理 Excel 文件，读取 Excel 文件中的数据，并转换成 OperateBatchDomainRequestDomainRecordInfo 结构体
	excelInfos := file.HandleExcel(excelFile, d.DomainName)

	for _, info := range excelInfos {
		logrus.Infoln(info.Type, info.Value, info.Host)

		domainRecordInfo.Type = tea.String(info.Type)
		domainRecordInfo.Value = tea.String(info.Value)
		domainRecordInfo.Rr = tea.String(info.Host)
		domainRecordInfo.Domain = tea.String(d.DomainName)

		domainRecordInfos = append(domainRecordInfos, &domainRecordInfo)
	}

	logrus.Info(domainRecordInfos)

	operateBatchDomainRequest := &alidns20150109.OperateBatchDomainRequest{
		Type:             tea.String(operateType),
		DomainRecordInfo: domainRecordInfos,
	}

	logrus.Info(operateBatchDomainRequest)

	result, err := client.OperateBatchDomainWithOptions(operateBatchDomainRequest, d.Runtime)
	if err != nil {
		panic(err)
	}
	logrus.Info(result)
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId string, accessKeySecret string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: &accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: &accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("alidns.cn-beijing.aliyuncs.com")
	// _result = &alidns20150109.Client{}
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

// 命令行标志
type Flags struct {
	DomainName string
	File       string
}

// 设置命令行标志
func (flags *Flags) AddFlags() {
	pflag.StringVarP(&flags.DomainName, "domain", "d", "", "域名")
	pflag.StringVarP(&flags.File, "file", "f", "desistdaydream.ltd.xlsx", "Excel文件")
}

func main() {
	operation := pflag.StringP("operation", "o", "", "操作类型")
	// 添加命令行标志
	f := &Flags{}
	f.AddFlags()
	pflag.Parse()

	h := NewDomainHandler(f)

	auth := NewAuthInfo()
	logrus.Info(auth)

	// 初始化账号Client
	client, err := CreateClient(auth.AccessKeyId, auth.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	logrus.Debugln(client)

	// Batch(client, "RR_DEL", "desistdaydream.ltd.xlsx")

	switch *operation {
	case "list":
		h.DomainRecordsList(client)
	case "add":
		// h.BatchDeleteAll(client)
		// TODO: 根据任务 ID 判断删除任务是否完成;删除任务完成后再执行添加任务
		h.BatchAddDomainRecord(client, f.File)
	case "del-all":
		h.BatchDeleteAll(client)
	default:
		panic("操作类型不存在")
	}

}
