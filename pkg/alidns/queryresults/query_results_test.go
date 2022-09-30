package queryresults

import (
	"testing"

	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/sirupsen/logrus"
)

func TestAlidnsQueryResults_QueryResults(t *testing.T) {
	// 准备测试数据
	domainName := "desistdaydream.ltd"
	auth := config.NewAuthInfo("../../../owner.yaml")
	ak := auth.AuthList["断灬念梦"].AccessKeyID
	sk := auth.AuthList["断灬念梦"].AccessKeySecret
	handler := alidns.NewAlidnsHandler(ak, sk, domainName, "alidns.cn-beijing.aliyuncs.com")
	q := NewQueryResults(handler)

	var taskID int64 = 4521793536
	got, err := q.QueryResults(taskID, "")
	if err != nil {
		logrus.Errorln(err)
	}

	logrus.Infoln(got)
}
