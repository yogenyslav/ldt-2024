package handler

import (
	"github.com/gofiber/contrib/websocket"
	chatresp "github.com/yogenyslav/ldt-2024/chat/pkg/chat_response"
)

// ProcessChunks обрабатывает чанки ответов.
func (h *Handler) ProcessChunks(out <-chan chatresp.Response, c *websocket.Conn) {
	for {
		select {
		case chunk := <-out:
			chatresp.Respond(c, chunk)
			if chunk.Finish {
				return
			}
		}
	}
}
