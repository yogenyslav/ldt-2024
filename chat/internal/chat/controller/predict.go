package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	ch "github.com/yogenyslav/ldt-2024/chat/internal/chat/handler"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Predict get prediction for query.
func (ctrl *Controller) Predict(ctx context.Context, out chan<- ch.Response, cancel <-chan struct{}, queryID int64) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Predict",
		trace.WithAttributes(attribute.Int64("queryID", queryID)),
	)
	defer span.End()

	if err := ctrl.repo.UpdateResponse(ctx, queryID, model.ResponseDao{
		Status: shared.StatusProcessing,
	}); err != nil {
		out <- ch.Response{
			Err: err,
			Msg: "predict failed",
		}
		return
	}

	withCancel, finish := context.WithCancel(ctx)
	defer finish()

	cnt := 0
	buff := strings.Builder{}
	for {
		select {
		case <-cancel:
			out <- ch.Response{
				Err:    nil,
				Msg:    "predict canceled",
				Finish: true,
			}
			if err := ctrl.repo.UpdateResponse(ctx, queryID, model.ResponseDao{
				Status: shared.StatusCanceled,
				Body:   buff.String(),
			}); err != nil {
				out <- ch.Response{
					Err:    err,
					Msg:    "cancel failed",
					Finish: true,
				}
			}
			return
		case <-withCancel.Done():
			out <- ch.Response{
				Err:    nil,
				Msg:    "finished",
				Finish: true,
			}
			if err := ctrl.repo.UpdateResponse(ctx, queryID, model.ResponseDao{
				Status: shared.StatusSuccess,
				Body:   buff.String(),
			}); err != nil {
				out <- ch.Response{
					Err:    err,
					Msg:    "failed to save response",
					Finish: true,
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
				finish()
			}
		}
	}
}
