package entity

import (
	"accounting/application/apperror"
	"time"
)

type Journal struct {
	ID          string    `` //
	BusinessID  string    `` //
	Date        time.Time `` //
	Description string    `` //
	JournalType string    `` //
	UserID      string    `` //
}

type JournalRequest struct {
	GetOrderID func() string
	BusinessID  string `` //
	Date        string `` //
	Description string `` //
	JournalType string `` //
	UserID      string `` //
}

func NewJournal(req JournalRequest) (*Journal, error) {

	date, err := time.Parse("", req.Date)
	if err != nil {
		return nil, err
	}

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
		Date:        date,
		Description: req.Description,
		JournalType: req.JournalType,
		UserID:      req.UserID,
	}

	return &obj, nil
}
