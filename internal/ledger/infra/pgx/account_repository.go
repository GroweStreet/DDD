package pgx

import (
	"FinFlow/internal/ledger/domain"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type AccountRepository struct {
	db dbtx
}

func NewAccountRepository(db dbtx) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Save(ctx context.Context, acc *domain.Account) error {

	var err error

	if acc.Version() == 0 {
		_, err = r.db.Exec(ctx, `INSERT INTO accounts (id, user_id, currency, balance_amount, status, account_type, version, created_at)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8)`, acc.ID(), nullString(string(acc.UserID())),
			string(acc.Currency().Code()), acc.Balance().Amount(), string(acc.Status()), string(acc.Type()), acc.Version(), acc.Created())
		return err
	}

	tag, err := r.db.Exec(ctx,
		`UPDATE accounts SET balance_amount=$1, status=$2, version=$3                                                  
                 WHERE id=$4 AND version=$5`,
		acc.Balance().Amount(), acc.Status(), acc.Version(),
		string(acc.ID()), acc.Version()-1)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrVersionConflict
	}
	return nil
}

func (r *AccountRepository) FindByID(ctx context.Context, id domain.AccountID) (*domain.Account, error) {

	var (
		accID     string
		userID    *string
		currency  string
		balance   int64
		status    string
		accType   string
		version   int64
		createdAt time.Time
	)

	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, currency, balance_amount, status, account_type, version, created_at                       
                 FROM accounts WHERE id = $1`, string(id)).
		Scan(&accID, &userID, &currency, &balance, &status, &accType, &version, &createdAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	c, err := domain.NewCurrency(currency)
	if err != nil {
		return nil, err
	}
	money := domain.NewMoney(balance, c)

	uid := domain.UserID("")
	if userID != nil {
		uid = domain.UserID(*userID)
	}
	account := domain.ReconstituteAccount(domain.AccountID(accID), uid, c, money, domain.AccountStatus(status), version, domain.AccountType(accType), createdAt)

	return account, nil
}

func (r *AccountRepository) FindSystemAccount(ctx context.Context, currency domain.Currency) (*domain.Account, error) {

	var (
		accID        string
		userID       *string
		currencyCode string
		balance      int64
		status       string
		accType      string
		version      int64
		createdAt    time.Time
	)

	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, currency, balance_amount, status, account_type, version, created_at                       
                 FROM accounts WHERE account_type = 'SystemAccount' AND currency = $1`, currency.Code()).
		Scan(&accID, &userID, &currencyCode, &balance, &status, &accType, &version, &createdAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSystemAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	c, err := domain.NewCurrency(currencyCode)
	if err != nil {
		return nil, err
	}
	money := domain.NewMoney(balance, c)

	uid := domain.UserID("")
	if userID != nil {
		uid = domain.UserID(*userID)
	}

	account := domain.ReconstituteAccount(domain.AccountID(accID), uid, c, money, domain.AccountStatus(status), version, domain.AccountType(accType), createdAt)

	return account, nil
}

func nullString(s string) *string {

	if s == "" {
		return nil
	}

	return &s
}
