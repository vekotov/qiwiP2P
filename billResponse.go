package qiwiP2P

// BillResponse
// Object used for storing response about bill from API
//
// Fields:
//   - `SiteId` : id of seller
//   - `BillId` : id of bill
//   - `Amount` : amount of money customer needs to pay
//   - `Status` : status of bill (paid/waiting/rejected etc.)
//   - `Customer` : info about customer
//   - `CustomFields` : custom additional fields with application data
//   - `Comment` : comment of bill (displayed to user)
//   - `CreationDateTime` : timestamp of bill creation moment (in format 2006-01-02T15:04:05+00:00)
//   - `ExpirationDateTime` : timestamp of bill expiration moment (in format 2006-01-02T15:04:05+00:00)
//   - `PayUrl` : link to payment form, should be sent to customer
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

// Status
// Object used for storing bill status
//
// Fields:
//   - `Value` : status of bill ("PAID"/"WAITING"/"REJECTED"/"EXPIRED")
//   - `ChangedDateTime` : timestamp of moment, when bill status was changed last time (in format 2006-01-02T15:04:05+00:00)
type Status struct {
	Value           string `json:"value"`
	ChangedDateTime string `json:"changedDateTime"`
}

// PaymentUpdate
// Object used for storing webhook updates
//
// Fields:
//   - `Bill` : object of updated bill
//   - `Version` : version of Webhook API
type PaymentUpdate struct {
	Bill    BillResponse `json:"bill"`
	Version string       `json:"version"`
}

// RequestError
// Object used for returning errors, implements error interface
//
// Fields:
//   - `ServiceName` : name of service, where error happened
//   - `ErrorCode` : string code of error
//   - `Description` : description of error
//   - `UserMessage` : description of error, which can be showed to user/customer
//   - `DateTime` : timestamp of moment, when error was happened (in format 2006-01-02T15:04:05+00:00)
//   - `TraceId` : id of error, can be sent to tech support of QIWI P2P
type RequestError struct {
	ServiceName string `json:"serviceName"`
	ErrorCode   string `json:"errorCode"`
	Description string `json:"description"`
	UserMessage string `json:"userMessage"`
	DateTime    string `json:"dateTime"`
	TraceId     string `json:"traceId"`
}

// Error
// Method returns description of happened error
func (c RequestError) Error() string {
	return c.Description
}
