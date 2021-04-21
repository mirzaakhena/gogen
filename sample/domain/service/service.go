package service

import (
	"accounting/application/apperror"
	"accounting/domain/vo"
	"context"
)

type GenerateUUIDService interface {
	GenerateUUID(ctx context.Context) string
}

type JournalBalanceChecker struct {
	debet  float64
	credit float64
}

func NewJournalBalanceChecker() *JournalBalanceChecker {
	return &JournalBalanceChecker{debet: 0, credit: 0}
}

func (b *JournalBalanceChecker) Add(amount float64, side vo.AccountSide) error {
	if amount == 0 {
		return apperror.BalanceAmountMustNotZero
	}

	if side == vo.ActivaAccountSideEnum {
		if amount > 0 {
			b.debet += amount

		} else if amount < 0 {
			b.credit += -amount
		}
	} else if side == vo.PassivaAccountSideEnum {
		if amount > 0 {
			b.credit += amount

		} else if amount < 0 {
			b.debet += -amount
		}
	} else {
		return apperror.SideIsNotDefined
	}

	return nil
}

func (b *JournalBalanceChecker) IsBalance() bool {
	return b.debet == b.credit
}
