package file

import (
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

// Excel 文件中的信息
type ExcelRowData struct {
	Type       string `json:"type"`
	Host       string `json:"host"`
	ISPLine    string `json:"isp"`
	Value      string `json:"value"`
	MXPriority string `json:"mxPriority"`
	TTL        string `json:"ttl"`
	Status     string `json:"status"`
	Remark     string `json:"remark"`
}

type ExcelData struct {
	Data []ExcelRowData `json:"data"`
}

func NewExcelInfo(file string, domainName string) (ed *ExcelData) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logrus.Errorln(err)
		}
	}()

	// 从第二行开始逐行读取Excel文件
	rows, _ := f.GetRows(domainName)

	for k, row := range rows {
		if k == 0 {
			continue
		}
		excelInfo := new(ExcelRowData)
		// 读取每一行的数据
		excelInfo.Type = row[0]
		excelInfo.Host = row[1]
		excelInfo.ISPLine = row[2]
		excelInfo.Value = row[3]
		excelInfo.MXPriority = row[4]
		excelInfo.TTL = row[5]
		excelInfo.Status = row[6]
		// 尝试获取7号元素，若无法获取则设置为空
		if len(row) > 7 {
			excelInfo.Remark = row[7]
		} else {
			excelInfo.Remark = ""
		}

		ed.Data = append(ed.Data, *excelInfo)
	}

	return &ExcelData{
		Data: ed.Data,
	}
}

// 处理 Excel 文件
func HandleExcel(fileName string, domainName string) (excelInfos []*ExcelRowData) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logrus.Errorln(err)
		}
	}()

	// 从第二行开始逐行读取Excel文件
	rows, _ := f.GetRows(domainName)

	for k, row := range rows {
		if k == 0 {
			continue
		}
		excelInfo := new(ExcelRowData)
		// 读取每一行的数据
		excelInfo.Type = row[0]
		excelInfo.Host = row[1]
		excelInfo.ISPLine = row[2]
		excelInfo.Value = row[3]
		excelInfo.MXPriority = row[4]
		excelInfo.TTL = row[5]
		excelInfo.Status = row[6]
		// 尝试获取7号元素，若无法获取则设置为空
		if len(row) > 7 {
			excelInfo.Remark = row[7]
		} else {
			excelInfo.Remark = ""
		}

		excelInfos = append(excelInfos, excelInfo)
	}

	return excelInfos
}
