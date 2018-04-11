package payzengo

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gologi/payzengo/vars"
)

type PayzenSite struct {
	Site            string
	SiteId          uint64
	CertificateDev  string
	CertificateProd string
}

func (site *PayzenSite) GetCertificate(SiProduction bool) string {
	if SiProduction {
		return site.CertificateProd
	}
	return site.CertificateDev
}

type PayzenPaiement interface {
	GetSignature() string
	GetForm(string) string
}

type PaymentConfig struct {
	Site                   *PayzenSite
	Currency               vars.Currency
	PaymentType            vars.PaymentType
	ActionMode             vars.ActionMode
	SiProduction           bool
	URLSuccess             string
	URLRefused             string
	URLCancel              string
	URLError               string
	URLReturn              string
	Language               vars.Language
	AvailableLanguages     []vars.Language
	ReturnMode             vars.ReturnMode
	PageAction             vars.PageAction
	Version                string
	RedirectErrorMessage   string
	RedirectErrorTimeout   string
	RedirectSuccessTimeout string
	RedirectSuccessMessage string
}

func (c *PaymentConfig) init() {
	if c.PaymentType == "" {
		c.PaymentType = vars.PaymentTypeSingle
	}
	if c.Currency == 0 {
		c.Currency = vars.CurrencyEuro
	}
	if c.ActionMode == "" {
		c.ActionMode = vars.ActionModeInteractive
	}
	if c.PageAction == "" {
		c.PageAction = vars.PageActionPayment
	}
	if c.Language == "" {
		c.Language = "fr"
	}
	if c.Version == "" {
		c.Version = "V2"
	}
}

type payzenPaiement struct {
	Config *PaymentConfig
	Amount uint64
	//TransactionID Il est constitué de 6 caractères numériques et doit être unique pour chaque transaction
	// pour une boutique donnée sur la journée. Remarque : l’unicité de l’identifiant de transaction se base
	//sur l’heure universelle (UTC). Il est à la charge du site marchand de garantir cette unicité
	//sur la journée. Il doit être compris entre 000000 et 899999.
	// La tranche 900000 et 999999 est reservée aux remboursements et aux opérations effectuées depuis le Back Office.
	TransactionID   int
	TransactionDate string
	CaptureDelay    string
	CustomerID      int
	OrderID         int
	ClientID        int
	signature       string
}

func getBoolValue(x bool, truevalue, falsevalue string) string {
	if x {
		return truevalue
	}
	return falsevalue
}

func (p *payzenPaiement) GetSignature() string {
	if p.signature == "" {
		p.genSignature()
	}
	return p.signature
}

func (p *PaymentConfig) SetAutomaticReturn(msg string, timeout int, msgerror string, timeouterror int) {
	p.ReturnMode = vars.ReturnModeGet
	if msg == "" {
		msg = "Vous allez être rédirigé vers le site marchand"
	}
	if msgerror == "" {
		msgerror = msg
	}
	p.RedirectSuccessMessage = msg
	p.RedirectSuccessTimeout = strconv.Itoa(timeout)
	p.RedirectErrorMessage = msgerror
	p.RedirectErrorTimeout = strconv.Itoa(timeouterror)
}

