package repository

import (
	"codebase-app/internal/module/common/entity"
	"codebase-app/pkg/errmsg"
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
)

func (r *commonRepo) GetHTIBrands(ctx context.Context) ([]entity.CommonResponse, error) {
	query := `
		SELECT DISTINCT
			brand as name
		FROM
			trade_in_trends
		ORDER BY
			brand ASC
		`

	var data = make([]entity.CommonResponse, 0)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query))
	if err != nil {
		log.Error().Err(err).Msg("repo::GetHTIBrands - Failed to get HTI brands")
		return nil, err
	}

	return data, nil
}

func (r *commonRepo) GetHTIModels(ctx context.Context, req *entity.GetHTIModelsRequest) ([]entity.CommonResponse, error) {
	query := `
		SELECT DISTINCT
			model as name
		FROM
			trade_in_trends
		WHERE
			brand = ?
		ORDER BY
			model ASC
		`

	var data = make([]entity.CommonResponse, 0)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.Brand)
	if err != nil {
		log.Error().Err(err).Str("brand", req.Brand).Msg("repo::GetHTIModels - Failed to get HTI models")
		return nil, err
	}

	return data, nil
}

func (r *commonRepo) GetHTITypes(ctx context.Context, req *entity.GetHTITypesRequest) ([]entity.CommonResponse, error) {
	query := `
		SELECT DISTINCT
			type as name
		FROM
			trade_in_trends
		WHERE
			brand = ? AND model = ?
		ORDER BY
			type ASC
		`

	var data = make([]entity.CommonResponse, 0)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.Brand, req.Model)
	if err != nil {
		log.Error().Err(err).Str("brand", req.Brand).Str("model", req.Model).Msg("repo::GetHTITypes - Failed to get HTI types")
		return nil, err
	}

	return data, nil
}

func (r *commonRepo) GetHTIYears(ctx context.Context, req *entity.GetHTIYearsRequest) ([]entity.CommonResponse, error) {
	query := `
		SELECT DISTINCT
			CAST(year AS VARCHAR) as name
		FROM
			trade_in_trends
		WHERE
			brand = ? AND model = ? AND type = ?
		ORDER BY
			CAST(year AS VARCHAR) ASC
		`

	var data = make([]entity.CommonResponse, 0)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.Brand, req.Model, req.Type)
	if err != nil {
		log.Error().Err(err).Str("brand", req.Brand).Str("model", req.Type).Str("type", req.Type).Msg("repo::GetHTIYears - Failed to get HTI years")
		return nil, err
	}

	return data, nil
}

func (r *commonRepo) GetHTIPurchase(ctx context.Context, req *entity.GetHTIPurchaseRequest) (entity.GetHTIPurchaseResponse, error) {
	query := `
		SELECT
			min_purchase,
			max_purchase
		FROM
			trade_in_trends
		WHERE
			brand = ? AND model = ? AND type = ? AND year = CAST(? AS INTEGER)
		`

	var data entity.GetHTIPurchaseResponse

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), req.Brand, req.Model, req.Type, req.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Str("brand", req.Brand).Str("model", req.Model).Str("type", req.Type).Str("year", req.Year).Msg("repo::GetHTIPurchase - Data not found")
			return data, errmsg.NewCustomErrors(404).SetMessage("Data tidak ditemukan")
		}
		log.Error().Err(err).Str("brand", req.Brand).Str("model", req.Model).Str("type", req.Type).Str("year", req.Year).Msg("repo::GetHTIPurchase - Failed to get HTI purchase")
		return data, err
	}

	return data, nil
}

func (r *commonRepo) GetHTIValuations(ctx context.Context, req *entity.GetHTIValuationsRequest) (entity.GetHTIValuationsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.HTIValuationItem
	}

	var (
		res  entity.GetHTIValuationsResponse
		data = make([]dao, 0)
	)
	res.Items = make([]entity.HTIValuationItem, 0, req.Paginate)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			type,
			year,
			min_purchase,
			max_purchase
		FROM
			trade_in_trends
		WHERE
			brand = ? AND model = ?
		ORDER BY
			year DESC, type ASC
		LIMIT ? OFFSET ?
		`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.Brand, req.Model, req.Paginate, (req.Page-1)*req.Paginate)
	if err != nil {
		log.Error().Err(err).Str("brand", req.Brand).Str("model", req.Model).Msg("repo::GetHTIValuations - Failed to get HTI valuations")
		return res, err
	}

	if len(data) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	for _, item := range data {
		res.Items = append(res.Items, item.HTIValuationItem)
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}
