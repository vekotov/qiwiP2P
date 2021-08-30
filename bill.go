package qiwiP2P

import (
	"fmt"
	"time"
)

type Bill struct {
	Currency string
	Value float32
	Comment string
	ExpirationDateTime time.Time
	CustomerPhone string
	CustomerEmail string
	CustomerAccount string
	CustomFields map[string]string
}

func CreateBill() *Bill {
	return &Bill{
		Currency: "RUB",
		Value: 1,
		ExpirationDateTime: time.Now().UTC().Add(time.Hour * 3),
	}
}

func (b *Bill) SetCurrency(currency string) *Bill {
	b.Currency = currency
	return b
}

func (b *Bill) SetValue(value float32) *Bill {
	b.Value = value
	return b
}

func (b *Bill) toJSON() string {
	return fmt.Sprintf(
		"{\"amount\":{\"currency\":\"%s\",\"value\": \"%.2f\"}," +
			"\"expirationDateTime\": \"%s\"}",
			b.Currency,
			b.Value,
			b.ExpirationDateTime.Format("2006-01-02T15:01:05+00:00"))
}
