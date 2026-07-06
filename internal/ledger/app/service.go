package app

import (
	"FinFlow/internal/ledger/domain"
	"context"
	"errors"
	"slices"
	"time"
)

type TransferCommand struct {
	From        domain.AccountID
	To          domain.AccountID
	Amount      int64
	Key         domain.IdempotencyKey
	At          time.Time
	Description string
}

type DepositCommand struct {
	To          domain.AccountID
	Amount      int64
	Key         domain.IdempotencyKey
	At          time.Time
	Description string
}

type WithdrawCommand struct {
	From        domain.AccountID
	Amount      int64
	Key         domain.IdempotencyKey
	At          time.Time
	Description string
}

type LedgerService struct {
	accounts domain.AccountRepository
	txns     domain.TransactionRepository
	ids      domain.IDGenerator
	uow      UnitOfWork
}

func NewLedgerService(accounts domain.AccountRepository, txns domain.TransactionRepository, ids domain.IDGenerator, uow UnitOfWork) *LedgerService {
	return &LedgerService{
		accounts: accounts,
		txns:     txns,
		ids:      ids,
		uow:      uow,
	}
}

func (s *LedgerService) Transfer(ctx context.Context, cmd TransferCommand) (*domain.Transaction, []domain.DomainEvent, error) {

	var result *domain.Transaction
	var events []domain.DomainEvent

	err := s.uow.Do(ctx, func(r Repositories) error {

		t, err := r.Transactions.FindByIdempotencyKey(ctx, cmd.Key)

		if err != nil && !errors.Is(err, domain.ErrTransactionNotFound) {
			return err
		}

		if t != nil {
			result = t
			return nil
		}

		if cmd.Amount <= 0 {
			return domain.ErrInvalidAmount
		}

		if cmd.From == cmd.To {
			return domain.ErrSameAccount
		}

		from, err := r.Accounts.FindByID(ctx, cmd.From)
		if err != nil {
			return err
		}

		to, err := r.Accounts.FindByID(ctx, cmd.To)
		if err != nil {
			return err
		}

		if from.Currency() != to.Currency() {
			return domain.ErrMismatchCurrency
		}

		money := domain.NewMoney(cmd.Amount, from.Currency())

		tx, evs, err := s.post(ctx, r, from, to, money, cmd.At, cmd.Key, cmd.Description)

		if err != nil {
			return err
		}
		result, events = tx, evs
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return result, events, nil
}

func (s *LedgerService) Deposit(ctx context.Context, cmd DepositCommand) (*domain.Transaction, []domain.DomainEvent, error) {

	var result *domain.Transaction
	var events []domain.DomainEvent

	err := s.uow.Do(ctx, func(r Repositories) error {

		t, err := r.Transactions.FindByIdempotencyKey(ctx, cmd.Key)

		if err != nil && !errors.Is(err, domain.ErrTransactionNotFound) {
			return err
		}

		if t != nil {
			result = t
			return nil
		}

		if cmd.Amount <= 0 {
			return domain.ErrInvalidAmount
		}

		to, err := r.Accounts.FindByID(ctx, cmd.To)
		if err != nil {
			return err
		}

		if to.Type() != domain.AccountUserType {
			return domain.ErrNotUserAccount
		}

		from, err := r.Accounts.FindSystemAccount(ctx, to.Currency())
		if err != nil {
			return err
		}

		money := domain.NewMoney(cmd.Amount, to.Currency())

		tx, evs, err := s.post(ctx, r, from, to, money, cmd.At, cmd.Key, cmd.Description)

		if err != nil {
			return err
		}

		result, events = tx, evs
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return result, events, nil

}

func (s *LedgerService) Withdraw(ctx context.Context, cmd WithdrawCommand) (*domain.Transaction, []domain.DomainEvent, error) {

	var result *domain.Transaction
	var events []domain.DomainEvent

	err := s.uow.Do(ctx, func(r Repositories) error {

		t, err := r.Transactions.FindByIdempotencyKey(ctx, cmd.Key)

		if err != nil && !errors.Is(err, domain.ErrTransactionNotFound) {
			return err
		}

		if t != nil {
			result = t
			return nil
		}

		if cmd.Amount <= 0 {
			return domain.ErrInvalidAmount
		}

		from, err := r.Accounts.FindByID(ctx, cmd.From)
		if err != nil {
			return err
		}

		if from.Type() != domain.AccountUserType {
			return domain.ErrNotUserAccount
		}

		to, err := r.Accounts.FindSystemAccount(ctx, from.Currency())
		if err != nil {
			return err
		}

		money := domain.NewMoney(cmd.Amount, from.Currency())

		tx, evs, err := s.post(ctx, r, from, to, money, cmd.At, cmd.Key, cmd.Description)
		if err != nil {
			return err
		}

		result, events = tx, evs
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return result, events, nil

}

func (s *LedgerService) post(ctx context.Context, r Repositories, from, to *domain.Account, money domain.Money, at time.Time, key domain.IdempotencyKey, desc string) (*domain.Transaction, []domain.DomainEvent, error) {

	postingFrom, err := domain.NewPosting(from.ID(), money.Negate())
	if err != nil {
		return nil, nil, err
	}

	postingTo, err := domain.NewPosting(to.ID(), money)
	if err != nil {
		return nil, nil, err
	}

	if err = from.Apply(postingFrom); err != nil {
		return nil, nil, err
	}

	if err = to.Apply(postingTo); err != nil {
		return nil, nil, err
	}

	tx, err := domain.NewTransaction(s.ids.Generate(), []domain.Posting{postingFrom, postingTo}, at, key, desc)
	if err != nil {
		return nil, nil, err
	}

	if err = r.Accounts.Save(ctx, from); err != nil {
		return nil, nil, err
	}

	if err = r.Accounts.Save(ctx, to); err != nil {
		return nil, nil, err
	}

	if err = r.Transactions.Save(ctx, tx); err != nil {
		return nil, nil, err
	}

	eventsFrom := from.PullEvents()
	eventsTo := to.PullEvents()
	eventsTxns := tx.PullEvents()

	events := slices.Concat(eventsFrom, eventsTo, eventsTxns)

	return tx, events, nil
}
