package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"stariboss-script/env"
	"testing"
)

func TestPost(t *testing.T) {

}

func TestGet(t *testing.T) {

}

func TestChurnRequest(t *testing.T) {
	rInfo, err := json.Marshal(ChurnInfo{
		PartnerIds:    []string{"3882206"},
		CheckCanChurn: false,
	})
	if err != nil {
		panic(err)
	}

	bytes := Post(environment{
		Url:    env.PartnerChurn.ChurnUrl,
		Cookie: env.PartnerChurn.Cookie,
	}, rInfo)

	log.Info("result ", string(bytes))
}

type ChurnInfo struct {
	PartnerIds    []string `json:"partnerIds,omitempty"`
	CheckCanChurn bool     `json:"checkCanChurn"`
}
