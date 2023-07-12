package alidns

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func alidnsListCommand() *cobra.Command {
	alidnsListCmd := &cobra.Command{
		Use:   "list",
		Short: "列出域名下的全部解析记录",
		Run:   runAlidnsList,
	}

	return alidnsListCmd
}

func runAlidnsList(cmd *cobra.Command, args []string) {
	domainRecords, err := r.DomainRecordsList()
	if err != nil {
		logrus.Fatalf("列出记录失败，原因: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"类型", "记录", "值", "ID", "描述"})

	for _, domainRecord := range domainRecords.Record {
		var remark string
		if domainRecord.Remark != nil {
			remark = *domainRecord.Remark
		} else {
			remark = ""
		}
		table.Append([]string{*domainRecord.Type, *domainRecord.RR, *domainRecord.Value, *domainRecord.RecordId, remark})
	}

	table.Render()

	logrus.Infof("共有 %d 条记录", len(domainRecords.Record))
}
