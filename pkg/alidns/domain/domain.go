package domain

import (
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/fileparse"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
)

type AlidnsDomain struct {
	AlidnsHandler *alidns.AlidnsHandler
}

// 实例化 alidnsdomain
func NewAlidnsDomain(alidnsHandler *alidns.AlidnsHandler) *AlidnsDomain {
	return &AlidnsDomain{
		AlidnsHandler: alidnsHandler,
	}
}

// 批量操作
// operateType 可用值如下：
// DOMAIN_ADD：批量添加域名
// DOMAIN_DEL：批量删除域名
// RR_ADD：批量添加解析
// RR_DEL：批量删除解析（删除满足N.RR、N.VALUE、N.RR&amp;N.VALUE条件的解析记录。如果无N.RR&&N.VALUE则清空参数DomainRecordInfo.N.Domain下的解析记录）
func (d *AlidnsDomain) Batch(operateType string, file string) {
	var domainRecordInfos []*alidns20150109.OperateBatchDomainRequestDomainRecordInfo

	// 处理 Excel 文件，读取 Excel 文件中的数据，并转换成 OperateBatchDomainRequestDomainRecordInfo 结构体
	data := fileparse.NewExcelData(file, d.AlidnsHandler.DomainName)

	for _, row := range data.Rows {
		var domainRecordInfo alidns20150109.OperateBatchDomainRequestDomainRecordInfo
		domainRecordInfo.Type = tea.String(row.Type)
		domainRecordInfo.Value = tea.String(row.Value)
		domainRecordInfo.Rr = tea.String(row.Host)
		domainRecordInfo.Domain = tea.String(d.AlidnsHandler.DomainName)

		domainRecordInfos = append(domainRecordInfos, &domainRecordInfo)
	}

	logrus.Debugln(domainRecordInfos)

	operateBatchDomainRequest := &alidns20150109.OperateBatchDomainRequest{
		Type:             tea.String(operateType),
		DomainRecordInfo: domainRecordInfos,
	}

	result, err := d.AlidnsHandler.Client.OperateBatchDomainWithOptions(operateBatchDomainRequest, d.AlidnsHandler.Runtime)
	if err != nil {
		panic(err)
	}
	logrus.Info(result)
}

// 判断批量操作类型是否合法
func (d *AlidnsDomain) IsBatchOperationExist(operateType string) bool {
	switch operateType {
	case "DOMAIN_ADD":
		return true
	case "DOMAIN_DEL":
		return true
	case "RR_ADD":
		return true
	case "RR_DEL":
		return true
	default:
		return false
	}
}
