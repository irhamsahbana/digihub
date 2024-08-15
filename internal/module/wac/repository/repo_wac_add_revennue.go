package repository

import (
	"codebase-app/internal/module/wac/entity"
	"context"

	"github.com/rs/zerolog/log"
)

func (r *wacRepository) AddRevenue(ctx context.Context, req *entity.AddWACRevenueRequest) error {
	query := `
		UPDATE
			walk_around_checks
		SET
			invoice_number = ?,
			revenue = ?,
			status = 'completed',
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.InvoiceNumber, req.TotalRevenue, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenue - failed to add revenue")
		return err
	}

	return nil
}
