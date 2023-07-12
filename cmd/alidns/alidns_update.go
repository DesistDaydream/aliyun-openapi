package alidns

import (
	"github.com/DesistDaydream/aliyun-openapi/pkg/fileparse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func alidnsUpdateCommand() *cobra.Command {
	alidnsUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "",
		Run:   runAlidnsUpdate,
	}

	return alidnsUpdateCmd
}

func runAlidnsUpdate(cmd *cobra.Command, args []string) {
	// domainRecords, err := r.DomainRecordsList()
	// if err != nil {
	// 	logrus.Fatalf("列出记录失败，原因: %v", err)
	// }

	excelData, err := fileparse.NewExcelData(alidnsFlags.rrFile, alidnsFlags.domainName)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	for _, row := range excelData.Rows {
		if row.ID != "" {
			updateDomainRecordResponse, err := h.Client.UpdateDomainRecordWithOptions(&alidns20150109.UpdateDomainRecordRequest{
				RecordId: tea.String(row.ID),
				RR:       tea.String(row.Host),
				Type:     tea.String(row.Type),
				Value:    tea.String(row.Value),
			}, &util.RuntimeOptions{})
			if err != nil && tea.ToMap(err)["Code"] != "DomainRecordDuplicate" {
				logrus.Errorf("更新记录失败，原因: %v", tea.ToMap(err)["Code"])
			}

			logrus.Infof("%v", updateDomainRecordResponse)
		} else {
			addDomainRecordResponse, err := h.Client.AddDomainRecordWithOptions(&alidns20150109.AddDomainRecordRequest{
				RR:         tea.String(row.Host),
				Type:       tea.String(row.Type),
				Value:      tea.String(row.Value),
				DomainName: tea.String(alidnsFlags.domainName),
			}, &util.RuntimeOptions{})
			if err != nil && tea.ToMap(err)["Code"] != "DomainRecordDuplicate" {
				logrus.Errorf("添加记录失败，原因: %v", tea.ToMap(err)["Code"])
			}

			logrus.Infof("%v", addDomainRecordResponse)
		}
	}
}
