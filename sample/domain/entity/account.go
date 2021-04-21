package entity

import "accounting/domain/vo"

type Account struct {
	ID   string         `` //
	Side vo.AccountSide `` //
}

type AccountRequest struct {
}

func NewAccount(req AccountRequest) (*Account, error) {

	var obj Account

	return &obj, nil
}
