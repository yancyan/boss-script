package main

import (
	"encoding/csv"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"stariboss-script/env"
	"sync"
)

const FilePath = "file/加蓬_20220430_待流失代理商列表.csv"

func init() {
	initLogMethod("partner-churn")
	//	测试授权
	// page=1&limit=25&count=true&companyId=2&name=&start=0
	bytes := Get(environment{
		Url:    env.PartnerChurn.TestQueryPartnerTypeUrl,
		Cookie: env.PartnerChurn.Cookie,
	}, map[string]string{
		"page":      "1",
		"limit":     "25",
		"companyId": "2",
		"start":     "0",
	})
	log.Info("result is " + string(bytes))
	log.Info("test send request success.")
}

func churn(filePath string) {

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
	log.Infof("read data is %s", data)

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
	for _, info := range infos {
		waitGroup.Add(1)
		go churnRequest(info, &waitGroup)
	}

	waitGroup.Wait()

}

func churnRequest(info DataInfo, group *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	defer group.Done()

	log.Infof("request params is %+v", info)

	rInfo, err := json.Marshal(ChurnRequest{
		PartnerIds:    []string{info.PartnerId},
		CheckCanChurn: true,
	})
	if err != nil {
		panic(err)
	}

	resultBytes := Post(environment{
		Url:    env.PartnerChurn.ChurnUrl,
		Cookie: env.PartnerChurn.Cookie,
	}, rInfo)

	log.Info("result ", string(resultBytes))
}

type DataInfo struct {
	Line      int
	CompanyId string
	PartnerId string
}

type ChurnRequest struct {
	PartnerIds    []string `json:"partnerIds,omitempty"`
	CheckCanChurn bool     `json:"checkCanChurn"`
}
