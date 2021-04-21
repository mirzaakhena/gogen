package apperror

const (
	FailUnmarshalResponseBodyError ErrorType = "ER1000 Fail to unmarshal response body"  // used by controller
	ObjectNotFound                 ErrorType = "ER1001 Object %s is not found"           // used by injected repo in interactor
	UnrecognizedEnum               ErrorType = "ER1002 %s is not recognized %s enum"     // used by enum
	DatabaseNotFoundInContextError ErrorType = "ER1003 Database is not found in context" // used by repoimpl
	GetOrderIDMustNotEmpty         ErrorType = "ER1000 get order id must not empty"      //
	JournalMustNotEmpty            ErrorType = "ER1000 journal must not empty"           //
	NominalMustNotZero             ErrorType = "ER1000 nominal must not zero"            //
	JournalIDMustNotEmpty          ErrorType = "ER1000 journal id must not empty"        //
	JournalBalanceMustNotEmpty     ErrorType = "ER1000 journal balance must not empty"   //
	BalanceAmountMustNotZero       ErrorType = "ER1000 balance amount must not zero"     //
	JournalIsNotBalance            ErrorType = "ER1000 journal is not balance"           //
	SideIsNotDefined               ErrorType = "ER1000 side is not defined"              //
	GakBisaSaveJournal             ErrorType = "ER1000 gak bisa save journal"            //
)
