package resolve

import (
	"testing"

	"github.com/DesistDaydream/aliyun-openapi/pkg/aliclient"
	"github.com/DesistDaydream/aliyun-openapi/pkg/alidns"
	"github.com/DesistDaydream/aliyun-openapi/pkg/config"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	"github.com/sirupsen/logrus"
)

func TestAlidnsResolve_DomainRecordsList(t *testing.T) {
	type fields struct {
		AlidnsHandler *alidns.AlidnsHandler
	}
	tests := []struct {
		name    string
		fields  fields
		want    *alidns20150109.DescribeDomainRecordsResponseBodyDomainRecords
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "test",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userName := "a"
			region := "alidns.cn-beijing.aliyuncs.com"
			domainName := "b"
			auth := config.NewAuthInfo("owner.yaml")
			aliclient.Info = &aliclient.ClientInfo{
				AK:     auth.AuthList[userName].AccessKeyID,
				SK:     auth.AuthList[userName].AccessKeySecret,
				Region: region,
			}

			h := alidns.NewAlidnsHandler(aliclient.Info.AK, aliclient.Info.SK, domainName, aliclient.Info.Region)
			r := NewAlidnsResolve(h)

			domainRecords, err := r.DomainRecordsList()
			if err != nil {
				logrus.Fatalln(err)
			}

			for _, r := range domainRecords.Record {
				logrus.WithFields(logrus.Fields{
					"domainName": *r.DomainName,
					"rr":         *r.RR,
				}).Infoln("域名记录")
			}

			logrus.Infof("共有 %d 条记录", len(domainRecords.Record))
		})
	}
}
