# payzengo
Go support for the PayZen online payment solution provided by lyra


```

site:=payzengo.PayzenSite{
	Site:            "TEST",
	SiteID:          11111111,
  CertificateDev: "22222",
  CertificateProd: "22222",
  }
confPayzen := (&payzengo.PaymentConfig{
	SiProduction: true,
	Site:         site,
	ReturnMode:   "POST",
})


amount:=100
transactionid:=4000 // unique transaction id
orderid:=1000
clientid:=1
  
p, err := payzengo.GetNewPaiement(confPayzen, time.Now(), amount, transactionid, orderid, clientid)
if err != nil {
	log.Println(err)
}
p.SetData("custom_form_var", "hop")
paymentForm:= p.GetForm("<input type=submit/><script>payzenform.submit();</script>"))
```
