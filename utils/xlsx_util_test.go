package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestXlsxUtil_NewXlsxFile(t *testing.T) {
	log.Info(fmt.Sprintf("the data %d id %d", []int{2, 4}, 12))
	XlsxUtil.fileDir = ""
	fileName := "churn_result_01.xlsx"
	XlsxUtil.NewXlsxFile(fileName)
	XlsxUtil.AppendRowData(fileName, []interface{}{"company_id", "agent_id", "result", "fail_reason"})
}

func TestXlsxUtil_AddRowData(t *testing.T) {
	XlsxUtil.fileDir = ""
	fileName := "churn_result_01.xlsx"
	//XlsxUtil.NewXlsxFile(fileName)
	//XlsxUtil.AppendRowData(fileName, []interface{}{"company_id", "agent_id", "result", "fail_reason"})
	//XlsxUtil.AppendRowData(fileName, []interface{}{"32", 12, true})
	//XlsxUtil.AppendRowData(fileName, []interface{}{"32", 12, true, ""})

	for i := 0; i < 50; i++ {
		go func() {
			XlsxUtil.AppendRowData(fileName, []interface{}{"32", 12, false, "fail content."})
		}()
	}

	time.Sleep(20 * time.Second)
}
