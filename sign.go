package payzengo

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"strings"
)

func (p *payzenPaiement) genSignature() string {
	if p.Config.SiHMAC {
		p.signature = p.getMap().GenSignatureHMAC(p.Config.Site.GetCertificate(p.Config.SiProduction))
	} else {
		p.signature = p.getMap().GenSignature(p.Config.Site.GetCertificate(p.Config.SiProduction))
	}
	return p.signature
}

func (arr argsArr) GenSignatureHMAC(key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	var sar []string
	for _, b := range arr {
		sar = append(sar, b[1])
	}
	str := strings.Join(sar, "+")
	mac.Write([]byte(str))
	return string(mac.Sum(nil))
}

func (arr argsArr) GenSignature(certificate string) string {
	var sar []string
	for _, b := range arr {
		sar = append(sar, b[1])
	}
	str := strings.Join(sar, "+") + "+" + certificate
	return getSha1(str)

}

func getSha1(s string) string {
	sha := sha1.New()
	sha.Write([]byte(s))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

type argsArr [][]string

func (a argsArr) Len() int { return len(a) }
func (a argsArr) Less(i, j int) bool {
	return strings.ToLower(a[i][0]) < strings.ToLower(a[j][0])
}
func (a argsArr) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
