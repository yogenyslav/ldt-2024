package handler

type chatController interface {
}

type Handler struct {
	ctrl chatController
}

func New(ctrl chatController) *Handler {
	return &Handler{ctrl: ctrl}
}
