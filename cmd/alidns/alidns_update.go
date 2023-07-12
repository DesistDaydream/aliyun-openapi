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
		Short: "更新解析记录，若存在记录ID则更新，不存在则创建",
		Run:   runAlidnsUpdate,
	}

	return alidnsUpdateCmd
}

func runAlidnsUpdate(cmd *cobra.Command, args []string) {
	checkFile(alidnsFlags.rrFile)

	excelData, err := fileparse.NewExcelData(alidnsFlags.rrFile, alidnsFlags.domainName)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	for _, row := range excelData.Rows {
		if row.ID != "" {
			resp, err := h.Client.UpdateDomainRecordWithOptions(&alidns20150109.UpdateDomainRecordRequest{
				RecordId: tea.String(row.ID),
				RR:       tea.String(row.Host),
				Type:     tea.String(row.Type),
				Value:    tea.String(row.Value),
			}, &util.RuntimeOptions{})
			if err != nil {
				if tea.ToMap(err)["Code"] == "DomainRecordDuplicate" {
					logrus.Warnf("%v 记录无变化，无需更新", row.Host)
					continue
				} else {
					logrus.Errorf("更新记录 %v 失败，原因: %v", row.Host, tea.ToMap(err)["Code"])
				}
			}

			logrus.Infof("ID 为 %v 的记录更新成功", *resp.Body.RecordId)
		} else {
			resp, err := h.Client.AddDomainRecordWithOptions(&alidns20150109.AddDomainRecordRequest{
				RR:         tea.String(row.Host),
				Type:       tea.String(row.Type),
				Value:      tea.String(row.Value),
				DomainName: tea.String(alidnsFlags.domainName),
			}, &util.RuntimeOptions{})
			if err != nil {
				if tea.ToMap(err)["Code"] == "DomainRecordDuplicate" {
					logrus.Warnf("%v 已存在相同记录，无需添加", row.Host)
					continue
				} else {
					logrus.Errorf("添加记录 %v 失败，原因: %v", row.Host, tea.ToMap(err)["Code"])
				}
			}

			logrus.Infof("成功添加 ID 为 %v 的记录", *resp.Body.RecordId)
		}
	}
}
