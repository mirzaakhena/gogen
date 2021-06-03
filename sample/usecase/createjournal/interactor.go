package createjournal

import (
	"accounting/application/apperror"
	"accounting/domain/entity"
	"accounting/domain/repository"
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

	err := repository.ExecuteTransaction(ctx, r.outport, func(ctx context.Context) (err error) {

		if len(req.JournalBalances) == 0 {
			return apperror.JournalBalanceMustNotEmpty
		}

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
			return err
		}

		// prepare last account balance
		accountCodes := make([]string, 0)

		for _, jb := range req.JournalBalances {
			accountCodes = append(accountCodes, jb.AccountCode)
		}

		journalBalanceMap, err := r.outport.FindAllLastJournalBalance(ctx, req.BusinessID, accountCodes)
		if err != nil {
			return err
		}

		for i, jb := range req.JournalBalances {

			if _, exist := journalBalanceMap[jb.AccountCode]; !exist {
				return apperror.MissingAccountCode.Var(jb.AccountCode)
			}

			journalBalanceObj, err := entity.NewJournalBalance(entity.JournalBalanceRequest{
				Journal:     journalObj,
				LastBalance: journalBalanceMap[jb.AccountCode],
				AccountCode: jb.AccountCode,
				Nominal:     jb.Nominal,
				Sequence:    i + 1,
			})
			if err != nil {
				return err
			}

			err = journalObj.Add(journalBalanceObj)
			if err != nil {
				return err
			}

		}

		// Check journal balance
		{
			accountObjMap, err := r.outport.FindAllAccountSideByCodes(ctx, req.BusinessID, accountCodes)
			if err != nil {
				return err
			}

			for _, code := range accountCodes {
				if _, exist := accountObjMap[code]; !exist {
					return apperror.MissingAccountCode.Var(code)
				}
			}

			err = journalObj.ValidateJournalBalance(accountObjMap)
			if err != nil {
				return err
			}
		}

		err = r.outport.SaveJournal(ctx, journalObj)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
