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

// IsWACStatus checks if a Walk Around Check (WAC) with the given ID has the specified status.
// It returns true if the status matches, otherwise false. If an error occurs during the query,
// it logs the error and returns false along with the error.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - WACId: The ID of the Walk Around Check to be checked.
//   - status: The status to be checked against the WAC.
//
// Returns:
//   - bool: True if the WAC has the specified status, otherwise false.
//   - error: An error if one occurs during the database query.
func (r *wacRepository) IsWACStatus(ctx context.Context, WACId, status string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM walk_around_checks
			WHERE id = ? AND status = ?
		)
	`

	var isStatus bool
	err := r.db.GetContext(ctx, &isStatus, r.db.Rebind(query), WACId, status)
	if err != nil {
		log.Error().Err(err).Msg("repository::IsWACStatus - An error occurred")
		return false, err
	}

	return isStatus, nil
}
