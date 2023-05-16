package alidns

import (
	"os"
	"time"

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
	AlidnsCmd := &cobra.Command{
		Use:              "alidns",
		Short:            "云解析",
		Long:             long,
		PersistentPreRun: alidnsPersistentPreRun,
		Run:              runAlidns,
	}

	AlidnsCmd.Flags().StringVarP(&alidnsFlags.operation, "operation", "o", "list", "操作类型")
	AlidnsCmd.Flags().StringVarP(&alidnsFlags.batchType, "batch-type", "O", "", "批量操作类型")
	AlidnsCmd.Flags().StringVarP(&alidnsFlags.domainName, "domain-name", "d", "", "域名")
	AlidnsCmd.Flags().StringVarP(&alidnsFlags.rrFile, "rr-file", "f", "", "存有域名资源记录的文件")

	return AlidnsCmd
}

// 执行 alidns 子命令之前需要执行的操作
func alidnsPersistentPreRun(cmd *cobra.Command, args []string) {
	// 执行根命令的初始化操作
	parent := cmd.Parent()
	if parent.PersistentPreRun != nil {
		parent.PersistentPreRun(parent, args)
	}
}

// 执行 alidns 子命令
func runAlidns(cmd *cobra.Command, args []string) {
	if alidnsFlags.domainName == "" {
		logrus.Fatal("请使用 -d 标志指定要操作的域名")
	}

	// 实例化各种 API 处理器
	h := alidns.NewAlidnsHandler(aliclient.Info.AK, aliclient.Info.SK, alidnsFlags.domainName, aliclient.Info.Region)
	d := domain.NewAlidnsDomain(h)
	r := resolve.NewAlidnsResolve(h)
	q := queryresults.NewQueryResults(h)

	switch alidnsFlags.operation {
	case "list":
		domainRecords, err := r.DomainRecordsList()
		if err != nil {
			panic(err)
		}

		for index, domainRecord := range domainRecords.Record {
			logrus.WithFields(logrus.Fields{
				"类型": *domainRecord.Type,
				"记录": *domainRecord.RR,
				"值":  *domainRecord.Value,
			}).Infof("%v 域名的第 %v 条资源记录", alidnsFlags.domainName, index+1)
		}
		logrus.Infof("共有 %d 条记录", len(domainRecords.Record))
	case "full-update":
		// 检查文件是否存在
		checkFile(alidnsFlags.rrFile)
		// 从文件中获取需要批量添加的解析记录
		domainRecordInfos, err := handleFile(alidnsFlags.rrFile, alidnsFlags.domainName)
		if err != nil {
			panic(err)
		}
		fullUpdate(domainRecordInfos, r, q, d, alidnsFlags.domainName)
	case "update":
		// TODO: 只更新
	case "batch":
		// 检查文件是否存在
		checkFile(alidnsFlags.rrFile)
		domainRecordInfos, err := handleFile(alidnsFlags.rrFile, alidnsFlags.domainName)
		if err != nil {
			panic(err)
		}
		batch(domainRecordInfos, r, q, d, alidnsFlags.domainName, alidnsFlags.batchType)
	default:
		logrus.Fatalln("操作类型不存在，请使用 -o 指定操作类型")
	}
}

