package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/therealadik/bank-api/internal/models/transaction"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, accountID int64, amount decimal.Decimal,
	txType transaction.Type, status transaction.Status) (*transaction.Transaction, error) {
	query := `
		INSERT INTO transactions (account_id, amount, type, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, account_id, amount, type, status, created_at
	`
	var tx transaction.Transaction
	err := r.db.QueryRow(ctx, query, accountID, amount, txType, status).Scan(
		&tx.ID, &tx.AccountID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionRepository) GetTransactionsByAccountID(ctx context.Context, accountID int64) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, account_id, amount, type, status, created_at
		FROM transactions
		WHERE account_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*transaction.Transaction
	for rows.Next() {
		var tx transaction.Transaction
		if err := rows.Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetTransactionsByUserID(ctx context.Context, userID int64) ([]*transaction.Transaction, error) {
	query := `
		SELECT t.id, t.account_id, t.amount, t.type, t.status, t.created_at
		FROM transactions t
		JOIN accounts a ON t.account_id = a.id
		WHERE a.user_id = $1
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*transaction.Transaction
	for rows.Next() {
		var tx transaction.Transaction
		if err := rows.Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}
