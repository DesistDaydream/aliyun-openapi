package resolve

import (
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/fileparse"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
)

type AlidnsResolve struct {
	AlidnsHandler *alidns.AlidnsHandler
}

// 实例化 alidnsdomain
func NewAlidnsResolve(alidnsHandler *alidns.AlidnsHandler) *AlidnsResolve {
	return &AlidnsResolve{
		AlidnsHandler: alidnsHandler,
	}
}

// 获取解析记录列表 DescribeDomainRecords
func (d *AlidnsResolve) DomainRecordsList() (*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecords, error) {
	// 一次查询可以返回的最大条目数量
	var pageSize int64 = 20
	var pageNumber int64 = 1
	// var domainRecords *alidns20150109.DescribeDomainRecordsResponseBodyDomainRecords
	// 发起 DescribeDomainRecords 请求时需要携带的参数
	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(d.AlidnsHandler.DomainName),
		PageSize:   tea.Int64(pageSize),   // 一次查询可以返回的最大条目数量，取值范围为1~500，默认为20
		PageNumber: tea.Int64(pageNumber), // 分页查询的页码，默认为1
	}

	// 使用参数调用 DescribeDomainRecords 接口
	dd, err := d.AlidnsHandler.Client.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, d.AlidnsHandler.Runtime)
	if err != nil {
		return nil, err
	}

	// 如果查询到的记录条数大于 pageSize 的值，那么需要分页查询。并将查询到的记录合并
	if *dd.Body.TotalCount/pageSize >= 1 && *dd.Body.TotalCount%pageSize != 0 {
		page := int(*dd.Body.TotalCount/pageSize + 1)

		for i := 2; i <= page; i++ {
			describeDomainRecordsRequest.PageNumber = tea.Int64(int64(i))
			dr, err := d.AlidnsHandler.Client.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, d.AlidnsHandler.Runtime)
			if err != nil {
				return nil, err
			}

			dd.Body.DomainRecords.Record = append(dd.Body.DomainRecords.Record, dr.Body.DomainRecords.Record...)
		}
	}

	return dd.Body.DomainRecords, nil
}

// 逐一添加解析记录 AddDomainRecord
func (d *AlidnsResolve) OnebyoneAddDomainRecord(file string) {
	// 处理 Excel 文件，读取 Excel 文件中的数据，并转换成 OperateBatchDomainRequestDomainRecordInfo 结构体
	data, err := fileparse.NewExcelData(file, d.AlidnsHandler.DomainName)
	if err != nil {
		logrus.Errorln(err)
	}

	for _, row := range data.Rows {
		// 发起 AddDomainRecord 请求时需要携带的参数
		addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
			DomainName: tea.String(d.AlidnsHandler.DomainName),
			Type:       tea.String(row.Type),
			Value:      tea.String(row.Value),
			RR:         tea.String(row.Host),
		}
		dd, err := d.AlidnsHandler.Client.AddDomainRecordWithOptions(addDomainRecordRequest, d.AlidnsHandler.Runtime)
		if err != nil {
			logrus.Errorf("添加记录失败\n%v", err)
		} else {
			logrus.WithFields(logrus.Fields{
				"记录类型": row.Type,
				"记录值":  row.Value,
				"主机记录": row.Host,
			}).Infof("记录添加成功")

			logrus.Debugf("检查添加成功的响应结果: %v", dd)
		}
	}
}
