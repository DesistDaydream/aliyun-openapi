package fileparse

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewExcelData(t *testing.T) {
	type args struct {
		file       string
		domainName string
	}
	tests := []struct {
		name    string
		args    args
		want    *ExcelData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "测试",
			args: args{
				file:       "../../desistdaydream.ltd.xlsx",
				domainName: "desistdaydream.ltd",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewExcelData(tt.args.file, tt.args.domainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExcelData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, r := range got.Rows {
				logrus.Infoln(r)
			}
		})
	}
}
