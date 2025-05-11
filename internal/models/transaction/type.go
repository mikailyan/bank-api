package transaction

type Type string

const (
	DEPOSIT    Type = "DEPOSIT"
	WITHDRAWAL Type = "WITHDRAWAL"
	TRANSFER   Type = "TRANSFER"
)
