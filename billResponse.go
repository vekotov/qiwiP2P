package qiwiP2P

type BillResponse struct {
	SiteId             string            `json:"siteId"`
	BillId             string            `json:"billId"`
	Amount             Amount            `json:"amount"`
	Status             Status            `json:"status"`
	Customer           Customer          `json:"customer"`
	CustomFields       map[string]string `json:"customFields"`
	Comment            string            `json:"comment"`
	CreationDateTime   string            `json:"creationDateTime"`
	ExpirationDateTime string            `json:"expirationDateTime"`
	PayUrl             string            `json:"payUrl"`
}

type RequestError struct {
	ServiceName string `json:"serviceName"`
	ErrorCode   string `json:"errorCode"`
	Description string `json:"description"`
	UserMessage string `json:"userMessage"`
	DateTime    string `json:"dateTime"`
	TraceId     string `json:"traceId"`
}

type Status struct {
	Value           string `json:"value"`
	ChangedDateTime string `json:"changedDateTime"`
}
