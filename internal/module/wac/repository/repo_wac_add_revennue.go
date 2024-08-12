package repository

import (
	"codebase-app/internal/module/wac/entity"
	"context"

	"github.com/rs/zerolog/log"
)

func (r *wacRepository) AddWACRevenue(ctx context.Context, req *entity.AddWACRevenueRequest) (entity.AddWACRevenueResponse, error) {
	var res entity.AddWACRevenueResponse
	query := `
		UPDATE walk_around_checks
		SET invoice_number = ?, total_revenue = ?, status = 'completed'
		WHERE
			id = ?
			AND user_id = ?
			AND status = 'wip'
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.InvoiceNumber, req.TotalRevenue, req.Id, req.UserId)
	if err != nil {
		log.Error().Err(err).Msg("repository::AddWACRevenue - An error occurred")
		return res, err
	}

	res.Id = req.Id
	return res, nil
}
