package repository

import (
	"codebase-app/internal/module/wac/entity"
	"codebase-app/pkg/errmsg"
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (r *wacRepository) AddRevenue(ctx context.Context, req *entity.AddWACRevenueRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenue - failed to begin transaction")
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenue - failed to rollback transaction")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenue - failed to commit transaction")
			}
		}
	}()

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

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.InvoiceNumber, req.TotalRevenue, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenue - failed to add revenue")
		return err
	}

	var isNeedFollowUp bool
	query = `
		SELECT EXISTS (
			SELECT
				1
			FROM
				walk_around_check_conditions
			WHERE
				walk_around_check_id = ?
				AND is_interested = FALSE
		)
	`

	err = tx.GetContext(ctx, &isNeedFollowUp, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenue - failed to check need follow up")
		return err
	}

	if isNeedFollowUp { // add 7 days from now in utc
		followUpAt := time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02 15:04:05")
		query = `
			UPDATE
				walk_around_checks
			SET
				is_needs_follow_up = TRUE,
				updated_at = NOW(),
				follow_up_at = ?
			WHERE
				id = ?
		`

		_, err = tx.ExecContext(ctx, r.db.Rebind(query), followUpAt, req.Id)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenue - failed to follow up")
			return err
		}

		query = `
			INSERT INTO
				wac_follow_up_logs (id, walk_around_check_id, notes)
			VALUES (?, ?, ?)
		`

		_, err = tx.ExecContext(ctx, r.db.Rebind(query),
			ulid.Make().String(),
			req.Id,
			"perlu follow up",
		)
	}

	return nil
}

func (r *wacRepository) AddRevenues(ctx context.Context, req *entity.AddWACRevenuesRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to begin transaction")
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to rollback transaction")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to commit transaction")
			}
		}
	}()

	IsStillNeedRevenue, err := r.IsStillNeedAddRevenue(ctx, tx, req.Id)
	if err != nil {
		return err
	}

	if !IsStillNeedRevenue { // if all revenue added, reject the request to add revenue
		log.Warn().Any("payload", req).Msg("repo::AddRevenues - all revenue already added")
		return errmsg.NewCustomErrors(400, errmsg.WithMessage("Semua revenue sudah diinput, tidak dapat mengubah revenue lagi"))
	}

	err = r.updateConditions(ctx, tx, req)
	if err != nil {
		return err
	}

	IsStillNeedRevenue, err = r.IsStillNeedAddRevenue(ctx, tx, req.Id)
	if err != nil {
		return err
	}

	err = r.countTotalLeadsCompleted(ctx, tx, req.Id)
	if err != nil {
		return err
	}

	if !IsStillNeedRevenue { // if all revenue added, update status to completed and follow up if needed
		err = r.completingWAC(ctx, tx, req)
		if err != nil {
			return err
		}

		isNeedFollowUp, err := r.isNeedFollowUp(ctx, tx, req.Id)
		if err != nil {
			return err
		}

		if isNeedFollowUp { // add 7 days from now in utc
			err := r.updateNextFollowUpAt(ctx, tx, req)
			if err != nil {
				return err
			}

			err = r.createFollowUpLogs(ctx, tx, req)
			if err != nil {
				return err
			}
		}

		// create activity
		err = r.addCompletedActivity(ctx, tx, req)
		if err != nil {
			return err
		}
	}

	return nil
}
