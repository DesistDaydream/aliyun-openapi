package domain

import (
	"testing"

	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/sirupsen/logrus"
)

func TestAlidnsDomain_Batch(t *testing.T) {
	// 准备测试数据
	domainName := "desistdaydream.ltd"
	auth := config.NewAuthInfo("../../../auth.yaml")
	handler := alidns.NewAlidnsHandler(auth, "断灬念梦", domainName, "alidns.cn-beijing.aliyuncs.com")

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
			d := &AlidnsDomain{
				AlidnsHandler: tt.fields.AlidnsHandler,
			}
			got, err := d.Batch(tt.args.operateType, tt.args.file)
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
