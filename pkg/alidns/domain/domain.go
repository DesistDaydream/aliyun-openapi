package domain

import (
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
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
func (d *AlidnsDomain) Batch(operateType string, domainRecordInfos []*alidns20150109.OperateBatchDomainRequestDomainRecordInfo) (int64, error) {
	operateBatchDomainRequest := &alidns20150109.OperateBatchDomainRequest{
		Type:             tea.String(operateType),
		DomainRecordInfo: domainRecordInfos,
	}

	result, err := d.AlidnsHandler.Client.OperateBatchDomainWithOptions(operateBatchDomainRequest, d.AlidnsHandler.Runtime)
	if err != nil {
		logrus.Fatalln(err)
		return 0, err
	}

	logrus.WithFields(logrus.Fields{
		"任务ID": *result.Body.TaskId,
		"请求ID": *result.Body.RequestId,
	}).Info("批量任务运行信息")

	return *result.Body.TaskId, nil
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
