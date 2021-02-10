package shared

const (
	NotAllowedOrderStatusTransitionError ErrorType = "E1021 Order Status Transition from %s to %s is not allowed"
	OrderlineMustNotEmptyError           ErrorType = "E1022 Orderline must not empty"
	InvalidPhoneNumberFormatError        ErrorType = "E1023 %s is not valid phone number format"
	BalanceIsNotEnoughError              ErrorType = "E1024 User balance is not enough to continue the payment"
)
