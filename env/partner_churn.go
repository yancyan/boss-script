package env

const (
	partnerChurnTestQueryPartnerTypeUrl = "/partner-ui/api/partnerType/by-condition"
	partnerChurnChurnUrl                = "/partner-ui/api/partnerLoss/patchLoss"
)

type partnerChurn struct {
	TestQueryPartnerTypeUrl string
	ChurnUrl                string
	Cookie                  string
}

func initPartnerChurn() partnerChurn {
	return partnerChurn{
		TestQueryPartnerTypeUrl: Config.Domain + partnerChurnTestQueryPartnerTypeUrl,
		ChurnUrl:                Config.Domain + partnerChurnChurnUrl,
		Cookie:                  Config.Cookie,
	}
}
