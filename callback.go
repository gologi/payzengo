package payzengo

import (
	"log"
	"net/http"
	"sort"
	"strconv"
)

type CallbackInfo struct {
	Site       *PayzenSite
	Production bool
	Valid      bool
}

func CallbackCheck(r *http.Request) *CallbackInfo {
	var args argsArr
	r.ParseForm()
	log.Println(r.Method)
	if r.Method == "GET" {
		args = getFormData(r)
	} else {
		args = getPostData(r)
	}
	log.Printf("%v\n", args)
	siteid, _ := strconv.ParseUint(r.FormValue("vads_site_id"), 0, 64)
	SiProduction := r.FormValue("vads_ctx_mode") != "TEST"
	certificate, site := GetCertificate(siteid, SiProduction)
	signature := args.GenSignature(certificate)
	signatureForm := r.FormValue("signature")
	log.Println(signature, signatureForm)
	if signatureForm == "" || signature == "" || certificate == "" || site == nil {
		return &CallbackInfo{
			Valid:      false,
			Production: SiProduction,
			Site:       &PayzenSite{},
		}
	}
	return &CallbackInfo{
		Site:       site,
		Production: SiProduction,
		Valid:      signature == signatureForm,
	}
}

func getFormData(r *http.Request) argsArr {
	var x argsArr
	for a, b := range r.Form {
		if len(a) > 5 && a[:5] == "vads_" {
			x = append(x, []string{a, b[0]})
		}
	}
	sort.Sort(&x)
	return x
}

func getPostData(r *http.Request) argsArr {
	var x argsArr
	for a, b := range r.PostForm {
		if len(a) > 5 && a[:5] == "vads_" {
			x = append(x, []string{a, b[0]})
		}
	}
	sort.Sort(&x)
	return x
}
