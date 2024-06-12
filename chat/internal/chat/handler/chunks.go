package handler

import (
	"github.com/gofiber/contrib/websocket"
)

// ProcessChunks обрабатывает чанки ответов.
func (h *Handler) ProcessChunks(out <-chan Response, c *websocket.Conn) {
	for {
		select {
		case chunk := <-out:
			respond(c, chunk)
			if chunk.Finish {
				return
			}
		}
	}
}
