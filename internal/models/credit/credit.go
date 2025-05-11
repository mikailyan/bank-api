package credit

import (
	"github.com/shopspring/decimal"
	"time"
)

type Credit struct {
	ID           int64           `db:"id"            json:"id"`
	AccountID    int64           `db:"account_id"    json:"account_id"`
	Principal    decimal.Decimal `db:"principal" json:"principal"`
	InterestRate float64         `db:"interest_rate" json:"interest_rate"`
	TermMonths   int             `db:"term_months"   json:"term_months"`
	StartDate    time.Time       `db:"start_date"    json:"start_date"`
	Status       Status          `db:"status"        json:"status"`
	CreatedAt    time.Time       `db:"created_at"    json:"created_at"`
}
