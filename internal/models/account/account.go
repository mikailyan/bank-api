package account

import (
	"github.com/shopspring/decimal"
	"time"
)

type Account struct {
	ID        int64           `db:"id"       json:"id"`
	UserID    int64           `db:"user_id"  json:"user_id"`
	Balance   decimal.Decimal `db:"balance"  json:"balance"`
	Currency  Currency        `db:"currency" json:"currency"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
}
