package fileparse

import (
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

// Excel 中每行的数据
type ExcelRowData struct {
	Type       string `json:"type" xlsx:"记录类型"`
	Host       string `json:"host" xlsx:"主机记录"`
	ISPLine    string `json:"isp" xlsx:"解析线路"`
	Value      string `json:"value" xlsx:"记录值"`
	MXPriority string `json:"mxPriority" xlsx:"Mx优先级"`
	TTL        string `json:"ttl" xlsx:"TTL值"`
	Status     string `json:"status" xlsx:"状态(暂停/正常)"`
	Remark     string `json:"remark" xlsx:"备注"`
	ID         string `json:"id" xlsx:"ID"`
}

// Excel 中的数据
type ExcelData struct {
	Rows []ExcelRowData `json:"data"`
}

func NewExcelData(file string, domainName string) (*ExcelData, error) {
	var ed ExcelData
	f, err := excelize.OpenFile(file)
	if err != nil {
		logrus.Fatalf("打开 Excel 文件异常，原因: %v", err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logrus.Errorln(err)
			return
		}
	}()

	// 逐行读取Excel文件
	rows, err := f.GetRows(domainName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file":  file,
			"sheet": domainName,
		}).Errorf("读取中sheet页异常: %v", err)
		return nil, err
	}

	for k, row := range rows {
		// 跳过第一行
		if k == 0 {
			continue
		}
		logrus.WithFields(logrus.Fields{
			"k":   k,
			"row": row,
		}).Debugf("检查每一条需要处理的解析记录")

		// 尝试第8列的值，若无法获取则设置为空
		var id string

		if len(row) > 8 {
			id = row[8]
		} else {
			id = ""
		}

		// 将每一行中的的每列数据赋值到结构体中
		ed.Rows = append(ed.Rows, ExcelRowData{
			Type:       row[0],
			Host:       row[1],
			ISPLine:    row[2],
			Value:      row[3],
			MXPriority: row[4],
			TTL:        row[5],
			Status:     row[6],
			Remark:     row[7],
			ID:         id,
		})

	}

	return &ExcelData{
		Rows: ed.Rows,
	}, nil
}
