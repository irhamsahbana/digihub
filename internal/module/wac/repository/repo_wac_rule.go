package repository

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (r *wacRepository) IsWACCreator(ctx context.Context, userId, WACId string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM walk_around_checks
			WHERE id = ? AND user_id = ?
		)
	`

	var isCreator bool
	err := r.db.GetContext(ctx, &isCreator, r.db.Rebind(query), WACId, userId)
	if err != nil {
		log.Error().Err(err).Msg("repository::IsWACCreator - An error occurred")
		return false, err
	}

	return isCreator, nil
}

func (r *wacRepository) IsWACOffered(ctx context.Context, WACId string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM walk_around_checks
			WHERE id = ? AND status != 'created'
		)
	`

	var isOffered bool
	err := r.db.GetContext(ctx, &isOffered, r.db.Rebind(query), WACId)
	if err != nil {
		log.Error().Err(err).Msg("repository::IsWACOffered - An error occurred")
		return false, err
	}

	return isOffered, nil
}
