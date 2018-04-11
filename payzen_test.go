package payzengo

import (
	"testing"
	"time"
)

func TestSignature(t *testing.T) {

	testSite := &PayzenSite{
		Site:            "TEST",
		SiteId:          12345678,
		CertificateDev:  "1122334455667788",
		CertificateProd: "1122334455667789",
	}

	config := &PaymentConfig{
		SiProduction: false,
		Site:         testSite,
	}

	pz, err := GetNewPaiement(config,
		time.Date(2017, 1, 29, 13, 0, 25, 0, time.UTC),
		5124,
		123456,
		0,
		0)

	if err != nil {
		t.Error(err)
	}
	signature := pz.GetSignature()
	if signature != "59c96b34c74b9375c332b0b6a32e6deeec87de2b" {
		t.Error("mauvaise génération de signature")
	}
}

func TestTransactionID(t *testing.T) {

	testSite := &PayzenSite{
		Site:            "TEST",
		SiteId:          12345678,
		CertificateDev:  "1122334455667788",
		CertificateProd: "1122334455667789",
	}

	config := &PaymentConfig{
		SiProduction: false,
		Site:         testSite,
	}

	_, err := GetNewPaiement(config,
		time.Date(2017, 1, 29, 13, 0, 25, 0, time.UTC),
		5124,
		999999,
		0,
		0)
	if err != PaymentErrorBadTransactionID {
		t.Error("Mauvaise transaction non détectée")
	}
}
