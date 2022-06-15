package queryresults

import (
	"testing"

	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	"github.com/alibabacloud-go/tea-utils/service"
)

func TestAlidnsQueryResults_QueryResults(t *testing.T) {
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
		taskID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int32
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: domainName,
			fields: fields{
				AlidnsHandler: &alidns.AlidnsHandler{
					DomainName: domainName,
					Runtime:    &service.RuntimeOptions{},
					Client:     client,
				},
			},
			args: args{
				taskID: 3715007136,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := &AlidnsQueryResults{
				AlidnsHandler: tt.fields.AlidnsHandler,
			}
			got, err := qr.QueryResults(tt.args.taskID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlidnsQueryResults.QueryResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AlidnsQueryResults.QueryResults() = %v, want %v", got, tt.want)
			}
		})
	}
}
