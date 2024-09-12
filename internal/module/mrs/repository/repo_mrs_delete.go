package repository

import (
	"codebase-app/internal/module/mrs/entity"
	"context"
	"strconv"

	"github.com/rs/zerolog/log"
)

func (r *mrsRepository) DeleteFollowUp(ctx context.Context, req *entity.DeleteFollowUpRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteFollowUp - Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::DeleteFollowUp - Failed to rollback transaction")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::DeleteFollowUp - Failed to commit transaction")
			}
		}
	}()

	queryCountNotInterestedLeft := `
		SELECT
			COUNT(*)
		FROM
			walk_around_check_conditions
		WHERE
			walk_around_check_id = ?
			AND is_interested = FALSE
	`
	var countNotInterestedLeft int
	err = tx.GetContext(ctx, &countNotInterestedLeft, r.db.Rebind(queryCountNotInterestedLeft), req.WacId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteFollowUp - Failed to get count not interested left")
		return err
	}

	query := `
		UPDATE
			walk_around_checks
		SET
			is_needs_follow_up = FALSE,
			is_followed_up = TRUE,
			total_follow_ups = 0,
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.WacId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteFollowUp - Failed to delete follow up")
		return err
	}

	left := strconv.Itoa(countNotInterestedLeft)
	err = r.createFollowUpLog(ctx, tx,
		req.WacId,
		"dihapus dari daftar penawaran dengan sisa penawaran yang tidak tertarik sebanyak "+left+" item",
	)

	return nil
}
