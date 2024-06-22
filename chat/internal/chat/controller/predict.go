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
func (ctrl *Controller) Predict(ctx context.Context, out chan<- chatresp.Response, cancel <-chan struct{}, queryID int64, org string) {
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
		out <- chatresp.Response{Err: err.Error(), Msg: "update response to processing failed"}
		return
	}

	meta, err := ctrl.cr.FindQueryMeta(ctx, queryID)
	if err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "find query meta failed"}
		return
	}

	predict, err := ctrl.predictor.Predict(ctx, &pb.PredictReq{
		Type:         pb.QueryType(meta.Type),
		Product:      meta.Product,
		Period:       meta.Period,
		Organization: org,
	})
	if err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "predict failed"}
		return
	}

	if err = ctrl.cr.UpdateResponseData(ctx, model.ResponseDao{
		QueryID:  queryID,
		Data:     predict.GetData(),
		DataType: meta.Type,
	}); err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "update response data failed"}
		return
	}

	data := make(map[string]any)
	if err = json.Unmarshal(predict.GetData(), &data); err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "failed to unmarshal response data"}
		return
	}

	out <- chatresp.Response{Msg: "predict succeeded", Data: data, DataType: meta.Type.ToString()}
	if meta.Type == shared.TypeStock {
		if err = ctrl.cr.UpdateResponse(ctx, model.ResponseDao{
			QueryID: queryID,
			Status:  shared.StatusSuccess,
		}); err != nil {
			out <- chatresp.Response{Err: err.Error(), Msg: "update response failed"}
		}
		return
	}

	err = ctrl.respondStream(ctx, out, cancel, predict.GetData(), queryID)
}

func (ctrl *Controller) respondStream(ctx context.Context, out chan<- chatresp.Response, cancel <-chan struct{}, prompt []byte, queryID int64) error {
	withCancel, finish := context.WithCancel(ctx)
	defer finish()

	in := &pb.StreamReq{Prompt: prompt}
	stream, err := ctrl.prompter.RespondStream(withCancel, in)
	if err != nil {
		out <- chatresp.Response{Err: err.Error(), Msg: "failed to respond stream", Finish: true}
		return err
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
				return err
			}
			return nil
		case <-withCancel.Done():
			out <- chatresp.Response{Msg: "finished", Finish: true}
			if err := ctrl.cr.UpdateResponse(ctx, model.ResponseDao{
				QueryID: queryID,
				Status:  shared.StatusSuccess,
				Body:    buff.String(),
			}); err != nil {
				out <- chatresp.Response{Err: err.Error(), Msg: "failed to save response", Finish: true}
				return err
			}
			return nil
		default:
			resp, err := stream.Recv()
			if err == io.EOF {
				finish()
				break
			}
			if err != nil {
				out <- chatresp.Response{Err: err.Error(), Msg: "failed to receive response", Finish: true}
				return err
			}
			chunk := resp.GetChunk()
			out <- chatresp.Response{Data: struct {
				Info string `json:"info"`
			}{chunk}, Chunk: true}
			buff.WriteString(chunk)
		}
	}
}
