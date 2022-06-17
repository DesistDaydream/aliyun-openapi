// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/domain"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/queryresults"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns/resolve"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/DesistDaydream/aliyun-openapi/pkg/logging"
)

// LogInit 日志功能初始化，若指定了 log-output 命令行标志，则将日志写入到文件中
func LogInit(level, file, format string) error {
	switch format {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   "2006-01-02 15:04:05",
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
			DataKey:           "",
			// FieldMap:          map[logrus.fieldKey]string{},
			// CallerPrettyfier: func(*runtime.Frame) (string, string) {},
			PrettyPrint: false,
		})
	}

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)

	if file != "" {
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		logrus.SetOutput(f)
	}

	return nil
}

func main() {
	operation := pflag.StringP("operation", "o", "", "操作类型: [add, list, batch]")
	batchOperation := pflag.StringP("batch-operation", "O", "", "批量操作类型: [RR_ADD,RR_DEL,DOMAIN_ADD,DOMAIN_DEL]")

	// 添加命令行标志
	alidnsFlags := &alidns.AlidnsFlags{}
	logFlags := logging.LoggingFlags{}
	logFlags.AddFlags()
	alidnsFlags.AddFlags()
	pflag.Parse()

	// 初始化日志
	if err := LogInit(logFlags.LogLevel, logFlags.LogOutput, logFlags.LogFormat); err != nil {
		logrus.Fatal("set log level error")
	}

	if alidnsFlags.DomainName == "" {
		logrus.Fatal("请使用 -d 标志指定要操作的域名")
	}

	// 获取认证信息
	auth := config.NewAuthInfo(alidnsFlags.AuthFile)

	// 判断传入的域名是否存在在认证信息中
	if !auth.IsDomainExist(alidnsFlags.DomainName) {
		logrus.Fatalf("认证信息中不存在 %v 域名, 请检查认证信息文件或命令行参数的值", alidnsFlags.DomainName)
	}

	// 初始化账号Client
	client, err := aliclient.CreateClient(auth.AuthList[alidnsFlags.DomainName].AccessKeyID, auth.AuthList[alidnsFlags.DomainName].AccessKeySecret)
	if err != nil {
		panic(err)
	}

	// 实例化各种 API 处理器
	h := alidns.NewAlidnsHandler(alidnsFlags, client)
	d := domain.NewAlidnsDomain(h)
	r := resolve.NewAlidnsResolve(h)
	q := queryresults.NewQueryResults(h)

	switch *operation {
	case "list":
		r.DomainRecordsList()
	case "add":
		// 检查文件是否存在
		if _, err := os.Stat(alidnsFlags.RRFile); os.IsNotExist(err) {
			logrus.Fatal("文件不存在")
		}

		// 添加解析记录前先删除全部解析记录
		taskID, err := d.Batch("RR_DEL", alidnsFlags.RRFile)
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
				break
			}
			time.Sleep(time.Second * 1)
		}
		// 这里为了测试目的，使用的是逐一添加的方式，实际应用中可以使用批量添加的方式
		r.OnebyoneAddDomainRecord(alidnsFlags.RRFile)
	case "batch":
		// 检查文件是否存在
		if _, err := os.Stat(alidnsFlags.RRFile); os.IsNotExist(err) {
			logrus.Fatal("文件不存在")
		}
		// 判断批量操作类型是否存在
		if *batchOperation == "" {
			logrus.Fatal("请使用 -O 标志指定批量操作类型")
		}
		// 判断批量操作类型是否合法
		if !d.IsBatchOperationExist(*batchOperation) {
			logrus.Fatal("批量操作类型不存在，可用的值有: RR_ADD,RR_DEL,DOMAIN_ADD,DOMAIN_DEL")
		}
		d.Batch(*batchOperation, alidnsFlags.RRFile)
	default:
		logrus.Fatalln("操作类型不存在，请使用 -o 指定操作类型")
	}
}
