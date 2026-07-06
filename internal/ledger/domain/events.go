package domain

import "time"

type DomainEvent interface {
	OccurredAt() time.Time
}

type AccountOpened struct {
	accountID AccountID
	userID    UserID
	currency  Currency
	at        time.Time
}

func (a AccountOpened) OccurredAt() time.Time {
	return a.at
}

type AccountFrozen struct {
	accountID AccountID
	at        time.Time
}

func (a AccountFrozen) OccurredAt() time.Time {
	return a.at
}

type AccountUnfreeze struct {
	accountID AccountID
	at        time.Time
}

func (a AccountUnfreeze) OccurredAt() time.Time {
	return a.at
}

type AccountClosed struct {
	accountID AccountID
	at        time.Time
}

func (a AccountClosed) OccurredAt() time.Time {
	return a.at
}

type TransactionPosted struct {
	transactionID TransactionID
	postings      []Posting
	currency      Currency
	at            time.Time
}

func (t TransactionPosted) OccurredAt() time.Time {
	return t.at
}

type MoneyTransferred struct {
	transactionID TransactionID
	from          UserID
	to            UserID
	amount        Money
	at            time.Time
}

func (m MoneyTransferred) OccurredAt() time.Time {
	return m.at
}

type MoneyDeposited struct {
	transactionID TransactionID
	accountID     AccountID
	amount        Money
	at            time.Time
}

func (m MoneyDeposited) OccurredAt() time.Time {
	return m.at
}

type MoneyWithdrawn struct {
	transactionID TransactionID
	accountID     AccountID
	amount        Money
	at            time.Time
}

func (m MoneyWithdrawn) OccurredAt() time.Time {
	return m.at
}
