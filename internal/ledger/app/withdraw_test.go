package app

import (
	"FinFlow/internal/ledger/domain"
	"FinFlow/internal/ledger/infra/memory"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithdrawHappyPath(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", currency, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenSystemAccount("acc-2", currency, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	wc := WithdrawCommand{
		From:        "acc-1",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	dc := DepositCommand{
		To:          "acc-1",
		Amount:      1000,
		Key:         "key-2",
		At:          someTime,
		Description: "",
	}

	service := NewLedgerService(accounts, txns, &fakeIDGen{})
	tx, events, err := service.Deposit(t.Context(), dc)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, 1, len(events))

	tx, events, err = service.Withdraw(t.Context(), wc)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, 1, len(events))

	acc, err := accounts.FindByID(t.Context(), "acc-1")
	assert.NotNil(t, acc)
	assert.Equal(t, int64(0), acc.Balance().Amount())
}

func TestWithdraw_IdempotencyKey(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", currency, someTime)
	assert.Nil(t, err)

	system, err := domain.OpenSystemAccount("acc-2", currency, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), system)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	wc := WithdrawCommand{
		From:        "acc-1",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	tx, events, err := service.Withdraw(t.Context(), wc)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, 1, len(events))

	wc = WithdrawCommand{
		From:        "acc-1",
		Amount:      2000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	tx, events, err = service.Withdraw(t.Context(), wc)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
	assert.Nil(t, events)

	gotFrom, _ := accounts.FindByID(t.Context(), "acc-1")
	assert.Equal(t, int64(-1000), gotFrom.Balance().Amount())
}

func TestWithdraw_InvalidAmount(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", currency, someTime)
	assert.Nil(t, err)

	system, err := domain.OpenSystemAccount("acc-2", currency, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), system)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	wc := WithdrawCommand{
		From:        "acc-1",
		Amount:      -1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Withdraw(t.Context(), wc)
	assert.ErrorIs(t, err, domain.ErrInvalidAmount)
	assert.Equal(t, 0, len(events))
}

func TestWithdraw_NotUserAccount(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenSystemAccount("acc-1", currency, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	wc := WithdrawCommand{
		From:        "acc-1",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Withdraw(t.Context(), wc)
	assert.ErrorIs(t, err, domain.ErrNotUserAccount)
	assert.Equal(t, 0, len(events))
}

func TestWithdraw_FromNotFound(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	wc := WithdrawCommand{
		From:        "acc-1",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Withdraw(t.Context(), wc)
	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Equal(t, 0, len(events))
}

func TestWithdraw_SystemAccountNotFound(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", currency, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	wc := WithdrawCommand{
		From:        "acc-1",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Withdraw(t.Context(), wc)
	assert.ErrorIs(t, err, domain.ErrSystemAccountNotFound)
	assert.Equal(t, 0, len(events))
}

func TestWithdraw_OverdraftLimit(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", currency, someTime)
	assert.Nil(t, err)

	system, err := domain.OpenSystemAccount("acc-2", currency, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), system)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	wc := WithdrawCommand{
		From:        "acc-1",
		Amount:      6000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Withdraw(t.Context(), wc)
	assert.ErrorIs(t, err, domain.ErrOverdraftLimit)
	assert.Equal(t, 0, len(events))
}
