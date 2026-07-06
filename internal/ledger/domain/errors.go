package domain

import "errors"

var ErrUnknownCurrency = errors.New("currency does not exist")

var ErrCurrencyRequired = errors.New("currency is required")

var ErrMismatchingMoneyTypes = errors.New("money type is different")

var ErrZeroPosting = errors.New("amount cannot be zero")

var ErrNotEnoughPostings = errors.New("postings cant be less than 2")

var ErrDifferentCurrency = errors.New("different currency in transaction")

var ErrNotZeroSum = errors.New("sum of posting not equal to zero")

var ErrNotActiveAccount = errors.New("account is not active")

var ErrOverdraftLimit = errors.New("hit overdraft limit")

var ErrFreezeAccount = errors.New("account cant be freeze")

var ErrUnfreezeAccount = errors.New("cant unfreeze active or closed account")

var ErrClosedAccount = errors.New("cant close closed account")

var ErrInvalidID = errors.New("invalid account id")

var ErrSameAccount = errors.New("can not transfer to same account")

var ErrMismatchCurrency = errors.New("mismatched currency")

var ErrNotFound = errors.New("data does not found")

var ErrInvalidAmount = errors.New("amount cant be negative or zero")

var ErrSystemAccountNotFound = errors.New("account is not found")

var ErrNotUserAccount = errors.New("not a user account")

var ErrVersionConflict = errors.New("version conflict")

var ErrAccountNotFound = errors.New("account is not found")

var ErrTransactionNotFound = errors.New("transaction is not found")

var ErrTransactionExits = errors.New("transaction is already exist")
