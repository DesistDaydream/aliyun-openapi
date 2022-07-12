package main

import (
	"fmt"

	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/queryresults"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/sirupsen/logrus"
)

func main() {
	taskID := int64(12345)
	domainName := "desistdaydream.ltd"
	auth := config.NewAuthInfo("owner.yaml")
	ak := auth.AuthList["断灬念梦"].AccessKeyID
	sk := auth.AuthList["断灬念梦"].AccessKeySecret
	handler := alidns.NewAlidnsHandler(ak, sk, domainName, "alidns.cn-beijing.aliyuncs.com")
	q := queryresults.NewQueryResults(handler)

	task, err := q.QueryResults(taskID)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Println(task)
}
