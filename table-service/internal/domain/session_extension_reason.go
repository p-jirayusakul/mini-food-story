package domain

type ListSessionExtensionReason struct {
	ID       int64   `json:"id,string"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	NameEN   string  `json:"nameEN"`
	Category *string `json:"category"`
	ModeCode *string `json:"modeCode"`
}
