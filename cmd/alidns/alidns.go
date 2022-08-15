package alidns

import (
	"os"
	"time"

	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/domain"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/queryresults"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/resolve"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	long := `域名解析记录可用的操作类型：
	update 更新域名的解析记录，先全部删除再逐一添加
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

	AlidnsCmd.Flags().StringP("operation", "o", "list", "操作类型")
	AlidnsCmd.Flags().StringP("batch-operation", "O", "", "批量操作类型")
	AlidnsCmd.Flags().StringP("domain", "d", "", "域名")
	AlidnsCmd.Flags().StringP("rr-file", "f", "", "存有域名资源记录的文件")

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
	operation, _ := cmd.Flags().GetString("operation")
	batchOperation, _ := cmd.Flags().GetString("batch-operation")
	domainName, err := cmd.Flags().GetString("domain")
	if err != nil || domainName == "" {
		logrus.Fatal("请使用 -d 标志指定要操作的域名")
	}
	rrFile, _ := cmd.Flags().GetString("rr-file")

	// 实例化各种 API 处理器
	h := alidns.NewAlidnsHandler(aliclient.Info.AK, aliclient.Info.SK, domainName, aliclient.Info.Region)
	d := domain.NewAlidnsDomain(h)
	r := resolve.NewAlidnsResolve(h)
	q := queryresults.NewQueryResults(h)

	switch operation {
	case "list":
		domainRecords, err := r.DomainRecordsList()
		if err != nil {
			panic(err)
		}
		logrus.Infoln(domainRecords)
	case "update":
		// 检查文件是否存在
		if _, err := os.Stat(rrFile); os.IsNotExist(err) {
			logrus.Fatalf("【%v】文件不存在，请使用 -f 指定域名的记录规则文件", rrFile)
		}

		// 添加解析记录前先删除全部解析记录
		taskID, err := d.Batch("RR_DEL", rrFile)
		if err != nil {
			logrus.Fatal(err)
		}

		// 根据 taskID 持续查询删除任务完成状态，任务完成后再执行后续代码
		for {
			task, err := q.QueryResults(taskID)
			if err != nil {
				logrus.Fatal(err)
			}
			if task == 1 {
				logrus.Infof("域名解析记录已批量删除，开始执行添加解析记录的操作")
				break
			}
			time.Sleep(time.Second * 1)
		}
		// 这里为了测试目的，使用的是逐一添加的方式，实际应用中可以使用批量添加的方式
		r.OnebyoneAddDomainRecord(rrFile)
	case "batch":
		// 检查文件是否存在
		if _, err := os.Stat(rrFile); os.IsNotExist(err) {
			logrus.Fatalf("【%v】文件不存在，请使用 -f 指定域名的记录规则文件", rrFile)
		}
		// 判断批量操作类型是否存在
		if batchOperation == "" {
			logrus.Fatal("请使用 -O 标志指定批量操作类型")
		}
		// 判断批量操作类型是否合法
		if !d.IsBatchOperationExist(batchOperation) {
			logrus.Fatal("批量操作类型不存在，可用的值有: RR_ADD,RR_DEL,DOMAIN_ADD,DOMAIN_DEL")
		}
		taskID, err := d.Batch(batchOperation, rrFile)
		if err != nil {
			logrus.Errorf("执行【%v】操作失败，错误信息: %v", batchOperation, err)
		}
		// 根据 taskID 持续查询删除任务完成状态，任务完成后再执行后续代码
		for {
			task, err := q.QueryResults(taskID)
			if err != nil {
				logrus.Fatal(err)
			}
			if task == 1 {
				logrus.Infof("执行【%v】操作成功", batchOperation)
				break
			}
			time.Sleep(time.Second * 1)
		}
	default:
		logrus.Fatalln("操作类型不存在，请使用 -o 指定操作类型")
	}
}
