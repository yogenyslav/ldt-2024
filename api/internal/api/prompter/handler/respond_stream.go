package handler

import (
	"context"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RespondStream обработка стрима ответов.
func (h *Handler) RespondStream(in *pb.StreamReq, out pb.Prompter_RespondStreamServer) error {
	ctx, span := h.tracer.Start(context.Background(), "Handler.RespondStream")
	defer span.End()

	stream, err := h.prompter.RespondStream(ctx, in)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to receive response")
			return status.Error(codes.Internal, err.Error())
		}

		if err := out.Send(resp); err != nil {
			log.Error().Err(err).Msg("failed to send response")
			return status.Error(codes.Internal, err.Error())
		}
	}
}
