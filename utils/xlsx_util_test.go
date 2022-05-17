package utils

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"os"
)

type xlsxUtil struct {
}

var XlsxUtil = xlsxUtil{}

func (xl xlsxUtil) ReadAllData(filePath string) [][]string {

	log.Info("parse file " + filePath)

	open, err := os.Open(filePath)
	if err != nil {
		panic("open " + filePath + " error " + err.Error())
	}

	reader := csv.NewReader(open)
	if reader == nil {
		panic("reader nil.")
	}

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	log.Infof("parse data is %s", data)
	return data
}

func (xl *xlsxUtil) newXlsxFile(fileName string) {
	f := excelize.NewFile()
	// 创建一个工作表
	index := f.NewSheet(fileName)
	// 设置单元格的值
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs("/result/" + fileName + ".xlsx"); err != nil {
		log.Errorln(err)
	}
}
