package alidns

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func alidnsListCommand() *cobra.Command {
	alidnsListCmd := &cobra.Command{
		Use:   "list",
		Short: "",
		Run:   runAlidnsList,
	}

	return alidnsListCmd
}

func runAlidnsList(cmd *cobra.Command, args []string) {
	domainRecords, err := r.DomainRecordsList()
	if err != nil {
		panic(err)
	}

	for index, domainRecord := range domainRecords.Record {
		logrus.WithFields(logrus.Fields{
			"类型": *domainRecord.Type,
			"记录": *domainRecord.RR,
			"值":  *domainRecord.Value,
		}).Infof("%v 域名的第 %v 条资源记录", alidnsFlags.domainName, index+1)
	}
	logrus.Infof("共有 %d 条记录", len(domainRecords.Record))
}
