package models

// sent by the web-browser
type RegistrationInfoResponse struct {
	Data RegistrationInfo `json:"data"`
}

type RegistrationInfo struct {
	BrowserName         string `json:"browserName"`
	BrowserEngine       string `json:"browserEngine"`
	BrowserVersion      string `json:"browserVersion"`
	ProfileId           string `json:"profileId"`
	ProfileRank         int    `json:"profileRank"`
	ProfileName         string `json:"profileName"`
	ProfileAlias        string `json:"profileAlias"`
	ProfileCommandAlias string `json:"profileCommandAlias"`
	UserAgent           string `json:"userAgent"`
	RegisteredAt        string `json:"registeredAt"`
}
