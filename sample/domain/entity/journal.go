package entity

import (
	"accounting/application/apperror"
	"accounting/domain/vo"
	"time"
)

type Journal struct {
	ID              string            `` //
	BusinessID      string            `` //
	Date            time.Time         `` //
	Description     string            `` //
	JournalType     string            `` //
	UserID          string            `` //
	JournalBalances []*JournalBalance `` //
}

type JournalRequest struct {
	GetOrderID  func() string
	BusinessID  string    `` //
	Date        time.Time `` //
	Description string    `` //
	JournalType string    `` //
	UserID      string    `` //
}

func NewJournal(req JournalRequest) (*Journal, error) {

	if req.GetOrderID == nil {
		return nil, apperror.GetOrderIDMustNotEmpty
	}

	journalID := req.GetOrderID()
	if journalID == "" {
		return nil, apperror.JournalIDMustNotEmpty
	}

	obj := Journal{
		ID:          journalID,
		BusinessID:  req.BusinessID,
		Date:        req.Date,
		Description: req.Description,
		JournalType: req.JournalType,
		UserID:      req.UserID,
	}

	return &obj, nil
}

func (r *Journal) Add(jb *JournalBalance) error {
	r.JournalBalances = append(r.JournalBalances, jb)
	return nil
}

func (r *Journal) ValidateJournalBalance(accountSideMap map[string]vo.AccountSide) error {

	var debet float64
	var credit float64

	for _, jb := range r.JournalBalances {

		side := accountSideMap[jb.AccountCode]

		if side == vo.ActivaAccountSideEnum {

			if jb.Amount > 0 {
				debet += jb.Amount

			} else if jb.Amount < 0 {
				credit += -jb.Amount

			}

		} else if side == vo.PassivaAccountSideEnum {

			if jb.Amount > 0 {
				credit += jb.Amount

			} else if jb.Amount < 0 {
				debet += -jb.Amount

			}
		}

	}

	if debet != credit {
		return apperror.JournalIsNotBalance
	}

	return nil
}
