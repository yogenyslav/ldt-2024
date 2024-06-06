package model

// QueryCreateReq is a struct for creating new query request.
type QueryCreateReq struct {
	Prompt  string `json:"prompt,omitempty" validate:"gte=5"`
	Command string `json:"command,omitempty"`
}
