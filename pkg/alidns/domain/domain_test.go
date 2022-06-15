package domain

import (
	"testing"

	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/alibabacloud-go/tea-utils/service"
)

func TestAlidnsDomain_Batch(t *testing.T) {
	// 准备测试数据
	domainName := "desistdaydream.ltd"
	auth := config.NewAuthInfo("../../../auth.yaml")
	client, err := aliclient.CreateClient(auth.AuthList[domainName].AccessKeyID, auth.AuthList[domainName].AccessKeySecret)
	if err != nil {
		panic(err)
	}

	// 使用 gtotests 工具生成的测试代码
	type fields struct {
		AlidnsHandler *alidns.AlidnsHandler
	}
	type args struct {
		operateType string
		file        string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "测试批量操作",
			fields: fields{
				AlidnsHandler: &alidns.AlidnsHandler{
					DomainName: domainName,
					Runtime:    &service.RuntimeOptions{},
					Client:     client,
				},
			},
			args: args{
				operateType: "RR_ADD",
				file:        "../../../desistdaydream.ltd.xlsx",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AlidnsDomain{
				AlidnsHandler: tt.fields.AlidnsHandler,
			}
			d.Batch(tt.args.operateType, tt.args.file)
		})
	}
}
