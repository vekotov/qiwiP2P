package qiwiP2P

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	token  string
	client http.Client
	ch     chan PaymentUpdate
}

func CreateClient(key string) *Client {
	return &Client{token: key, client: http.Client{}, ch: make(chan PaymentUpdate, 10)}
}

func (c *Client) SetSecretKey(key string) *Client {
	c.token = key
	return c
}

func (c *Client) PutBill(b *Bill) (result *BillResponse, err error) {
	billId := pseudoUUID()
	res, code, err := c.makeRequest(
		"https://api.qiwi.com/partner/bill/v1/bills/"+billId,
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

	return parseResponse(res)
}

func (c *Client) GetBill(id string) (result *BillResponse, err error) {
	res, code, err := c.makeRequest(
		"https://api.qiwi.com/partner/bill/v1/bills/"+id,
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

	return parseResponse(res)
}

func (c *Client) RejectBill(id string) (result *BillResponse, err error) {
	res, code, err := c.makeRequest(
		"https://api.qiwi.com/partner/bill/v1/bills/"+id+"/reject",
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

func (c *Client) makeRequest(url string, method string, data string) (json string, code int, err error) {
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

func (c *Client) StartWebhook(path string, port int) chan PaymentUpdate {
	go c.startListening(path, port)
	return c.ch
}

func (c *Client) startListening(path string, port int) {
	http.HandleFunc(path, c.onWebhook)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

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

	c.ch <- upd
}
