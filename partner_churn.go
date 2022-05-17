package main

import (
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"stariboss-script/env"
	"stariboss-script/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

func authTest() {
	//	测试授权
	utils.RequestUtil.Get(utils.RequestContext{
		Url:    env.Config.Domain + "/partner-ui/api/partnerType/by-condition",
		Cookie: env.Config.PartnerChurn.Cookie,
	}, map[string]string{
		"page":      "1",
		"limit":     "25",
		"companyId": "1",
		"start":     "0",
	})
	log.Info("test send request success.")
}

func churn(filePath string, resultPath string) {
	authTest()

	data := utils.XlsxUtil.ReadAllData(filePath)

	var infos []DataInfo
	for i, datum := range data {
		if i == 0 {
			continue
		}
		log.Infof("%d info %s \n", i, datum)
		infos = append(infos, DataInfo{
			Line:      i,
			CompanyId: datum[0],
			PartnerId: datum[1],
		})
	}

	var waitGroup sync.WaitGroup

	utils.XlsxUtil.NewXlsxFile(resultPath)
	utils.XlsxUtil.AppendRowData(resultPath, []interface{}{"line", "company_id", "agent_id", "result", "reason"})

	waitGroup.Add(len(infos))
	count := 0
	for _, info := range infos {
		churnRequest(info, resultPath, &waitGroup)
		count++
		if count%20 == 0 {
			time.Sleep(5 * time.Second)
		}
	}

	waitGroup.Wait()

}

func churnNewAutoDir(fileName string) {
	churnNew("file/"+fileName, "Result_"+fileName)
}

func churnNew(filePath string, resultPath string) {
	resultPath = strings.Replace(resultPath, ".csv", ".xlsx", 1)
	authTest()
	data := utils.XlsxUtil.ReadAllData(filePath)
	var infos []DataInfo
	for i, datum := range data {
		if i == 0 {
			continue
		}
		log.Infof("%d info %s \n", i, datum)
		infos = append(infos, DataInfo{
			Line:      i,
			CompanyId: datum[0],
			PartnerId: datum[1],
		})
	}

	var waiter sync.WaitGroup

	utils.XlsxUtil.NewXlsxFile(resultPath)

	waiter.Add(len(infos))

	var resultData [][]interface{}
	resultData = append(resultData, []interface{}{"line", "company_id", "agent_id", "result", "reason"})

	var lock sync.Mutex

	churnPool, _ := ants.NewPoolWithFunc(5, func(i interface{}) {
		defer waiter.Done()
		param, ok := i.(ChurnParam)
		if !ok {
			panic("the param cannot parse to ChurnParam")
		}
		resultD := churnRequestNew(param)
		lock.Lock()
		resultData = append(resultData, resultD)
		lock.Unlock()

	})

	for _, info := range infos {
		err := churnPool.Invoke(ChurnParam{
			info:       info,
			resultPath: resultPath,
		})
		if err != nil {
			log.Errorln(err)
		}

	}

	waiter.Wait()
	if resultData != nil && len(resultData) > 0 {
		log.Info("start write result. size is", len(resultData))

		for index, datum := range resultData {
			if index%5 == 0 || index == len(resultData)-1 {
				log.Info("write line", index)
			}
			utils.XlsxUtil.AppendRowData(resultPath, datum)
		}
	}

	log.Infof("end.")
}

func churnRequest(info DataInfo, resultPath string, group *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	defer group.Done()

	log.Infof("the request params is %+v", info)

	rInfo, err := json.Marshal(ChurnRequest{
		PartnerIds:    []string{info.PartnerId},
		CheckCanChurn: true,
	})
	if err != nil {
		panic(err)
	}

	resultBytes, er := utils.RequestUtil.Post(utils.RequestContext{
		Url:    env.Config.PartnerChurn.ChurnUrl,
		Cookie: env.Config.PartnerChurn.Cookie,
	}, rInfo)
	log.Info("the resp body is ", string(resultBytes))
	if er != nil {
		utils.XlsxUtil.AppendRowData(resultPath, []interface{}{info.Line, info.CompanyId, info.PartnerId, false, er.Error()})
	} else {

		result := ChurnResult{}
		_ = json.Unmarshal(resultBytes, &result)

		for _, prob := range result.ProblemList {
			utils.XlsxUtil.AppendRowData(resultPath, []interface{}{info.Line, info.CompanyId, info.PartnerId, false, prob.Reason})
		}

		partnerId, _ := strconv.Atoi(info.PartnerId)
		contains, _ := inArray(partnerId, result.RightIds)
		if contains {
			utils.XlsxUtil.AppendRowData(resultPath, []interface{}{info.Line, info.CompanyId, info.PartnerId, true, ""})
		} else if len(result.ProblemList) == 0 {
			utils.XlsxUtil.AppendRowData(resultPath, []interface{}{info.Line, info.CompanyId, info.PartnerId, false,
				fmt.Sprintf("error, the rights{%v} cannot contains partner_id[%d]", result.RightIds, partnerId)})
		}
	}

}

type ChurnParam struct {
	info       DataInfo
	resultPath string
}

type DataInfo struct {
	Line      int
	CompanyId string
	PartnerId string
}

type ChurnRequest struct {
	PartnerIds    []string `json:"partnerIds"`
	CheckCanChurn bool     `json:"checkCanChurn"`
}

type ChurnResult struct {
	RightIds    []int         `json:"rightIds"`
	ProblemList []ProblemInfo `json:"problemList"`
}

type ProblemInfo struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func inArray(val int, array []int) (exists bool, index int) {
	exists = false
	index = -1

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}
	return
}

func churnRequestNew(param ChurnParam) []interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()

	info := param.info

	log.Infof("the request params is %+v", param.info)

	rInfo, err := json.Marshal(ChurnRequest{
		PartnerIds:    []string{info.PartnerId},
		CheckCanChurn: true,
	})
	if err != nil {
		panic(err)
	}

	resultBytes, er := utils.RequestUtil.Post(utils.RequestContext{
		Url:    env.Config.PartnerChurn.ChurnUrl,
		Cookie: env.Config.PartnerChurn.Cookie,
	}, rInfo)
	log.Info("the resp body is ", string(resultBytes))
	if er != nil {
		return []interface{}{info.Line, info.CompanyId, info.PartnerId, false, er.Error()}
	} else {

		result := ChurnResult{}
		_ = json.Unmarshal(resultBytes, &result)

		for _, prob := range result.ProblemList {
			return []interface{}{info.Line, info.CompanyId, info.PartnerId, false, prob.Reason}
		}

		partnerId, _ := strconv.Atoi(info.PartnerId)
		contains, _ := inArray(partnerId, result.RightIds)
		if contains && len(result.ProblemList) == 0 {
			return []interface{}{info.Line, info.CompanyId, info.PartnerId, true, ""}
		}
	}
	return nil
}
