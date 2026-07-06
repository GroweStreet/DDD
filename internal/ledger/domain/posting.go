package domain

type Posting struct {
	accountID AccountID
	money     Money
}

func NewPosting(id AccountID, money Money) (Posting, error) {

	if money.IsZero() {
		return Posting{}, ErrZeroPosting
	}

	return Posting{
		accountID: id,
		money:     money,
	}, nil
}

func ReconstitutePosting(accountID AccountID, money Money) Posting {
	return Posting{
		accountID: accountID,
		money:     money,
	}
}

func (p Posting) AccountID() AccountID {
	return p.accountID
}

func (p Posting) Money() Money {
	return p.money
}
