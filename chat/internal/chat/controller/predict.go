package controller

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	chatresp "github.com/yogenyslav/ldt-2024/chat/pkg/chat_response"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Predict получить предикт по запросу.
//
//nolint:funlen // let it be long
func (ctrl *Controller) Predict(ctx context.Context, out chan<- chatresp.Response, cancel <-chan struct{}, queryID int64) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Predict",
		trace.WithAttributes(attribute.Int64("queryID", queryID)),
	)
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			if err = ctrl.cr.UpdateResponse(ctx, model.ResponseDao{
				QueryID: queryID,
				Status:  shared.StatusError,
			}); err != nil {
				out <- chatresp.Response{Err: err.Error(), Msg: "failed to set response to error"}
			}
		}
	}()

	if err = ctrl.cr.UpdateResponse(ctx, model.ResponseDao{
		QueryID: queryID,
		Status:  shared.StatusProcessing,
	}); err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "predict failed"}
		return
	}

	meta, err := ctrl.cr.FindQueryMeta(ctx, queryID)
	if err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "predict failed"}
		return
	}

	predict, err := ctrl.predictor.Predict(ctx, &pb.PredictReq{
		Type:    pb.QueryType(meta.Type),
		Product: meta.Product,
		Period:  meta.Period,
	})
	if err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "predict failed"}
		return
	}

	data := make(map[string]any)
	if err = json.Unmarshal(predict.GetData(), &data); err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "predict failed"}
		return
	}

	out <- chatresp.Response{Msg: "predict succeeded", Data: data, DataType: meta.Type}
	ctrl.respondStream(ctx, out, cancel, predict.GetData(), queryID)
}

func (ctrl *Controller) respondStream(ctx context.Context, out chan<- chatresp.Response, cancel <-chan struct{}, prompt []byte, queryID int64) {
	withCancel, finish := context.WithCancel(ctx)
	defer finish()

	in := &pb.StreamReq{Prompt: prompt}
	stream, err := ctrl.prompter.RespondStream(withCancel, in)
	if err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "failed to respond stream", Finish: true}
		return
	}

	buff := strings.Builder{}
	for {
		select {
		case <-cancel:
			out <- chatresp.Response{Msg: "predict canceled", Finish: true}
			if err := ctrl.cr.UpdateResponse(ctx, model.ResponseDao{
				QueryID: queryID,
				Status:  shared.StatusCanceled,
				Body:    buff.String(),
			}); err != nil {
				out <- chatresp.Response{Err: err.Error(), Msg: "cancel failed", Finish: true}
			}
			return
		case <-withCancel.Done():
			out <- chatresp.Response{Msg: "finished", Finish: true}
			if err := ctrl.cr.UpdateResponse(ctx, model.ResponseDao{
				QueryID: queryID,
				Status:  shared.StatusSuccess,
				Body:    buff.String(),
			}); err != nil {
				out <- chatresp.Response{Err: err.Error(), Msg: "failed to save response", Finish: true}
			}
			return
		default:
			resp, err := stream.Recv()
			if err == io.EOF {
				finish()
			}
			if err != nil {
				out <- chatresp.Response{Err: err.Error(), Msg: "failed to receive response", Finish: true}
				return
			}
			chunk := resp.GetChunk()
			out <- chatresp.Response{Data: struct {
				Info string `json:"info"`
			}{chunk}, Chunk: true}
			buff.WriteString(chunk)
		}
	}
}
