package service

import (
	"codebase-app/internal/module/wac/entity"
	"codebase-app/internal/module/wac/ports"
	"codebase-app/pkg/errmsg"
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type wacService struct {
	repo ports.WACRepository
}

var privateFolder = "storage/private"

var _ ports.WACService = &wacService{}

func NewWACService(repo ports.WACRepository) *wacService {
	return &wacService{
		repo: repo,
	}
}

func (s *wacService) CreateWAC(ctx context.Context, req *entity.CreateWACRequest) error {
	errs := errmsg.NewCustomErrors(http.StatusBadRequest)

	for i, vc := range req.VehicleConditions {
		idx := strconv.Itoa(i)

		fileContent, err := base64.StdEncoding.DecodeString(vc.Image)
		if err != nil {
			log.Error().Err(err).Msg("service::CreateWAC - Failed to decode base64 image")
			errs.Add("vehicle_conditions["+idx+"].image", "gagal mendecode gambar.")
			continue
		}

		mimeType := detectMimeType(fileContent)

		if !acceptMimeType(mimeType) {
			log.Error().Err(err).Msg("service::CreateWAC - Invalid image format")
			errs.Add("vehicle_conditions["+idx+"].image", "format gambar tidak valid.")
			continue
		}

		ext := extensionFromMimeType(mimeType)

		filename := ulid.Make().String() + "." + ext
		filePath := filepath.Join(privateFolder, filename)
		err = os.WriteFile(filePath, fileContent, 0644)
		if err != nil {
			log.Error().Err(err).Msg("service::CreateWAC - Failed to write file")
			errs.Add("vehicle_conditions["+idx+"].image", "gambar gagal disimpan.")
			continue
		}
	}

	if errs.HasErrors() {
		return errs
	}

	return s.repo.CreateWAC(ctx, req)
}

func detectMimeType(data []byte) string {
	mimeType := http.DetectContentType(data)

	return mimeType
}

func extensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return ""
	}
}

func acceptMimeType(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}
