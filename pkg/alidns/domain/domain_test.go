package domain

import (
	"testing"

	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/DesistDaydream/aliyun-openapi/pkg/fileparse"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
)

func TestAlidnsDomain_Batch(t *testing.T) {
	// 准备测试数据
	domainName := "desistdaydream.ltd"
	auth := config.NewAuthInfo("../../../owner.yaml")
	ak := auth.AuthList["断灬念梦"].AccessKeyID
	sk := auth.AuthList["断灬念梦"].AccessKeySecret
	handler := alidns.NewAlidnsHandler(ak, sk, domainName, "alidns.cn-beijing.aliyuncs.com")

	// 使用 gtotests 工具生成的测试代码
	type fields struct {
		AlidnsHandler *alidns.AlidnsHandler
	}
	type args struct {
		operateType string
		file        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "测试批量操作",
			fields: fields{
				AlidnsHandler: handler,
			},
			args: args{
				operateType: "RR_ADD",
				file:        "../../../desistdaydream.ltd.xlsx",
			},
			// 这里不用 want 为 0，因为结果肯定是不为 0
			// want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domainRecordInfos, err := HandleFile(tt.args.file, domainName)
			if err != nil {
				panic(err)
			}

			d := &AlidnsDomain{
				AlidnsHandler: tt.fields.AlidnsHandler,
			}
			got, err := d.Batch(tt.args.operateType, domainRecordInfos)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlidnsDomain.Batch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == 0 {
				t.Errorf("AlidnsDomain.Batch() = %v, want %v", got, tt.want)
			}
			logrus.Infof("测试成功,执行了 %v 行为，获取到的任务ID 为: %v", tt.args.operateType, got)
		})
	}
}

func HandleFile(file string, domainName string) ([]*alidns20150109.OperateBatchDomainRequestDomainRecordInfo, error) {
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

	logrus.Debugln(domainRecordInfos)

	return domainRecordInfos, nil
}
