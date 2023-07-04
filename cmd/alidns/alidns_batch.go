package alidns

import (
	"time"

	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/domain"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/queryresults"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/resolve"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func alidnsBatchCommand() *cobra.Command {
	alidnsBatchCmd := &cobra.Command{
		Use:   "bathc",
		Short: "",
		Run:   runAlidnsBatch,
	}

	return alidnsBatchCmd
}

func runAlidnsBatch(cmd *cobra.Command, args []string) {
	// 检查文件是否存在
	checkFile(alidnsFlags.rrFile)
	domainRecordInfos, err := handleFile(alidnsFlags.rrFile, alidnsFlags.domainName)
	if err != nil {
		panic(err)
	}
	batch(domainRecordInfos, r, q, d, alidnsFlags.domainName, alidnsFlags.batchType)
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
