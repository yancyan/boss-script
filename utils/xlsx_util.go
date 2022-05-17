package utils

import (
	"encoding/csv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"os"
)

type xlsxUtil struct {
	fileDir string
	//currentRowNum int32
}

var XlsxUtil = xlsxUtil{
	fileDir: "result_file/",
	//currentRowNum: 1,
}

func (xl *xlsxUtil) ReadAllData(filePath string) [][]string {

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

func (xl *xlsxUtil) NewXlsxFile(fileName string) {
	f := excelize.NewFile()
	a, _ := f.GetColWidth("Sheet1", "A")
	e, _ := f.GetColWidth("Sheet1", "E")
	_ = f.SetColWidth("Sheet1", "A", "A", a*1.5)
	_ = f.SetColWidth("Sheet1", "E", "E", e*6)
	if err := f.SaveAs(xl.fileDir + fileName); err != nil {
		log.Errorln(err)
		panic(err)
	}
}

func (xl *xlsxUtil) AppendRowData(fileName string, rowData []interface{}) {
	appendRowData(xl.fileDir+fileName, "A", rowData)
}

//var lock = &sync.Mutex{}

func appendRowData(fileName string, starCol string, rowData []interface{}) {

	//lock.Lock()
	//defer lock.Unlock()

	f, err := excelize.OpenFile(fileName, excelize.Options{})
	if err != nil {
		panic(err)
	}
	defer func(f *excelize.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	getRows, err := f.GetRows("Sheet1")
	if err != nil {
		panic(err)
	}
	lastIndex := len(getRows)
	err = f.SetSheetRow("Sheet1", fmt.Sprintf(starCol+"%d", lastIndex+1), &rowData)
	if err != nil {
		panic(err)
	}
	er := f.Save()
	if er != nil {
		panic(er)
	}
}
