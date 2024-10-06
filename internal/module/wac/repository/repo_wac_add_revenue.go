package repository

import (
	"codebase-app/internal/module/wac/entity"
	"codebase-app/pkg/errmsg"
	"context"

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

	isUsedCar := false
	query := `
		SELECT
			is_used_car
		FROM
			walk_around_checks
		WHERE
			id = ?
	`

	err = tx.GetContext(ctx, &isUsedCar, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenue - failed to get is_used_car")
		return err
	}

	if !isUsedCar {
		log.Warn().Any("payload", req).Msg("repo::AddRevenue - not used car")
		return errmsg.NewCustomErrors(400, errmsg.WithMessage("Bukan mobil bekas"))
	}

	query = `
		UPDATE
			walk_around_checks
		SET
			invoice_number = ?,
			revenue = ?,
			total_leads_completed = (
				SELECT
					COALESCE(SUM(1), 0)
				FROM
					walk_around_check_conditions
				WHERE
					walk_around_check_id = ?
			),
			status = 'completed',
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.InvoiceNumber, req.TotalRevenue, req.Id, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenue - failed to add revenue")
		return err
	}

	query = `
			SELECT
				id AS wac_id,
				user_id,
				total_potential_leads,
				total_leads,
				total_leads_completed AS total_completed_leads
				total_revenue
			FROM
				walk_around_checks
			WHERE
				id = ?
			`

	var a activity
	a.Status = "completed"
	err = tx.GetContext(ctx, &a, r.db.Rebind(query), req.Id, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to get walk around check record")
		return err
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
