package service

import (
	integstorage "codebase-app/internal/integration/localstorage"

	"codebase-app/internal/module/wac/entity"
	"codebase-app/internal/module/wac/ports"
	"codebase-app/pkg/errmsg"
	"context"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type wacService struct {
	repo    ports.WACRepository
	storage integstorage.LocalStorageContract
}

const PRIVATE_FOLDER = "storage/private"

var _ ports.WACService = &wacService{}

func NewWACService(
	repo ports.WACRepository,
	s integstorage.LocalStorageContract) *wacService {
	return &wacService{
		repo:    repo,
		storage: s,
	}
}

func (s *wacService) CreateWAC(ctx context.Context, req *entity.CreateWACRequest) (entity.CreateWACResponse, error) {
	var (
		resp entity.CreateWACResponse
		errs = errmsg.NewCustomErrors(http.StatusBadRequest)
	)

	for i, vc := range req.VehicleConditions {
		idx := strconv.Itoa(i)
		key := "vehicle_conditions[" + idx + "].image"

		path, err := s.storage.Save(vc.Image, PRIVATE_FOLDER)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateWAC - Failed to save image")
			errs.Add(key, "gambar gagal disimpan.")
			continue
		}

		req.VehicleConditions[i].Path = path
	}

	if errs.HasErrors() {
		return resp, errs
	}

	return s.repo.CreateWAC(ctx, req)
}

func (s *wacService) GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error) {
	return s.repo.GetWACs(ctx, req)
}

func (s *wacService) GetWAC(ctx context.Context, req *entity.GetWACRequest) (entity.GetWACResponse, error) {
	return s.repo.GetWAC(ctx, req)
}

func (s *wacService) OfferWAC(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error) {
	var (
		resp entity.OfferWACResponse
	)

	isCreator, err := s.repo.IsWACCreator(ctx, req.UserId, req.Id)
	if err != nil {
		return resp, err
	}

	if !isCreator {
		return resp, errmsg.NewCustomErrors(403, errmsg.WithMessage("Anda bukan pembuat walk around check ini"))
	}

	isOffered, err := s.repo.IsWACOffered(ctx, req.Id)
	if err != nil {
		return resp, err
	}

	if isOffered {
		return resp, errmsg.NewCustomErrors(403, errmsg.WithMessage("Walk around check sudah ditawarkan"))
	}

	return s.repo.OfferWAC(ctx, req)
}
