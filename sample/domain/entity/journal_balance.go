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
	ID          vo.JournalBalanceID `` //
	AccountCode string              `` //
	Amount      float64             `` //
	Balance     float64             `` //
	Direction   Direction           `` //
	Sequence    int                 `` //
}

type JournalBalanceRequest struct {
	Journal     *Journal `` //
	LastBalance float64  `` //
	AccountCode string   `` //
	Nominal     float64  `` //
	Sequence    int      `` //
}

func NewJournalBalance(req JournalBalanceRequest) (*JournalBalance, error) {

	if req.Journal == nil {
		return nil, apperror.JournalMustNotEmpty
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
		ID:          jbID,
		AccountCode: req.AccountCode,
		Amount:      req.Nominal,
		Balance:     req.LastBalance + req.Nominal,
		Direction:   direction,
		Sequence:    req.Sequence,
	}

	return &obj, nil
}
