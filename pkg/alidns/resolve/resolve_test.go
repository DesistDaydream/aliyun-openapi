package resolve

import (
	"fmt"
	"testing"

	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/DesistDaydream/aliyun-openapi/pkg/fileparse"
	"github.com/sirupsen/logrus"
)

var (
	domainName string = "datalake.cn"
	userName   string = "ehualu_oc"
	region     string = "alidns.cn-beijing.aliyuncs.com"
)

func createClient() *AlidnsResolve {
	auth := config.NewAuthInfo("../../../owner.yaml")
	aliclient.Info = &aliclient.ClientInfo{
		AK:     auth.AuthList[userName].AccessKeyID,
		SK:     auth.AuthList[userName].AccessKeySecret,
		Region: region,
	}

	h := alidns.NewAlidnsHandler(aliclient.Info.AK, aliclient.Info.SK, domainName, aliclient.Info.Region)
	d := NewAlidnsResolve(h)

	return d
}

func TestAlidnsResolve_DomainRecordsList(t *testing.T) {
	d := createClient()
	domainRecords, err := d.DomainRecordsList()
	if err != nil {
		logrus.Fatalln(err)
	}

	for _, r := range domainRecords.Record {
		logrus.WithFields(logrus.Fields{
			"ID":   *r.RecordId,
			"域名":   *r.DomainName,
			"资源记录": *r.RR,
		}).Infoln("域名记录")
	}

	logrus.Infof("共有 %d 条记录", len(domainRecords.Record))

}

func TestAlidnsResolve_OnebyoneSetDomainRecordStatus(t *testing.T) {
	file := fmt.Sprintf("/mnt/d/Documents/WPS Cloud Files/1054253139/团队文档/设备文档与服务信息/域名解析/%s.xlsx", domainName)
	// 处理 Excel 文件，读取 Excel 文件中的数据，并转换成 OperateBatchDomainRequestDomainRecordInfo 结构体
	data, err := fileparse.NewExcelData(file, domainName)
	if err != nil {
		logrus.Errorln(err)
	}
	d := createClient()
	d.OnebyoneSetDomainRecordStatusToDisable(data)
}
