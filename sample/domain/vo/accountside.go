package vo

import (
	"accounting/application/apperror"
	"strings"
)

type AccountSide string

const (
	ActivaAccountSideEnum  AccountSide = "ACTIVA"
	PassivaAccountSideEnum AccountSide = "PASSIVA"
)

var enumAccountSide = map[AccountSide]AccountSideDetail{
	ActivaAccountSideEnum:  {},
	PassivaAccountSideEnum: {},
}

type AccountSideDetail struct { //
}

func NewAccountSide(name string) (AccountSide, error) {
	name = strings.ToUpper(name)

	if _, exist := enumAccountSide[AccountSide(name)]; !exist {
		return "", apperror.UnrecognizedEnum.Var(name, "AccountSide")
	}

	return AccountSide(name), nil
}

func (r AccountSide) GetDetail() AccountSideDetail {
	return enumAccountSide[r]
}

func (r AccountSide) PossibleValues() []AccountSide {
	res := make([]AccountSide, len(enumAccountSide))
	for key, _ := range enumAccountSide {
		res = append(res, key)
	}
	return res
}
