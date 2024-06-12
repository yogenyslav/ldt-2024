package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/internal/api/prompter/model"
	"github.com/yogenyslav/ldt-2024/api/pkg"
)

// Extract вызывает метод Extract у Prompter.
func (ctrl *Controller) Extract(ctx context.Context, params model.ExtractReq) (*pb.ExtractedPrompt, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.Extract")
	defer span.End()

	ctx = pkg.PushSpan(ctx, span)
	return ctrl.prompter.Extract(ctx, &pb.ExtractReq{
		Prompt: params.Prompt,
	})
}
