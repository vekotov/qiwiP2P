package qiwiP2P

import (
	"crypto/rand"
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

func (c *Client) PutBill(b *Bill) (json string, retErr error) {
	billId := pseudoUUID()
	req, err := http.NewRequest(
		"PUT", "https://api.qiwi.com/partner/bill/v1/bills/"+billId, strings.NewReader(b.toJSON()),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("User-Agent", "GoQiwiP2P")
	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	buf := new(strings.Builder)
	io.Copy(buf, res.Body)

	return buf.String(), nil
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
