package createjournal

import (
	"accounting/application/apperror"
	"accounting/domain/entity"
	"accounting/domain/service"
	"context"
)

//go:generate mockery --name Outport -output mocks/

type createJournalInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase CreateJournal
func NewUsecase(outputPort Outport) Inport {
	return &createJournalInteractor{
		outport: outputPort,
	}
}

// Execute the usecase CreateJournal
func (r *createJournalInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	journalObj, err := entity.NewJournal(entity.JournalRequest{
		GetOrderID: func() string {
			return r.outport.GenerateUUID(ctx)
		},
		BusinessID:  req.BusinessID,
		Date:        req.Date,
		Description: req.Description,
		JournalType: req.JournalType,
		UserID:      req.UserID,
	})
	if err != nil {
		return nil, err
	}

	if len(req.JournalBalances) == 0 {
		return nil, apperror.JournalBalanceMustNotEmpty
	}

	err = r.outport.SaveJournal(ctx, journalObj)
	if err != nil {
		return nil, err
	}

	bc := service.NewJournalBalanceChecker()

	for i, jb := range req.JournalBalances {

		lastJournalBalanceObj, err := r.outport.FindLastJournalBalance(ctx, req.BusinessID, jb.AccountCode)
		if err != nil {
			return nil, err
		}

		journalBalanceObj, err := entity.NewJournalBalance(entity.JournalBalanceRequest{
			Journal:            journalObj,
			LastJournalBalance: lastJournalBalanceObj,
			AccountCode:        jb.AccountCode,
			Nominal:            jb.Nominal,
			Sequence:           i + 1,
		})
		if err != nil {
			return nil, err
		}

		accountObj, err := r.outport.FindOneAccountByCode(ctx, jb.AccountCode)
		if err != nil {
			return nil, err
		}
		if accountObj == nil {
			return nil, apperror.ObjectNotFound.Var(accountObj)
		}

		err = bc.Add(jb.Nominal, accountObj.Side)
		if err != nil {
			return nil, err
		}

		err = r.outport.SaveJournalBalance(ctx, journalBalanceObj)
		if err != nil {
			return nil, err
		}

	}

	if !bc.IsBalance() {
		return nil, apperror.JournalIsNotBalance
	}

	return res, nil
}
