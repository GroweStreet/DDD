package domain

import (
	"time"
)

type TransactionID string
type IdempotencyKey string

type Transaction struct {
	id          TransactionID
	postings    []Posting
	createdAt   time.Time
	key         IdempotencyKey
	description string
	events      []DomainEvent
}

func NewTransaction(id TransactionID, postings []Posting, created time.Time, key IdempotencyKey, description string) (*Transaction, error) {
	var err error
	cp := make([]Posting, len(postings))
	copy(cp, postings)

	if len(cp) < 2 {
		return nil, ErrNotEnoughPostings
	}

	transactionCurrency := cp[0].Money().currency
	var acc = Zero(transactionCurrency)

	for _, posting := range cp {
		if posting.Money().currency != transactionCurrency {
			return nil, ErrDifferentCurrency
		}

		acc, err = acc.Add(posting.money)
		if err != nil {
			return nil, err
		}
	}

	if !acc.IsZero() {
		return nil, ErrNotZeroSum
	}

	tt := TransactionPosted{
		transactionID: id,
		postings:      cp,
		currency:      transactionCurrency,
		at:            created,
	}

	transaction := Transaction{
		id:          id,
		postings:    cp,
		createdAt:   created,
		key:         key,
		description: description,
	}

	transaction.recordEvent(tt)
	return &transaction, nil
}

func ReconstituteTransaction(id TransactionID, postings []Posting, createdAt time.Time, key IdempotencyKey, description string) *Transaction {

	return &Transaction{
		id:          id,
		postings:    postings,
		createdAt:   createdAt,
		key:         key,
		description: description,
		events:      nil,
	}
}

func (t *Transaction) ID() TransactionID {
	return t.id
}

func (t *Transaction) Postings() []Posting {

	cp := make([]Posting, len(t.postings))
	copy(cp, t.postings)

	return cp
}

func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Transaction) Key() IdempotencyKey {
	return t.key
}

func (t *Transaction) Description() string {
	return t.description
}

func (t *Transaction) PullEvents() []DomainEvent {
	out := t.events
	t.events = nil

	return out
}

func (t *Transaction) recordEvent(e DomainEvent) {
	t.events = append(t.events, e)
}
