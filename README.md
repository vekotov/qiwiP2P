# Golang QIWI P2P library
Usage:
```go
import "github.com/vekotov/qiwiP2P"

// create API client
c := qiwiP2P.CreateClient(token)
// create bill object
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

// get bill info
result, err := c.GetBill(result.BillId)
if err != nil {
    println(err.Description)
    return
}

// print bill status
println(result.Status.Value)

// reject bill
result, err := c.RejectBill(result.BillId)
if err != nil {
    println(err.Description)
    return
}

// start webhook (you need to setup address on qiwi p2p page first)
ch := c.StartWebhook("/qiwi", 80)
for upd := range ch {
    // print bill id when bill is paid
    println(upd.Bill.BillId)
}

```
