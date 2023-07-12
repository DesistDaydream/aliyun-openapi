package alidns

import (
	"os"

	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/domain"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/queryresults"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/resolve"
	"github.com/DesistDaydream/aliyun-openapi/pkg/fileparse"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
)

type AliDNSFlags struct {
	operation  string
	batchType  string
	domainName string
	rrFile     string
}

var alidnsFlags AliDNSFlags

func CreateCommand() *cobra.Command {
	long := `域名解析记录可用的操作类型：
	full-update 【！！！高危操作！！！】全量更新域名的解析记录，先删除原有的解析记录，再添加新的解析记录
	list 列出所有记录规则
	batch 批量操作.包括如下几种: [RR_ADD,RR_DEL,DOMAIN_ADD,DOMAIN_DEL]

	注意：若想要增量更新域名的记录规则，使用 batch 的 RR_ADD 操作即可`
	alidnsCmd := &cobra.Command{
		Use:   "alidns",
		Short: "云解析",
		Long:  long,
		// Run:   runAlidns,
	}

	cobra.OnInitialize(initConfig)

	alidnsCmd.PersistentFlags().StringVarP(&alidnsFlags.operation, "operation", "o", "list", "操作类型")
	alidnsCmd.PersistentFlags().StringVarP(&alidnsFlags.batchType, "batch-type", "O", "", "批量操作类型")
	alidnsCmd.PersistentFlags().StringVarP(&alidnsFlags.domainName, "domain-name", "d", "", "域名")
	alidnsCmd.PersistentFlags().StringVarP(&alidnsFlags.rrFile, "rr-file", "f", "", "存有域名资源记录的文件")

	alidnsCmd.AddCommand(
		alidnsListCommand(),
		alidnsUpdateCommand(),
		alidnsFullUpdateCommand(),
		alidnsBatchCommand(),
	)

	return alidnsCmd
}

var (
	h *alidns.AlidnsHandler
	d *domain.AlidnsDomain
	r *resolve.AlidnsResolve
	q *queryresults.AlidnsQueryResults
)

func initConfig() {
	if alidnsFlags.domainName == "" {
		logrus.Fatal("请使用 -d 标志指定要操作的域名")
	}

	// 实例化各种 API 处理器
	h = alidns.NewAlidnsHandler(aliclient.Info.AK, aliclient.Info.SK, alidnsFlags.domainName, aliclient.Info.Region)
	d = domain.NewAlidnsDomain(h)
	r = resolve.NewAlidnsResolve(h)
	q = queryresults.NewQueryResults(h)
}

// 处理 Excel 文件，读取 Excel 文件中的数据，并转换成调用批量操作方法时所需的 OperateBatchDomainRequestDomainRecordInfo 结构体
func handleFile(file string, domainName string) ([]*alidns20150109.OperateBatchDomainRequestDomainRecordInfo, error) {
	var domainRecordInfos []*alidns20150109.OperateBatchDomainRequestDomainRecordInfo

	data, err := fileparse.NewExcelData(file, domainName)
	if err != nil {
		logrus.Errorf("fileparse.NewExcelData error: %v", err)
		return nil, err
	}

	for _, row := range data.Rows {
		domainRecordInfos = append(domainRecordInfos, &alidns20150109.OperateBatchDomainRequestDomainRecordInfo{
			Domain: tea.String(domainName),
			Rr:     tea.String(row.Host),
			Type:   tea.String(row.Type),
			Value:  tea.String(row.Value),
		})
	}

	return domainRecordInfos, nil
}

func checkFile(rrFile string) {
	if _, err := os.Stat(rrFile); os.IsNotExist(err) {
		logrus.Fatalf("【%v】文件不存在，请使用 -f 指定域名的记录规则文件", rrFile)
	}
}