func (p *payzenPaiement) genSignature() string {
	var sar []string

	for _, b := range p.getMap() {
		sar = append(sar, b[1])
	}
	str := strings.Join(sar, "+") + "+" + p.Config.Site.GetCertificate(p.Config.SiProduction)
	p.signature = getSha1(str)
	return p.signature
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

func (p *payzenPaiement) getMap() argsArr {
	m := argsArr{
		{vars.VadsActionMode, string(p.Config.ActionMode)},
		{vars.VadsAmount, strconv.FormatUint(p.Amount, 10)},
		{vars.VadsPageAction, string(p.Config.PageAction)},
		{vars.VadsCurrency, strconv.Itoa(int(p.Config.Currency))},
		{vars.VadsCtxMode, getBoolValue(p.Config.SiProduction, vars.PayzenCtxProduction, vars.PayzenCtxTest)},
		{vars.VadsPaymentConfig, string(p.Config.PaymentType)},
		{vars.VadsSiteID, strconv.FormatUint(p.Config.Site.SiteId, 10)},
		{vars.VadsTransDate, p.TransactionDate},
		{vars.VadsTransID, fmt.Sprintf("%06d", p.TransactionID)},
		{vars.VadsVersion, p.Config.Version},
		//{vars.VadsSignature, p.signature},
	}
	if p.OrderID > 0 {
		m = append(m, []string{vars.VadsOrderID, strconv.Itoa(p.OrderID)})
	}
	if p.Config.ReturnMode != "" {
		m = append(m,
			[]string{vars.VadsReturnMode, string(p.Config.ReturnMode)},
		)
	}
	if p.CaptureDelay != "" {
		m = append(m,
			[]string{vars.VadsCaptureDelay, p.CaptureDelay},
			[]string{vars.VadsValidationMode, "1"}, //strconv.Itoa(p.ValidationMode)},
		)
	}

	if p.Config.Language != "" {
		m = append(m,
			[]string{vars.VadsLanguage, string(p.Config.Language)})
	}

	if len(p.Config.AvailableLanguages) > 0 {
		var langs []string
		for lang := range p.Config.AvailableLanguages {
			langs = append(langs, string(lang))
		}
		m = append(m,
			[]string{vars.VadsAvailableLanguages, strings.Join(langs, ",")})

	}

	if p.Config.URLReturn != "" {
		m = append(m,
			[]string{vars.VadsURLReturn, p.Config.URLReturn},
		)
	}

	if p.Config.URLSuccess != "" {
		m = append(m,
			[]string{vars.VadsURLSuccess, p.Config.URLSuccess},
		)
	}

	if p.Config.URLRefused != "" {
		m = append(m,
			[]string{vars.VadsURLRefused, p.Config.URLRefused},
		)
	}

	if p.Config.URLCancel != "" {
		m = append(m,
			[]string{vars.VadsURLCancel, p.Config.URLCancel},
		)
	}

	if p.Config.URLError != "" {
		m = append(m,
			[]string{vars.VadsURLError, p.Config.URLError},
		)
	}
	// TODO
	/*
		if p.RedirectSuccessTimeout != "" {
			m = append(m,
				[]string{vars.VadsRedirectSuccessMessage, p.RedirectSuccessMessage},
				[]string{vars.VadsRedirectSuccessTimeout, p.RedirectSuccessTimeout},
				[]string{vars.VadsRedirectSuccessMessage, p.RedirectSuccessMessage},
				[]string{vars.VadsRedirectSuccessTimeout, p.RedirectSuccessTimeout},
			)
		}
	*/
	sort.Sort(&m)
	return m
}

func (p *payzenPaiement) GetForm(htmlSubmitButton string) string {
	return strings.Join([]string{
		`<form method="POST" action="https://secure.payzen.eu/vads-payment/">`,
		p.getFormInput(),
		htmlSubmitButton,
		"</form>"}, "")
}
func (p *payzenPaiement) getFormInput() string {
	sarray := []string{}

	for _, b := range append(p.getMap(), []string{vars.VadsSignature, p.GetSignature()}) {
		sarray = append(sarray, fmt.Sprintf(`<input type=hidden name=%s value="%s"/>`, b[0], b[1]))
	}
	return strings.Join(sarray, "")
}

func GetNewPaiement(payzenConfig *PaymentConfig, dt time.Time, amountCents uint64, transactionID, orderID, clientID int) (PayzenPaiement, error) {
	if transactionID < 1 || transactionID > 899999 {
		return nil, PaymentErrorBadTransactionID
	}
	p := &payzenPaiement{
		Config:          payzenConfig,
		TransactionID:   transactionID,
		OrderID:         orderID,
		ClientID:        clientID,
		Amount:          amountCents,
		TransactionDate: dt.UTC().Format("20060102150405"),
	}
	p.Config.init()
	return p, nil
}
