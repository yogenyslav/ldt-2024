package shared

var (
	// UsernameKey is a key for username in context.
	UsernameKey = "x-username"
)

type ResponseStatus int8

const (
	_ ResponseStatus = iota
	StatusProcessing
	StatusSuccess
	StatusError
	StatusCanceled
)

type QueryType int8

const (
	_ QueryType = iota
	TypePrediction
	TypeStock
)
