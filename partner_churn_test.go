package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"stariboss-script/env"
	"stariboss-script/utils"
	"testing"
)

func TestChurnRequest(t *testing.T) {
	rInfo, err := json.Marshal(ChurnInfo{
		PartnerIds:    []string{"3882206"},
		CheckCanChurn: false,
	})
	if err != nil {
		panic(err)
	}

	bytes, _ := utils.RequestUtil.Post(utils.RequestContext{
		Url:    env.Config.PartnerChurn.ChurnUrl,
		Cookie: env.Config.PartnerChurn.Cookie,
	}, rInfo)
	info := ChurnResult{}
	_ = json.Unmarshal(bytes, &info)

	log.Info("result ", string(bytes))
}

type ChurnInfo struct {
	PartnerIds    []string `json:"partnerIds,omitempty"`
	CheckCanChurn bool     `json:"checkCanChurn"`
}
