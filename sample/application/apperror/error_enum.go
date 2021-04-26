package apperror

const (
	FailUnmarshalResponseBodyError ErrorType = "ER1000 Fail to unmarshal response body"  // used by controller
	ObjectNotFound                 ErrorType = "ER1001 Object %s is not found"           // used by injected repo in interactor
	UnrecognizedEnum               ErrorType = "ER1002 %s is not recognized %s enum"     // used by enum
	DatabaseNotFoundInContextError ErrorType = "ER1003 Database is not found in context" // used by repoimpl
	GetOrderIDMustNotEmpty         ErrorType = "ER1004 get order id must not empty"      //
	JournalMustNotEmpty            ErrorType = "ER1005 journal must not empty"           //
	NominalMustNotZero             ErrorType = "ER1006 nominal must not zero"            //
	JournalIDMustNotEmpty          ErrorType = "ER1007 journal id must not empty"        //
	JournalBalanceMustNotEmpty     ErrorType = "ER1008 journal balance must not empty"   //
	BalanceAmountMustNotZero       ErrorType = "ER1009 balance amount must not zero"     //
	JournalIsNotBalance            ErrorType = "ER1010 journal is not balance"           //
	SideIsNotDefined               ErrorType = "ER1011 side is not defined"              //
	GakBisaSaveJournal             ErrorType = "ER1012 gak bisa save journal"            //
	AssertionFailed                ErrorType = "ER1000 assertion failed"                 //
	MissingAccountCode             ErrorType = "ER1000 missing account code %v"          //
	StockQuantityIsZero            ErrorType = "ER1000 stock quantity is zero"           //
	NotEnoughSufficientQuantity    ErrorType = "ER1000 not enough sufficient quantity"   //
	PriceMustNotBelowZero          ErrorType = "ER1000 price must not below zero"        //
)
