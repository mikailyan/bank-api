package service

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
	"github.com/therealadik/bank-api/internal/models/account"
	"github.com/therealadik/bank-api/internal/models/transaction"
	"github.com/therealadik/bank-api/internal/repository"
)

var (
	ErrInsufficientFunds = errors.New("недостаточно средств")
	ErrSameAccount       = errors.New("нельзя переводить деньги на тот же счет")
	ErrNegativeAmount    = errors.New("сумма не может быть отрицательной")
)

type AccountService struct {
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
}

func NewAccountService(accountRepo *repository.AccountRepository, transactionRepo *repository.TransactionRepository) *AccountService {
	return &AccountService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID int64, currency account.Currency) (*account.Account, error) {
	return s.accountRepo.CreateAccount(ctx, userID, currency)
}

func (s *AccountService) GetAccountByID(ctx context.Context, id int64, userID int64) (*account.Account, error) {
	acc, err := s.accountRepo.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if acc.UserID != userID {
		return nil, errors.New("счет не принадлежит пользователю")
	}

	return acc, nil
}

func (s *AccountService) GetAccountsByUserID(ctx context.Context, userID int64) ([]*account.Account, error) {
	return s.accountRepo.GetAccountsByUserID(ctx, userID)
}

func (s *AccountService) UpdateBalance(ctx context.Context, id int64, userID int64, amount decimal.Decimal) error {
	if amount.Equal(decimal.Zero) {
		return errors.New("сумма должна быть отлична от нуля")
	}

	acc, err := s.GetAccountByID(ctx, id, userID)
	if err != nil {
		return err
	}

	if amount.LessThan(decimal.Zero) && acc.Balance.Add(amount).LessThan(decimal.Zero) {
		return ErrInsufficientFunds
	}

	txType := transaction.WITHDRAWAL
	if amount.GreaterThan(decimal.Zero) {
		txType = transaction.DEPOSIT
	}

	err = s.accountRepo.UpdateBalance(ctx, id, amount)
	if err != nil {
		return err
	}

	absAmount := amount.Abs()
	_, err = s.transactionRepo.CreateTransaction(ctx, id, absAmount, txType, transaction.COMPLETED)

	return err
}

func (s *AccountService) Transfer(ctx context.Context, fromID, toID int64, userID int64, amount decimal.Decimal) error {
	if fromID == toID {
		return ErrSameAccount
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrNegativeAmount
	}

	fromAcc, err := s.GetAccountByID(ctx, fromID, userID)
	if err != nil {
		return err
	}

	if fromAcc.Balance.LessThan(amount) {
		return ErrInsufficientFunds
	}

	_, err = s.accountRepo.GetAccountByID(ctx, toID)
	if err != nil {
		return err
	}

	err = s.accountRepo.TransferBetweenAccounts(ctx, fromID, toID, amount)
	if err != nil {
		return err
	}

	_, err = s.transactionRepo.CreateTransaction(ctx, fromID, amount, transaction.DEPOSIT, transaction.COMPLETED)
	if err != nil {
		return err
	}

	_, err = s.transactionRepo.CreateTransaction(ctx, toID, amount, transaction.WITHDRAWAL, transaction.COMPLETED)
	return err
}

func (s *AccountService) GetTransactionsByAccountID(ctx context.Context, accountID int64, userID int64) ([]*transaction.Transaction, error) {
	_, err := s.GetAccountByID(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}

	return s.transactionRepo.GetTransactionsByAccountID(ctx, accountID)
}

func (s *AccountService) GetTransactionsByUserID(ctx context.Context, userID int64) ([]*transaction.Transaction, error) {
	return s.transactionRepo.GetTransactionsByUserID(ctx, userID)
}
