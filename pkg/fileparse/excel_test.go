package fileparse

import (
	"os"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

func TestNewExcelData(t *testing.T) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"类型", "记录", "值", "ID", "描述"})

	file := "/mnt/d/tmp/superstor.cn.xlsx"
	domainName := "superstor.cn"

	got, err := NewExcelData(file, domainName)
	if err != nil {
		logrus.Errorln(err)
	}

	for _, row := range got.Rows {
		table.Append([]string{row.Type, row.RR, row.Value, row.ID, row.Remark})
	}

	table.Render()
}
