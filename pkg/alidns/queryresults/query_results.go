package queryresults

import (
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/sirupsen/logrus"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

type AlidnsQueryResults struct {
	AlidnsHandler *alidns.AlidnsHandler
}

// 实例化 alidnsdomain
func NewQueryResults(alidnsHandler *alidns.AlidnsHandler) *AlidnsQueryResults {
	return &AlidnsQueryResults{
		AlidnsHandler: alidnsHandler,
	}
}

// 根据 taskID 查询批量处理任务的结果
func (qr *AlidnsQueryResults) QueryResults(taskID int64, batchType string) (int32, error) {
	describeBatchResultCountRequest := &alidns20150109.DescribeBatchResultCountRequest{
		TaskId: tea.Int64(taskID),
	}

	result, err := qr.AlidnsHandler.Client.DescribeBatchResultCountWithOptions(describeBatchResultCountRequest, qr.AlidnsHandler.Runtime)
	if err != nil {
		logrus.Errorln(err)
		return 0, err
	}

	if *result.Body.Status == -1 {
		logrus.Errorln("任务 ID 不存在")
		return 0, err
	}

	logrus.WithFields(logrus.Fields{
		"任务ID": *result.Body.TaskId,
		"任务状态": *result.Body.Status,
		"总数":   *result.Body.TotalCount,
		"成功数":  *result.Body.SuccessCount,
		"失败数":  *result.Body.FailedCount,
	}).Infof("查询 %v 批量任务结果", batchType)

	return *result.Body.Status, err
}
