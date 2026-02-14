package dto

type CreateRechargeRequest struct {
	Amount float64 `json:"amount"`
}

type CreateRechargeResponse struct {
	CheckoutURL string `json:"checkoutUrl"`
	SessionID   string `json:"sessionId"`
}
