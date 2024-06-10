package controller

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (ctrl *Controller) FindOne(ctx context.Context, id uuid.UUID, username string) (model.FindOneResp, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.FindContent",
		trace.WithAttributes(attribute.String("id", id.String())),
	)
	defer span.End()

	resp := model.FindOneResp{
		ID:       id,
		Editable: false,
	}

	meta, err := ctrl.repo.FindMeta(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp, shared.ErrNoSessionWithID
		}
		return resp, shared.ErrGetSession
	}
	if meta.IsDeleted {
		return resp, shared.ErrNoSessionWithID
	}
	if meta.Username == username {
		resp.Editable = true
	}
	resp.Tg = meta.Tg

	contentDB, err := ctrl.repo.FindContent(ctx, id)
	if err != nil {
		return resp, shared.ErrGetSession
	}

	content := make([]model.SessionContentDto, len(contentDB))
	for i := 0; i < len(contentDB); i++ {
		content[i].Query = contentDB[i].Query.ToDto()
		content[i].Response = contentDB[i].Response.ToDto()
	}
	resp.Content = content
	return resp, nil
}
