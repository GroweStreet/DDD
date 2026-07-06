package domain

import "time"

type AccountID string
type UserID string

const overdraft = 5000

type AccountStatus string

const (
	StatusActive AccountStatus = "Active"
	StatusFrozen AccountStatus = "Frozen"
	StatusClosed AccountStatus = "Closed"
)

type AccountType string

const (
	AccountUserType   AccountType = "UserAccount"
	AccountSystemType AccountType = "SystemAccount"
)

type Account struct {
	id          AccountID
	userID      UserID
	currency    Currency
	balance     Money
	status      AccountStatus
	version     int64
	accountType AccountType
	events      []DomainEvent
	createdAt   time.Time
}

func OpenUserAccount(id AccountID, userID UserID, currency Currency, created time.Time) (*Account, error) {

	if currency == (Currency{}) {
		return nil, ErrCurrencyRequired
	}

	if id == "" {
		return nil, ErrInvalidID
	}

	if userID == "" {
		return nil, ErrInvalidID
	}

	a := &Account{
		id:          id,
		userID:      userID,
		currency:    currency,
		balance:     Zero(currency),
		status:      StatusActive,
		accountType: AccountUserType,
		version:     0,
		createdAt:   created,
	}

	ao := AccountOpened{
		accountID: id,
		userID:    userID,
		currency:  currency,
		at:        created,
	}

	a.recordEvent(ao)

	return a, nil
}

func OpenSystemAccount(id AccountID, currency Currency, created time.Time) (*Account, error) {

	if currency == (Currency{}) {
		return nil, ErrCurrencyRequired
	}

	if id == "" {
		return nil, ErrInvalidID
	}

	a := &Account{
		id:          id,
		currency:    currency,
		balance:     Zero(currency),
		status:      StatusActive,
		accountType: AccountSystemType,
		version:     0,
	}

	ao := AccountOpened{
		accountID: id,
		currency:  currency,
		at:        created,
	}

	a.recordEvent(ao)

	return a, nil
}

func ReconstituteAccount(id AccountID, userID UserID, currency Currency, balance Money, status AccountStatus, version int64, accountType AccountType, createdAt time.Time) *Account {

	return &Account{
		id:          id,
		userID:      userID,
		currency:    currency,
		balance:     balance,
		status:      status,
		version:     version,
		accountType: accountType,
		createdAt:   createdAt,
	}
}

func (a *Account) Apply(posting Posting) error {

	var err error

	if a.status != StatusActive {
		return ErrNotActiveAccount
	}

	if a.currency != posting.Money().currency {
		return ErrDifferentCurrency
	}

	switch a.accountType {
	case AccountUserType:
		if a.balance.amount+posting.Money().amount < -overdraft {
			return ErrOverdraftLimit
		}
	}

	a.balance, err = a.balance.Add(posting.Money())
	if err != nil {
		return err
	}

	a.version++
	return nil
}

func (a *Account) Freeze(frozen time.Time) error {

	if a.status == StatusFrozen || a.status == StatusClosed {
		return ErrFreezeAccount
	}

	a.status = StatusFrozen
	a.version++

	af := AccountFrozen{
		accountID: a.id,
		at:        frozen,
	}

	a.recordEvent(af)
	return nil
}

func (a *Account) Unfreeze(unfreeze time.Time) error {
	if a.status == StatusActive || a.status == StatusClosed {
		return ErrUnfreezeAccount
	}

	a.status = StatusActive
	a.version++

	au := AccountUnfreeze{
		accountID: a.id,
		at:        unfreeze,
	}
	a.recordEvent(au)

	return nil
}

func (a *Account) Close(closed time.Time) error {

	if a.status == StatusClosed {
		return ErrClosedAccount
	}

	a.status = StatusClosed
	a.version++

	ac := AccountClosed{
		accountID: a.id,
		at:        closed,
	}

	a.recordEvent(ac)

	return nil
}

func (a *Account) ID() AccountID {
	return a.id
}

func (a *Account) UserID() UserID {
	return a.userID
}

func (a *Account) Balance() Money {
	return a.balance
}

func (a *Account) Status() AccountStatus {
	return a.status
}

func (a *Account) Currency() Currency {
	return a.currency
}

func (a *Account) Version() int64 {
	return a.version
}

func (a *Account) Type() AccountType {
	return a.accountType
}

func (a *Account) Created() time.Time {
	return a.createdAt
}

func (a *Account) PullEvents() []DomainEvent {
	out := a.events
	a.events = nil

	return out
}

func (a *Account) recordEvent(e DomainEvent) {
	a.events = append(a.events, e)
}
