package transaction

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID        int64           `db:"id"          json:"id"`
	AccountID int64           `db:"account_id"  json:"account_id"`
	Amount    decimal.Decimal `db:"amount" json:"amount"`
	Type      Type            `db:"type"        json:"type"`
	Status    Status          `db:"status"      json:"status"`
	CreatedAt time.Time       `db:"created_at"  json:"created_at"`
}