func fullUpdate(domainRecordInfos []*alidns20150109.OperateBatchDomainRequestDomainRecordInfo, r *resolve.AlidnsResolve, q *queryresults.AlidnsQueryResults, d *domain.AlidnsDomain, domainName string) {
	batchTypeDel := "RR_DEL"
	batchTypeAdd := "RR_ADD"
	// 列出所有域名记录
	domainRecords, err := r.DomainRecordsList()
	if err != nil {
		panic(err)
	}

	// 如果列出的域名记录不为空，则先批量删除所有列出的域名记录
	if len(domainRecords.Record) > 0 {
		var needDeleteRecords []*alidns20150109.OperateBatchDomainRequestDomainRecordInfo
		for _, record := range domainRecords.Record {
			needDeleteRecords = append(needDeleteRecords, &alidns20150109.OperateBatchDomainRequestDomainRecordInfo{
				Type:   record.Type,
				Value:  record.Value,
				Domain: record.DomainName,
			})
		}

		logrus.Debugf("需要删除 %v 条资源记录", len(needDeleteRecords))

		delTaskID, err := d.Batch(batchTypeDel, needDeleteRecords)
		if err != nil {
			logrus.Fatal(err)
		}

		// 根据 taskID 持续查询删除任务完成状态，任务完成后再执行后续代码
		for {
			task, err := q.QueryResults(delTaskID, batchTypeDel)
			if err != nil {
				logrus.Fatal(err)
			}
			if task == 1 {
				logrus.Infof("域名解析记录已批量删除，开始执行添加解析记录的操作")
				break
			}
			time.Sleep(time.Second * 1)
		}
	}

	// 批量添加解析记录
	addTaskID, err := d.Batch(batchTypeAdd, domainRecordInfos)
	if err != nil {
		logrus.Fatal(err)
	}

	// 根据 taskID 持续查询添加任务完成状态
	for {
		task, err := q.QueryResults(addTaskID, batchTypeDel)
		if err != nil {
			logrus.Fatal(err)
		}
		if task == 1 {
			logrus.Infof("域名解析记录已批量添加")
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func batch(domainRecordInfos []*alidns20150109.OperateBatchDomainRequestDomainRecordInfo, r *resolve.AlidnsResolve, q *queryresults.AlidnsQueryResults, d *domain.AlidnsDomain, domainName string, batchType string) {
	// 判断批量操作类型是否存在
	if batchType == "" {
		logrus.Fatal("请使用 -O 标志指定批量操作类型")
	}
	// 判断批量操作类型是否合法
	if !d.IsBatchOperationExist(batchType) {
		logrus.Fatal("批量操作类型不存在，可用的值有: RR_ADD,RR_DEL,DOMAIN_ADD,DOMAIN_DEL")
	}

	taskID, err := d.Batch(batchType, domainRecordInfos)
	if err != nil {
		logrus.Errorf("执行【%v】操作失败，错误信息: %v", batchType, err)
	}
	// 根据 taskID 持续查询删除任务完成状态，任务完成后再执行后续代码
	for {
		task, err := q.QueryResults(taskID, batchType)
		if err != nil {
			logrus.Fatal(err)
		}
		if task == 1 {
			logrus.Infof("执行【%v】操作成功", batchType)
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func handleFile(file string, domainName string) ([]*alidns20150109.OperateBatchDomainRequestDomainRecordInfo, error) {
	var domainRecordInfos []*alidns20150109.OperateBatchDomainRequestDomainRecordInfo

	// 处理 Excel 文件，读取 Excel 文件中的数据，并转换成 OperateBatchDomainRequestDomainRecordInfo 结构体
	data, err := fileparse.NewExcelData(file, domainName)
	if err != nil {
		logrus.Errorf("fileparse.NewExcelData error: %v", err)
		return nil, err
	}

	for _, row := range data.Rows {
		var domainRecordInfo alidns20150109.OperateBatchDomainRequestDomainRecordInfo
		domainRecordInfo.Type = tea.String(row.Type)
		domainRecordInfo.Value = tea.String(row.Value)
		domainRecordInfo.Rr = tea.String(row.Host)
		domainRecordInfo.Domain = tea.String(domainName)

		domainRecordInfos = append(domainRecordInfos, &domainRecordInfo)
	}

	return domainRecordInfos, nil
}

func checkFile(rrFile string) {
	if _, err := os.Stat(rrFile); os.IsNotExist(err) {
		logrus.Fatalf("【%v】文件不存在，请使用 -f 指定域名的记录规则文件", rrFile)
	}
}
