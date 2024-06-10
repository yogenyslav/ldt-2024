package shared

var (
	// UsernameKey is a key for username in context.
	UsernameKey = "x-username"
)

type ResponseStatus int8

const (
	_ ResponseStatus = iota
	StatusCreated
	StatusProcessing
	StatusSuccess
	StatusError
	StatusCanceled
)

// ToString return string representation of ResponseStatus.
func (s ResponseStatus) ToString() string {
	switch s {
	case StatusCreated:
		return "CREATED"
	case StatusProcessing:
		return "PROCESSING"
	case StatusSuccess:
		return "SUCCCESS"
	case StatusError:
		return "ERROR"
	case StatusCanceled:
		return "CANCELED"
	}
	return "UNDEFINED"
}

type QueryType int8

const (
	_ QueryType = iota
	TypePrediction
	TypeStock
)

// ToString return string representation of QueryType.
func (t QueryType) ToString() string {
	switch t {
	case TypePrediction:
		return "PREDICTION"
	case TypeStock:
		return "STOCK"
	}
	return "UNDEFINED"
}

type QueryCommand string

const (
	CommandValid   = "valid"
	CommandInvalid = "invalid"
	CommandCancel  = "cancel"
)

type QueryStatus int8

const (
	_ = iota
	StatusPending
	StatusValid
	StatusInvalid
)

// ToString return string representation of QueryStatus.
func (s QueryStatus) ToString() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusValid:
		return "VALID"
	case StatusInvalid:
		return "INVALID"
	}
	return "UNDEFINED"
}
