package handler

type chatController interface {
}

// Handler is the chat handler
type Handler struct {
	ctrl chatController
}

// New creates a new chat handler
func New(ctrl chatController) *Handler {
	return &Handler{ctrl: ctrl}
}
