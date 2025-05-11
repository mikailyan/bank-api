package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type PaymentSchedule struct {
	ID        int64           `db:"id"         json:"id"`
	CreditID  int64           `db:"credit_id"  json:"credit_id"`
	DueDate   time.Time       `db:"due_date"   json:"due_date"`
	Amount    decimal.Decimal `db:"amount" json:"amount"`
	Paid      bool            `db:"paid"       json:"paid"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
}
