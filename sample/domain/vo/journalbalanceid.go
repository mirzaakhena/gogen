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

	d := ""
	d += req.Date.Format("06")
	d += req.Date.Format("01")
	d += req.Date.Format("02")
	d += req.Date.Format("15")
	d += req.Date.Format("04")
	d += req.Date.Format("05")

	format := fmt.Sprintf("%v-%v-%v-%d", req.BusinessID, req.AccountCode, d, req.Sequence)
	obj := JournalBalanceID(format)

	return obj, nil
}

func (r JournalBalanceID) String() string {
	return string(r)
}
