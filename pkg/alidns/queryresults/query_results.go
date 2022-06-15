package queryresults

import (
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/sirupsen/logrus"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
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
func (qr *AlidnsQueryResults) QueryResults(taskID int64) (int32, error) {
	describeBatchResultCountRequest := &alidns20150109.DescribeBatchResultCountRequest{
		TaskId: tea.Int64(taskID),
	}

	result, err := qr.AlidnsHandler.Client.DescribeBatchResultCountWithOptions(describeBatchResultCountRequest, qr.AlidnsHandler.Runtime)
	if err != nil {
		logrus.Errorln(err)
		return 0, err
	}

	logrus.WithFields(logrus.Fields{
		"任务ID": *result.Body.TaskId,
		"任务结果": *result.Body.Status,
		"成功数":  *result.Body.SuccessCount,
		"失败数":  *result.Body.FailedCount,
	}).Info("查询批量任务结果")

	return *result.Body.Status, err
}
