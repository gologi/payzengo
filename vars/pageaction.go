package vars

//PageAction Définit l'action à réaliser
type PageAction string

const (
	//PageActionPayment Paiement (avec ou sans alias)
	PageActionPayment = "PAYMENT"
	//PageActionRegister Inscription sans paiement
	PageActionRegister = "REGISTER"
	//PageActionRegisterUpdate Mise à jour des informations du moyen de paiement
	PageActionRegisterUpdate = "REGISTER_UPDATE"

	//PageActionRegisterPay Inscription avec paiement
	PageActionRegisterPay = "REGISTER_PAY"

	//PageActionRegisterSubscribe Inscription avec souscription à un abonnement
	PageActionRegisterSubscribe = "REGISTER_SUBSCRIBE"

	//PageActionRegisterPaySuscribe Inscription avec paiement et souscription à un abonnement
	PageActionRegisterPaySuscribe = "REGISTER_PAY_SUBSCRIBE"

	//PageActionSubscribe Souscription à un abonnement
	PageActionSubscribe = "SUBSCRIBE"

	//PageActionRegisterUpdatePay Mise à jour des informations du moyen de paiement avec paiement
	PageActionRegisterUpdatePay = "REGISTER_UPDATE_PAY"
	//PageactionAskRegisterPay Paiement avec inscription optionnelle du porteur
	PageactionAskRegisterPay = "ASK_REGISTER_PAY"
)
