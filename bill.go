package qiwiP2P

import (
	"encoding/json"
	"fmt"
	"time"
)

type Bill struct {
	Amount             Amount            `json:"amount"`
	Comment            string            `json:"comment"`
	ExpirationDateTime string            `json:"expirationDateTime"`
	Customer           Customer          `json:"customer"`
	CustomFields       map[string]string `json:"customFields"`
}

type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

type Customer struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Account string `json:"account"`
}

func CreateBill() *Bill {
	return &Bill{
		Amount:             Amount{},
		Customer:           Customer{},
		ExpirationDateTime: time.Now().UTC().Add(time.Hour * 3).Format("2006-01-02T15:01:05+00:00"),
		CustomFields:       make(map[string]string),
	}
}

func (b *Bill) SetTheme(theme string) *Bill {
	b.CustomFields["themeCode"] = theme
	return b
}

func (b *Bill) SetPaySourcesFilter(filter string) *Bill {
	b.CustomFields["paySourcesFilter"] = filter
	return b
}

func (b *Bill) SetCurrency(currency string) *Bill {
	b.Amount.Currency = currency
	return b
}

func (b *Bill) SetValue(value float32) *Bill {
	b.Amount.Value = fmt.Sprintf("%.2f", value)
	return b
}

func (b *Bill) SetComment(comment string) *Bill {
	b.Comment = comment
	return b
}

func (b *Bill) SetExpirationDateTime(time time.Time) *Bill {
	b.ExpirationDateTime = time.UTC().Format("2006-01-02T15:04:05+00:00")
	return b
}

func (b *Bill) SetCustomerPhone(phone string) *Bill {
	b.Customer.Phone = phone
	return b
}

func (b *Bill) SetCustomerEmail(email string) *Bill {
	b.Customer.Email = email
	return b
}

func (b *Bill) SetCustomerAccount(account string) *Bill {
	b.Customer.Account = account
	return b
}

func (b *Bill) SetCustomField(field string, value string) *Bill {
	b.CustomFields[field] = value
	return b
}

func (b *Bill) toJSON() string {
	arr, _ := json.Marshal(b)
	return string(arr)
}
