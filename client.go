package qiwiP2P

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	token  string
	client http.Client
}

func CreateClient(key string) *Client {
	return &Client{token: key, client: http.Client{}}
}

func (c *Client) SetSecretKey(key string) *Client {
	c.token = key
	return c
}

func (c *Client) PutBill(b *Bill) (result *BillResponse, error *RequestError) {
	billId := pseudoUUID()
	res, code := c.makeRequest("https://api.qiwi.com/partner/bill/v1/bills/"+billId, "PUT", b.toJSON())

	if code == 401 {
		return nil, &RequestError{ErrorCode: "bad_token", Description: "Bad token"}
	}

	return parseResponse(res)
}

func (c *Client) GetBill(id string) (result *BillResponse, error *RequestError) {
	res, code := c.makeRequest("https://api.qiwi.com/partner/bill/v1/bills/"+id, "GET", "")

	if code == 401 {
		return nil, &RequestError{ErrorCode: "bad_token", Description: "Bad token"}
	}
	if code == 404 {
		return nil, &RequestError{ErrorCode: "bad_id", Description: "No such bill found"}
	}

	return parseResponse(res)
}

func (c *Client) RejectBill(id string) (result *BillResponse, error *RequestError) {
	res, code := c.makeRequest("https://api.qiwi.com/partner/bill/v1/bills/"+id+"/reject", "POST", "")

	if code == 401 {
		return nil, &RequestError{ErrorCode: "bad_token", Description: "Bad token"}
	}
	if code == 404 {
		return nil, &RequestError{ErrorCode: "bad_id", Description: "No such bill found"}
	}

	return parseResponse(res)
}

func parseResponse(jsonResponse string) (result *BillResponse, error *RequestError) {
	var re RequestError
	json.Unmarshal([]byte(jsonResponse), &error)
	if re.ErrorCode != "" {
		return nil, &re
	}

	var response BillResponse
	json.Unmarshal([]byte(jsonResponse), &response)

	return &response, nil
}

func (c *Client) makeRequest(url string, method string, data string) (json string, code int) {
	req, _ := http.NewRequest(
		method, url, strings.NewReader(data),
	)

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("content-type", "application/json")
	res, _ := c.client.Do(req)

	defer res.Body.Close()
	buf := new(strings.Builder)
	io.Copy(buf, res.Body)
	return buf.String(), res.StatusCode
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
