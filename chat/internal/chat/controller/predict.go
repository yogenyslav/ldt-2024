package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	ch "github.com/yogenyslav/ldt-2024/chat/internal/chat/handler"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

func (ctrl *Controller) Predict(ctx context.Context, out chan<- ch.Response, cancel <-chan struct{}, queryID int64) {
	if err := ctrl.repo.UpdateResponse(ctx, queryID, model.ResponseDao{
		Status: shared.StatusProcessing,
	}); err != nil {
		out <- ch.Response{
			Err: err,
			Msg: "predict failed",
		}
		return
	}

	ctx, cancelf := context.WithCancel(ctx)
	defer cancelf()

	cnt := 0
	buff := strings.Builder{}
	for {
		select {
		case <-cancel:
			out <- ch.Response{
				Err: nil,
				Msg: "predict canceled",
			}
			if err := ctrl.repo.UpdateResponse(ctx, queryID, model.ResponseDao{
				Status: shared.StatusCanceled,
				Body:   buff.String(),
			}); err != nil {
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
			if err := ctrl.repo.UpdateResponse(ctx, queryID, model.ResponseDao{
				Status: shared.StatusSuccess,
				Body:   buff.String(),
			}); err != nil {
				out <- ch.Response{
					Err: err,
					Msg: "failed to save response",
				}
			}
			return
		default:
			cnt++
			time.Sleep(time.Second * 1)
			msg := fmt.Sprintf("chunk %d", cnt)
			out <- ch.Response{
				Err: nil,
				Msg: msg,
			}
			buff.WriteString(msg)
			if cnt >= 10 {
				cancelf()
			}
		}
	}
}
