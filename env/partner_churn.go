package env

var PartnerChurn = partnerChurn{
	TestQueryPartnerTypeUrl: partnerChurnTestQueryPartnerTypeUrl,
	ChurnUrl:                partnerChurnChurnUrl,
	Cookie:                  partnerChurnCookie,
}

type partnerChurn struct {
	TestQueryPartnerTypeUrl string
	ChurnUrl                string
	Cookie                  string
}
