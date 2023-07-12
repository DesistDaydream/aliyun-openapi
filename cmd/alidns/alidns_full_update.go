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

func alidnsFullUpdateCommand() *cobra.Command {
	alidnsFullUpdateCmd := &cobra.Command{
		Use:   "full-update",
		Short: "全量更新指定域名的解析记录。注意: 该操作会批量删除后再批量创建！",
		Run:   runAlidnsFullUpdate,
	}

	return alidnsFullUpdateCmd
}

func runAlidnsFullUpdate(cmd *cobra.Command, args []string) {
	// 检查文件是否存在
	checkFile(alidnsFlags.rrFile)
	// 从文件中获取需要批量添加的解析记录
	domainRecordInfos, err := handleFile(alidnsFlags.rrFile, alidnsFlags.domainName)
	if err != nil {
		panic(err)
	}
	fullUpdate(domainRecordInfos, r, q, d, alidnsFlags.domainName)
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
