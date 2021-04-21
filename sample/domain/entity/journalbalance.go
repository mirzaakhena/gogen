package entity

import (
	"accounting/application/apperror"
	"accounting/domain/vo"
)

type Direction string

const (
	DirectionIncreased = "INCREASED"
	DirectionReduced   = "REDUCED"
)

type JournalBalance struct {
	ID                 vo.JournalBalanceID `` //
	BusinessID         string              `` //
	JournalID          string              `` //
	LastJournalBalance *JournalBalance     `` //
	AccountCode        string              `` //
	Amount             float64             `` //
	Balance            float64             `` //
	Direction          Direction           `` //
	Sequence           int                 `` //
}

type JournalBalanceRequest struct {
	Journal            *Journal        `` //
	LastJournalBalance *JournalBalance `` //
	AccountCode        string          `` //
	Nominal            float64         `` //
	Sequence           int             `` //
}

func NewJournalBalance(req JournalBalanceRequest) (*JournalBalance, error) {

	if req.Journal == nil {
		return nil, apperror.JournalMustNotEmpty
	}

	lastBalance := 0.0
	if req.LastJournalBalance != nil {
		lastBalance = req.LastJournalBalance.Balance
	}

	direction := Direction("")

	if req.Nominal > 0 {
		direction = DirectionIncreased

	} else if req.Nominal < 0 {
		direction = DirectionReduced

	} else {
		return nil, apperror.NominalMustNotZero

	}

	jbID, err := vo.NewJournalBalanceID(vo.JournalBalanceIDRequest{
		BusinessID:  req.Journal.BusinessID,
		AccountCode: req.AccountCode,
		Date:        req.Journal.Date,
		Sequence:    req.Sequence,
	})
	if err != nil {
		return nil, err
	}

	obj := JournalBalance{
		ID:                 jbID,
		BusinessID:         req.Journal.BusinessID,
		JournalID:          req.Journal.ID,
		LastJournalBalance: req.LastJournalBalance,
		AccountCode:        req.AccountCode,
		Amount:             req.Nominal,
		Balance:            lastBalance + req.Nominal,
		Direction:          direction,
		Sequence:           req.Sequence,
	}

	return &obj, nil
}
