# Golang QIWI P2P library
Usage:
```go
// create API client
c := qiwiP2P.CreateClient(token)
// create bill object
b := qiwiP2P.CreateBill()

// setup bill parameters
b.SetValue(10.57)
b.SetCurrency("RUB")
b.SetExpirationDateTime(time.Now().Add(time.Minute * 15))

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

```