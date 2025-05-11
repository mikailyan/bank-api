package dto

type CreateCardRequest struct {
	PGPKey string `json:"pgp_key"`
}

type CreateCardResponse struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	CreatedAt  string `json:"created_at"`
	CardNumber string `json:"card_number"`
	Expire     string `json:"expire"`
	CVV        string `json:"cvv"`
}

type CardResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type CardDetailsResponse struct {
	ID         int64  `json:"id"`
	CardNumber string `json:"card_number"`
	Expire     string `json:"expire"`
}

type CardListResponse struct {
	Cards []CardResponse `json:"cards"`
}

type CardPaymentRequest struct {
	CardID int64  `json:"card_id"`
	Amount string `json:"amount"`
	CVV    string `json:"cvv"`
	PGPKey string `json:"pgp_key"`
}

type CardPaymentResponse struct {
	Success     bool   `json:"success"`
	PaymentID   string `json:"payment_id,omitempty"`
	Description string `json:"description,omitempty"`
}
