package models

import "time"

type Card struct {
	ID         int64     `db:"id"        json:"id"`
	UserID     int64     `db:"user_id"   json:"user_id"`
	CardNumber []byte    `db:"card_number" json:"-"`
	Expire     []byte    `db:"expire"      json:"-"`
	CVVHash    string    `db:"cvv_hash"    json:"-"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
