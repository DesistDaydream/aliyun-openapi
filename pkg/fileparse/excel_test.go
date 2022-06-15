package fileparse

import (
	"reflect"
	"testing"
)

func TestNewExcelData(t *testing.T) {
	type args struct {
		file       string
		domainName string
	}
	tests := []struct {
		name string
		args args
		want *ExcelData
	}{
		// TODO: Add test cases.
		{
			name: "测试",
			args: args{
				file:       "desistdaydream.ltd.xlsx",
				domainName: "desistdaydream.ltd",
			},
			want: &ExcelData{
				Rows: []ExcelRowData{
					{
						Type:       "",
						Host:       "",
						ISPLine:    "",
						Value:      "",
						MXPriority: "",
						TTL:        "",
						Status:     "",
						Remark:     "",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExcelData(tt.args.file, tt.args.domainName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExcelData() = %v, want %v", got, tt.want)
			}
		})
	}
}
