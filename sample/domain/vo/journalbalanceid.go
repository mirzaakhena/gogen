package vo

import (
	"fmt"
	"time"
)

type JournalBalanceID string

type JournalBalanceIDRequest struct {
	BusinessID  string
	AccountCode string
	Date        time.Time
	Sequence    int
}

func NewJournalBalanceID(req JournalBalanceIDRequest) (JournalBalanceID, error) {

	d := req.Date.Format("060102150405")
	format := fmt.Sprintf("%v-%v-%v-%d", req.BusinessID, req.AccountCode, d, req.Sequence)
	obj := JournalBalanceID(format)

	return obj, nil
}

func (r JournalBalanceID) String() string {
	return string(r)
}
