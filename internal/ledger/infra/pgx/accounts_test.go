package pgx

import (
	"FinFlow/internal/ledger/app"
	"FinFlow/internal/ledger/domain"
	"FinFlow/internal/ledger/infra/idgen"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestAccounts_SaveHappyPath(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	pool, err := pgxpool.New(t.Context(), (os.Getenv("DATABASE_URL")))
	assert.Nil(t, err)

	uow := NewUnitOfWork(pool)

	accounts := NewAccountRepository(pool)
	transactions := NewTransactionRepository(pool)

	ids := idgen.UUIDGenerator{}

	currency, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-1", "u-1", currency, someTime)
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-2", "u-2", currency, someTime)
	assert.Nil(t, err)

	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)

	tc := app.TransferCommand{
		From:        "acc-1",
		To:          "acc-2",
		Amount:      500,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	service := app.NewLedgerService(accounts, transactions, ids)

	tx, events, err := service.Transfer(t.Context(), tc)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
	assert.NotNil(t, events)
}
