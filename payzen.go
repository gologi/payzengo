package payzengo

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gologi/payzengo/vars"
)

type PayzenSite struct {
	Site            string
	SiteID          uint64
	CertificateDev  string
	CertificateProd string
}

var payzenSites map[uint64]*PayzenSite

func init() {
	payzenSites = make(map[uint64]*PayzenSite)
}

func GetSite(siteid uint64) (*PayzenSite, bool) {
	s, ok := payzenSites[siteid]
	return s, ok
}
func GetCertificate(siteid uint64, SiProduction bool) (string, *PayzenSite) {
	site, ok := GetSite(siteid)
	if !ok {
		return "", nil
	}
	return site.GetCertificate(SiProduction), site
}

func (site *PayzenSite) Register() *PayzenSite {
	payzenSites[site.SiteID] = site
	return site
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
	SetData(string, string)
}

type PaymentConfig struct {
	Site                   *PayzenSite
	SiHMAC                 bool
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
	Data            map[string]string
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

func (p *PaymentConfig) SetAutomaticReturn(msg string, timeout int, msgerror string, timeouterror int) *PaymentConfig {
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
	return p
}

func (p *payzenPaiement) getMap() argsArr {
	m := argsArr{
		{vars.VadsActionMode, string(p.Config.ActionMode)},
		{vars.VadsAmount, strconv.FormatUint(p.Amount, 10)},
		{vars.VadsPageAction, string(p.Config.PageAction)},
		{vars.VadsCurrency, strconv.Itoa(int(p.Config.Currency))},
		{vars.VadsCtxMode, getBoolValue(p.Config.SiProduction, vars.PayzenCtxProduction, vars.PayzenCtxTest)},
		{vars.VadsPaymentConfig, string(p.Config.PaymentType)},
		{vars.VadsSiteID, strconv.FormatUint(p.Config.Site.SiteID, 10)},
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
	if p.Config.RedirectSuccessTimeout != "" {
		m = append(m,
			[]string{vars.VadsRedirectSuccessMessage, p.Config.RedirectSuccessMessage},
			[]string{vars.VadsRedirectSuccessTimeout, p.Config.RedirectSuccessTimeout},
			[]string{vars.VadsRedirectErrorMessage, p.Config.RedirectErrorMessage},
			[]string{vars.VadsRedirectErrorTimeout, p.Config.RedirectErrorTimeout},
		)
	}
	sort.Sort(&m)
	return m
}
func (p *payzenPaiement) SetData(key, value string) {
	key = strings.ToLower(key)
	p.Data[key] = value
}
func (p *payzenPaiement) GetForm(htmlSubmitButton string) string {
	return strings.Join([]string{
		`<form method="POST" action="https://secure.payzen.eu/vads-payment/" id="payzenform">`,
		p.getFormInput(),
		htmlSubmitButton,
		"</form>"}, "")
}
func (p *payzenPaiement) getFormInput() string {
	sarray := []string{}

	for _, b := range append(p.getMap(), []string{vars.VadsSignature, p.GetSignature()}) {
		sarray = append(sarray, fmt.Sprintf(`<input type=hidden name=%s value="%s"/>`, b[0], b[1]))
	}
	for k, v := range p.Data {
		sarray = append(sarray, fmt.Sprintf(`<input type=hidden name=%s value="%s"/>`, k, v))
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
		Data:            make(map[string]string),
		TransactionDate: dt.UTC().Format("20060102150405"),
	}
	p.Config.init()
	return p, nil
}
