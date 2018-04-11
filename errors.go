package payzengo

import "fmt"

//PaymentError retour d'erreur sur le module de paiement
type PaymentError int

func (pe PaymentError) Error() string {
	return fmt.Sprintf("Erreur Payment %#v\n", pe)
}

const (
	PaymentErrorBadTransactionID PaymentError = iota
)
