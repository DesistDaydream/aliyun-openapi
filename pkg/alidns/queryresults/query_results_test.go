package queryresults

import (
	"testing"

	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
)

func TestAlidnsQueryResults_QueryResults(t *testing.T) {
	// 准备测试数据
	domainName := "desistdaydream.ltd"
	auth := config.NewAuthInfo("../../../auth.yaml")
	handler := alidns.NewAlidnsHandler(auth, "断灬念梦", domainName, "alidns.cn-beijing.aliyuncs.com")

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
				AlidnsHandler: handler,
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
