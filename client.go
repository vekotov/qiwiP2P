package qiwiP2P

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Client
// Object used for storing client data
// token - Private key of QIWI P2P
// client - http.Client object, used for http requests to API
// ch - Channel of payment updates, used for webhook
type Client struct {
	token  string
	client http.Client
	ch     chan PaymentUpdate
}

// CreateClient
// Method creates Client object with given private key
func CreateClient(key string) *Client {
	return &Client{token: key, client: http.Client{}, ch: make(chan PaymentUpdate, 50)}
}

// SetSecretKey
// Method changes private key of Client to new key
func (c *Client) SetSecretKey(key string) *Client {
	c.token = key
	return c
}

// PutBill
// Methods sends Bill object to API, putting it to random ID
// Returns BillResponse on success, error on failed
func (c *Client) PutBill(b *Bill) (result *BillResponse, err error) {
	billId := pseudoUUID()
	res, code, err := c.makeRequest(
		billId,
		"PUT",
		b.toJSON(),
	)

	if err != nil {
		return nil, err
	}

	if code == 400 {
		return nil, RequestError{
			ErrorCode:   "bad_request",
			Description: "Bad request. Maybe you have bad expire time?",
		}
	}
	if code == 401 {
		return nil, RequestError{
			ErrorCode:   "bad_token",
			Description: "Bad token",
		}
	}
	if code != 200 {
		return nil, &RequestError{
			ErrorCode:   "Error " + strconv.Itoa(code),
			Description: "HTTP Error " + strconv.Itoa(code),
		}
	}

	return parseResponse(res)
}

// GetBill
// Method gets info about bill with given id
// You can get ID from BillResponse, which was returned in PutBill
// Returns BillResponse on success, error on failed
func (c *Client) GetBill(id string) (result *BillResponse, err error) {
	res, code, err := c.makeRequest(
		id,
		"GET",
		"",
	)

	if err != nil {
		return nil, err
	}

	if code == 401 {
		return nil, &RequestError{
			ErrorCode:   "bad_token",
			Description: "Bad token",
		}
	}
	if code == 404 {
		return nil, &RequestError{
			ErrorCode:   "bad_id",
			Description: "No such bill found",
		}
	}
	if code != 200 {
		return nil, &RequestError{
			ErrorCode:   "Error " + strconv.Itoa(code),
			Description: "HTTP Error " + strconv.Itoa(code),
		}
	}

	return parseResponse(res)
}

// RejectBill
// Method sets bill status to rejected. Id is needed in arguments
// You can get ID from BillResponse, which was returned in PutBill
// Returns BillResponse on success, error on failed
func (c *Client) RejectBill(id string) (result *BillResponse, err error) {
	res, code, err := c.makeRequest(
		id+"/reject",
		"POST",
		"",
	)

	if err != nil {
		return nil, err
	}

	if code == 401 {
		return nil, &RequestError{
			ErrorCode:   "bad_token",
			Description: "Bad token",
		}
	}
	if code == 404 {
		return nil, &RequestError{
			ErrorCode:   "bad_id",
			Description: "No such bill found",
		}
	}
	if code != 200 {
		return nil, &RequestError{
			ErrorCode:   "Error " + strconv.Itoa(code),
			Description: "HTTP Error " + strconv.Itoa(code),
		}
	}

	return parseResponse(res)
}

// parseResponse
// Parses BillResponse or RequestError JSON
// Returns either BillResponse or error object
func parseResponse(jsonResponse string) (result *BillResponse, error error) {
	var re RequestError
	err := json.Unmarshal([]byte(jsonResponse), &re)
	if err != nil {
		return nil, err
	}
	if re.ErrorCode != "" {
		return nil, re
	}

	var response BillResponse
	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// makeRequest
// Makes HTTP request to QIWI API server
// url - path to needed API method
// method - HTTP method used in API call
// data - JSON body data (for POST and PUT requests)
func (c *Client) makeRequest(url string, method string, data string) (json string, code int, err error) {
	url = "https://api.qiwi.com/partner/bill/v1/bills/" + url
	req, err := http.NewRequest(
		method, url, strings.NewReader(data),
	)
	if err != nil {
		return "", 0, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("content-type", "application/json")
	res, _ := c.client.Do(req)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return "", 0, err
	}
	return buf.String(), res.StatusCode, nil
}

// pseudoUUID
// Generates random combination of symbols and letters
// Used as bill ID
func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

// StartWebhook
// Starts webhook listening on given path and port.
// Returns channel with payment updates
// Usage:
// ```
// ch := c.StartWebhook("/qiwiWebhook", 80)
// ```
func (c *Client) StartWebhook(path string, port int) chan PaymentUpdate {
	go c.startListening(path, port)
	return c.ch
}

// startListening
// Starts webhook listening. Pauses current thread
func (c *Client) startListening(path string, port int) {
	http.HandleFunc(path, c.onWebhook)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

// onWebhook
// Called when new webhook request is caught
func (c *Client) onWebhook(w http.ResponseWriter, r *http.Request) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		log.Println("Error while webhook: ", err.Error())
		return
	}

	var upd PaymentUpdate
	err = json.Unmarshal([]byte(buf.String()), &upd)
	if err != nil {
		log.Println("Error while webhook: ", err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write([]byte("{\"error\":\"0\"}"))
	if err != nil {
		log.Println("Error while webhook: ", err.Error())
		return
	}
	if c.verifyWebhook(upd, r.Header.Get("X-Api-Signature-SHA256")) {
		c.ch <- upd
	}
}

// verifyWebhook
// Method verifies update, returns true if update is authorized and false if not
// update - PaymentUpdate object
// hash - X-Api-Signature-SHA256 header from webhook
func (c *Client) verifyWebhook(update PaymentUpdate, hash string) bool {
	invoiceParameters := ""
	invoiceParameters += update.Bill.Amount.Currency + "|"
	invoiceParameters += update.Bill.Amount.Value + "|"
	invoiceParameters += update.Bill.BillId + "|"
	invoiceParameters += update.Bill.SiteId + "|"
	invoiceParameters += update.Bill.Status.Value
	h := hmac.New(sha256.New, []byte(c.token))
	h.Write([]byte(invoiceParameters))
	sha := hex.EncodeToString(h.Sum(nil))
	if sha == hash {
		return true
	}
	return false
}
