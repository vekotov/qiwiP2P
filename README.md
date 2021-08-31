# Golang QIWI P2P library
This library can be used to access QIWI P2P Payments API (https://p2p.qiwi.com)

To get library run
```shell
go get -u github.com/vekotov/qiwiP2P
```
## Usage
### Importing library
```go
import "github.com/vekotov/qiwiP2P"
```
### Creating API Client
```go
c := qiwiP2P.CreateClient(token)
```
### Creating new payment bill
```go
b := qiwiP2P.CreateBill()

// setup bill parameters
b.SetValue(10.57)
b.SetCurrency("RUB")
b.SetExpirationDuration(time.Minute * 5)

// put bill
result, err := c.PutBill(b)
if err != nil {
	println(err.Description)
	return
}
// print bill url
println(result.PayUrl)
```
### Getting bill info
```go
result, err := c.GetBill(billId)
if err != nil {
    println(err.Description)
    return
}

// print bill status
println(result.Status.Value)
```
### Rejecting bill
```go
result, err := c.RejectBill(result.BillId)
if err != nil {
    println(err.Description)
    return
}
```
### Using webhook
For using webhook you need go to https://qiwi.com/p2p-admin/transfers/api 
and create new private key with webhook turned on.

Also don't forget that you need to add your backend URL in key settings.
```go
// start webhook (you need to setup address on qiwi p2p page first)
ch := c.StartWebhook("/qiwi", 80)
for upd := range ch {
    // print bill id when bill is paid
    println(upd.Bill.BillId)
}
```

## Authors
vekotov (tg: @vekotov), 2021

Library licensed under MIT License