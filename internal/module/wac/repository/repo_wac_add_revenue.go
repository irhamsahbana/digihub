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
	var (
		// this query is to check if all revenue added where invoice_number is null and is_interested is true
		queryCheckIsStillNeedRevenue = `
			SELECT EXISTS (
				SELECT
					1
				FROM
					walk_around_check_conditions
				WHERE
					walk_around_check_id = ?
					AND invoice_number IS NULL
					AND is_interested = TRUE
			)
		`
		IsStillNeedRevenue bool
	)

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

	err = tx.GetContext(ctx, &IsStillNeedRevenue, r.db.Rebind(queryCheckIsStillNeedRevenue), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to check all revenue added")
		return err
	}

	if !IsStillNeedRevenue { // if all revenue added, reject the request to add revenue
		log.Warn().Any("payload", req).Msg("repo::AddRevenues - all revenue already added")
		return errmsg.NewCustomErrors(400, errmsg.WithMessage("Semua revenue sudah diinput, tidak dapat mengubah revenue lagi"))
	}

	query := `
		UPDATE
			walk_around_check_conditions
		SET
			revenue = ?,
			invoice_number = ?,
			updated_at = NOW()
		WHERE
			id = ?
			AND walk_around_check_id = ?
			AND is_interested = TRUE
	`
	query = r.db.Rebind(query)

	for idx, revenue := range req.Revenues {
		_, err = tx.ExecContext(ctx, query,
			revenue.Revenue,
			revenue.InvoiceNumber,
			revenue.VehicleConditionId,
			req.Id,
		)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Any("item", req.Revenues[idx]).Msg("repo::AddRevenues - failed to add revenue")
			return err
		}
	}

	err = tx.GetContext(ctx, &IsStillNeedRevenue, r.db.Rebind(queryCheckIsStillNeedRevenue), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to check all revenue added")
		return err
	}

	if !IsStillNeedRevenue { // if all revenue added, update status to completed and follow up if needed
		query = `
			UPDATE
				walk_around_checks
			SET
				status = 'completed',
				updated_at = NOW()
			WHERE
				id = ?
		`

		_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.Id)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to update status")
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
			log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to check need follow up")
			return err
		}

		if isNeedFollowUp { // add 7 days from now in utc
			followUpAt := time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02 15:04:05")
			query = `
				UPDATE
					walk_around_checks
				SET
					is_needs_follow_up = TRUE,
					total_follow_ups = COALESCE((SELECT SUM(1) FROM walk_around_check_conditions WHERE walk_around_check_id = ? AND is_interested = FALSE), 0),
					updated_at = NOW(),
					follow_up_at = ?
				WHERE
					id = ?
			`

			_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.Id, followUpAt, req.Id)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to follow up")
				return err
			}
		}
	}

	return nil
}
