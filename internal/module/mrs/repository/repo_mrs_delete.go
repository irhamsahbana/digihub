package repository

import (
	"codebase-app/internal/module/mrs/entity"
	"context"

	"github.com/rs/zerolog/log"
)

func (r *mrsRepository) DeleteFollowUp(ctx context.Context, req *entity.DeleteFollowUpRequest) error {
	query := `
		UPDATE
			walk_around_checks
		SET
			is_needs_follow_up = FALSE,
			total_follow_ups = 0,
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.WacId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteFollowUp - Failed to delete follow up")
		return err
	}

	return nil
}
