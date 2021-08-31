package qiwiP2P

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Bill
// Object used for creating QIWI P2P bill with parameters
// Amount - money amount
// Comment - text comment (shown to customer)
// ExpirationDateTime - time of expiration of bill, in format 2006-01-02T15:04:05+00:00
// Customer - info about customer, his email, phone and id
// CustomFields - additional fields, can be used for application purposes (for example storing order id)
type Bill struct {
	Amount             Amount            `json:"amount"`
	Comment            string            `json:"comment"`
	ExpirationDateTime string            `json:"expirationDateTime"`
	Customer           Customer          `json:"customer"`
	CustomFields       map[string]string `json:"customFields"`
}

// Amount Object representing some money amount (currency + value)
type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

// Customer
// Object representing customers (this phone, email and id)
type Customer struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Account string `json:"account"`
}

// CreateBill
// Method creates Bill with default parameters
func CreateBill() *Bill {
	return &Bill{
		Amount:             Amount{},
		Customer:           Customer{},
		ExpirationDateTime: time.Now().UTC().Add(time.Hour * 3).Format("2006-01-02T15:04:05+00:00"),
		CustomFields:       make(map[string]string),
	}
}

// SetTheme
// Method sets new theme code in bill (decoration of payment form)
// You can get your theme code here https://qiwi.com/p2p-admin/transfers/link?settings=true
func (b *Bill) SetTheme(theme string) *Bill {
	b.CustomFields["themeCode"] = theme
	return b
}

// SetPaySourcesFilter
// Method sets pay sources filters
// filter variants:
// * "qw,card" - both QIWI Wallet and Credit Card will be available to customer
// * "card" - only Credit Card will be available to customer
// * "qw" - only QIWI Wallet will be available to customer
func (b *Bill) SetPaySourcesFilter(filter string) *Bill {
	b.CustomFields["paySourcesFilter"] = filter
	return b
}

// SetCurrency
// Method sets payment currency
// Currency variants:
// * "RUB" - Russian roubles
// * "KZT" - Kazakhstan tenges
func (b *Bill) SetCurrency(currency string) *Bill {
	b.Amount.Currency = currency
	return b
}

// SetValue
// Method sets value of payment (number of money you need to pay)
func (b *Bill) SetValue(value float32) *Bill {
	b.Amount.Value = fmt.Sprintf("%.2f", value)
	return b
}

// SetComment
// Method sets comment of payment
// Please note that comment will be visible by user
func (b *Bill) SetComment(comment string) *Bill {
	b.Comment = comment
	return b
}

// SetExpirationDateTime
// Method sets expiration moment of payment
// Can be used instead SetExpirationDuration
// time is moment, when payment will be automatically rejected
// Example of usage:
// ```
// b.SetExpirationDateTime(time.Now().Add(time.Day))
// ```
func (b *Bill) SetExpirationDateTime(time time.Time) *Bill {
	b.ExpirationDateTime = time.UTC().Format("2006-01-02T15:04:05+00:00")
	return b
}

// SetExpirationDuration
// Method sets how long payment will be waiting
// Can be used instead of SetExpirationDateTime
// Example of usage:
// ```
// b.SetExpirationDuration(time.Hour * 2)
// ```
func (b *Bill) SetExpirationDuration(duration time.Duration) *Bill {
	b.ExpirationDateTime = time.Now().UTC().Add(duration).Format("2006-01-02T15:04:05+00:00")
	return b
}

// SetCustomerPhone
// Method sets customer phone
// Customer info can be used only for app purposes, no API use of this field
func (b *Bill) SetCustomerPhone(phone string) *Bill {
	b.Customer.Phone = phone
	return b
}

// SetCustomerEmail
// Method sets customer email
// Customer info can be used only for app purposes, no API use of this field
func (b *Bill) SetCustomerEmail(email string) *Bill {
	b.Customer.Email = email
	return b
}

// SetCustomerAccount
// Method sets customer id
// Customer info can be used only for app purposes, no API use of this field
func (b *Bill) SetCustomerAccount(account string) *Bill {
	b.Customer.Account = account
	return b
}

// SetCustomField
// Method sets custom field in bill
// Can be used for storing order ids and other application data
func (b *Bill) SetCustomField(field string, value string) *Bill {
	b.CustomFields[field] = value
	return b
}

// toJSON
// Method converts bill object to JSON, for API transferring
func (b *Bill) toJSON() string {
	arr, err := json.Marshal(b)
	if err != nil {
		log.Println("Error on toJSON Bill: " + err.Error())
	}
	return string(arr)
}
