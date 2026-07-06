package idgen

import (
	"FinFlow/internal/ledger/domain"

	"github.com/google/uuid"
)

type UUIDGenerator struct{}

var _ domain.IDGenerator = (*UUIDGenerator)(nil)

func (UUIDGenerator) Generate() domain.TransactionID {
	return domain.TransactionID(uuid.NewString())
}
