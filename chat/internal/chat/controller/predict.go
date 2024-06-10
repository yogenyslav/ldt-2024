package controller

import (
	"context"
	"fmt"
	"time"

	ch "github.com/yogenyslav/ldt-2024/chat/internal/chat/handler"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

func (ctrl *Controller) Predict(ctx context.Context, out chan<- ch.Response, cancel <-chan struct{}, queryID int64) {
	if err := ctrl.repo.UpdateResponseStatus(ctx, queryID, shared.StatusProcessing); err != nil {
		out <- ch.Response{
			Err: err,
			Msg: "predict failed",
		}
		return
	}

	ctx, cancelf := context.WithCancel(ctx)
	defer cancelf()

	cnt := 0
	for {
		select {
		case <-cancel:
			out <- ch.Response{
				Err: nil,
				Msg: "predict canceled",
			}
			if err := ctrl.repo.UpdateResponseStatus(ctx, queryID, shared.StatusCanceled); err != nil {
				out <- ch.Response{
					Err: err,
					Msg: "cancel failed",
				}
			}
			return
		case <-ctx.Done():
			out <- ch.Response{
				Err: nil,
				Msg: "finished",
			}
			if err := ctrl.repo.UpdateResponseStatus(ctx, queryID, shared.StatusSuccess); err != nil {
				out <- ch.Response{
					Err: err,
					Msg: "cancel failed",
				}
			}
			return
		default:
			cnt++
			time.Sleep(time.Second * 1)
			out <- ch.Response{
				Err: nil,
				Msg: fmt.Sprintf("token %d", cnt),
			}
			if cnt >= 10 {
				cancelf()
			}
		}
	}
}
