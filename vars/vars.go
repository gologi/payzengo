package vars

const (
	VadsSiteID    string = "vads_site_id"
	VadsCtxMode   string = "vads_ctx_mode"
	VadsTransID   string = "vads_trans_id"
	VadsTransDate string = "vads_trans_date"
	VadsOrderID   string = "vads_order_id"
	VadsCustID    string = "vads_cust_id"
	//VadsCustEmail adresse e-mail de l’acheteur, nécessaire si vous souhaitez que la plateforme de paiement envoie un e-mail à l’acheteur. Pour que l'acheteur reçoive un e-mail, n'oubliez pas de poster ce paramètre dans le formulaire lorsque vous générez une demande de paiement
	VadsCustEmail              string = "vads_cust_email"
	VadsAmount                 string = "vads_amount"
	VadsCurrency               string = "vads_currency"
	VadsActionMode             string = "vads_action_mode"
	VadsPageAction             string = "vads_page_action"
	VadsVersion                string = "vads_version"
	VadsPaymentConfig          string = "vads_payment_config"
	VadsCaptureDelay           string = "vads_capture_delay"
	VadsValidationMode         string = "vads_validation_mode"
	VadsSignature              string = "signature"
	VadsURLSuccess             string = "vads_url_success"
	VadsURLRefused             string = "vads_url_refused"
	VadsURLCancel              string = "vads_url_cancel"
	VadsURLError               string = "vads_url_error"
	VadsURLReturn              string = "vads_url_return"
	VadsLanguage               string = "vads_language"
	VadsAvailableLanguages     string = "vads_available_languages"
	VadsReturnMode             string = "vads_return_mode"
	VadsRedirectSuccessTimeout string = "vads_redirect_success_timeout"
	VadsRedirectSuccessMessage string = "vads_redirect_success_message"
	VadsRedirectErrorTimeout   string = "vads_redirect_error_timeout"
	VadsRedirectErrorMessage   string = "vads_redirect_error_message"
)