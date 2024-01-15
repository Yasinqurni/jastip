package customerror

const (
	// Unique record that already exists
	CodeErrAlreadyExists string = "already_exists"
	// Invalid credentials
	CodeErrInvalidCredential string = "invalid_credential"
	// Invalid request
	CodeErrInvalidRequest string = "invalid_request"
	// Record not found
	CodeErrNotFound string = "not_found"
	// Invalid transaction type
	CodeErrInvalidTransactionType string = "invalid_transaction_type"
	// Insufficient balance
	CodeErrInsufficientBalance string = "insufficient balance"
	// Already approved
	CodeErrAlreadyApproved string = "already approved"
	// Already rejected
	CodeErrAlreadyRejected string = "already rejected"
	// Goldpay Trad is inactive
	CodeErrInactiveGoldpayTrade string = "inactive_goldpay_trade"
)

type Err struct {
	Code   string `json:"code"`
	Errors any    `json:"error"`
}

func (e *Err) Error() string {
	return e.Code
}
