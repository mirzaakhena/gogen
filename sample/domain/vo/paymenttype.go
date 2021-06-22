package vo

import (
	"accounting/application/apperror"
	"strings"
)

type PaymentType string

const (
	DANAPaymentTypeEnum  PaymentType = "DANA"
	GOPAYPaymentTypeEnum PaymentType = "GOPAY"
	OVOPaymentTypeEnum   PaymentType = "OVO"
)

var enumPaymentType = map[PaymentType]PaymentTypeDetail{
	DANAPaymentTypeEnum:  {},
	GOPAYPaymentTypeEnum: {},
	OVOPaymentTypeEnum:   {},
}

type PaymentTypeDetail struct { //
}

func NewPaymentType(name string) (PaymentType, error) {
	name = strings.ToUpper(name)

	if _, exist := enumPaymentType[PaymentType(name)]; !exist {
		return "", apperror.UnrecognizedEnum.Var(name, "PaymentType")
	}

	return PaymentType(name), nil
}

func (r PaymentType) GetDetail() PaymentTypeDetail {
	return enumPaymentType[r]
}

func (r PaymentType) PossibleValues() []PaymentType {
	res := make([]PaymentType, len(enumPaymentType))
	for key, _ := range enumPaymentType {
		res = append(res, key)
	}
	return res
}
