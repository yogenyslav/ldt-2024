package controller

import (
	"archive/zip"
	"context"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/api/pb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ImportData импортирует данные из архива.
func (ctrl *Controller) ImportData(ctx context.Context, mpArchive *multipart.FileHeader, id int64) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.ImportData",
		trace.WithAttributes(attribute.Int64("organizationID", id)),
	)
	defer span.End()

	rawArchive, err := mpArchive.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = rawArchive.Close()
	}()

	archive, err := zip.NewReader(rawArchive, mpArchive.Size)
	if err != nil {
		return err
	}

	organizationTitle := getOrganizationTitle(id)

	sources := make([]*pb.Source, 0, len(archive.File))
	for _, file := range archive.File {
		if file.FileInfo().IsDir() {
			continue
		}
		log.Debug().Str("file", file.FileInfo().Name()).Msg("processing file")
		f, err := file.Open()
		if err != nil {
			return err
		}

		if _, err := ctrl.s3.PutObject(
			ctx,
			organizationTitle,
			file.FileInfo().Name(),
			f,
			file.FileInfo().Size(),
			minio.PutObjectOptions{},
		); err != nil {
			log.Error().Err(err).Str("file", file.FileInfo().Name()).Msg("failed to put object")
			return err
		}

		url, err := ctrl.s3.PresignedGetObject(ctx, getOrganizationTitle(id), file.FileInfo().Name(), time.Hour, nil)
		if err != nil {
			log.Error().Err(err).Str("file", file.FileInfo().Name()).Msg("failed to get presigned url")
			return err
		}
		sources = append(sources, &pb.Source{
			Name: file.FileInfo().Name(),
			Path: url.String(),
		})
	}

	in := &pb.PrepareDataReq{
		Sources:      sources,
		Organization: organizationTitle,
	}
	if _, err := ctrl.predictor.PrepareData(ctx, in); err != nil {
		log.Error().Err(err).Msg("failed to prepare data")
		return err
	}
	return nil
}
