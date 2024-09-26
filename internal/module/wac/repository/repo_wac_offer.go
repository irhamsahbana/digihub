package repository

import (
	"codebase-app/internal/module/wac/entity"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r *wacRepository) OfferWAC(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error) {
	var (
		res           entity.OfferWACResponse
		isAnyInterest bool
	)

	if req.IsUsedCar {
		return r.OfferWACUsedCard(ctx, req)
	}

	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWAC - Failed to begin transaction")
		return res, err
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repo::OfferWAC - Failed to rollback transaction")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repo::OfferWAC - Failed to commit transaction")
			}
		}
	}()

	for _, c := range req.VConditions {
		if c.IsInterested {
			isAnyInterest = true
			break
		}
	}

	if !isAnyInterest {
		err = r.toCompletedWAC(ctx, tx, req)
		if err != nil {
			return res, err
		}

		return res, nil
	}

	query := `
		UPDATE
			walk_around_check_conditions
		SET
			is_interested = ?,
			notes = ?,
			updated_at = NOW()
		WHERE
			id = ?
			AND walk_around_check_id = ?
	`
	var totalLeads int
	for _, c := range req.VConditions {
		_, err = tx.ExecContext(ctx, r.db.Rebind(query), c.IsInterested, c.Notes, c.Id, req.Id)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::OfferWAC - Failed to update walk around check conditions")
			return res, err
		}

		if c.IsInterested {
			totalLeads = totalLeads + 1
		}
	}

	query = `
		UPDATE
			walk_around_checks
		SET
			status = 'wip',
			total_leads = ?,
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), totalLeads, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWAC - Failed to update walk around check record")
		return res, err
	}

	res.Id = req.Id

	// create activity
	query = `
		SELECT
			user_id,
			total_potential_leads,
			total_leads
		FROM
			walk_around_checks
		WHERE
			id = ?
	`

	var a activity
	a.WacId = req.Id
	a.Status = "wip"
	err = tx.GetContext(ctx, &a, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWAC - Failed to get walk around check record")
		return res, err
	}

	err = r.addActivity(ctx, tx, a)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *wacRepository) OfferWACUsedCard(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error) {
	var res entity.OfferWACResponse

	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWACUsedCard - Failed to begin transaction")
		return res, err
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repo::OfferWACUsedCard - Failed to rollback transaction")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repo::OfferWACUsedCard - Failed to commit transaction")
			}
		}
	}()

	query := `
		UPDATE
			walk_around_checks
		SET
			is_used_car = TRUE,
			total_leads = COALESCE((SELECT SUM(1) FROM walk_around_check_conditions WHERE walk_around_check_id = ?), 0),
			status = 'wip',
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.Id, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWACUsedCard - Failed to update walk around check record")
		return res, err
	}

	query = `
		UPDATE
			walk_around_check_conditions
		SET
			is_interested = ?,
			updated_at = NOW()
		WHERE
			walk_around_check_id = ?
	`

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), true, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWACUsedCard - Failed to update walk around check conditions")
		return res, err
	}

	res.Id = req.Id

	// create activity
	query = `
		SELECT
			user_id,
			total_potential_leads,
			total_leads
		FROM
			walk_around_checks
		WHERE
			id = ?
	`

	var a activity
	a.WacId = req.Id
	a.Status = "wip"
	err = tx.GetContext(ctx, &a, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWACUsedCard - Failed to get walk around check record")
		return res, err
	}

	err = r.addActivity(ctx, tx, a)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *wacRepository) toCompletedWAC(ctx context.Context, tx *sqlx.Tx, req *entity.OfferWACRequest) error {
	followUpAt := time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02 15:04:05")
	query := `
		UPDATE
			walk_around_checks
		SET
			status = 'completed',
			total_follow_ups = total_potential_leads,
			is_needs_follow_up = TRUE,
			follow_up_at = ?,
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err := tx.ExecContext(ctx, r.db.Rebind(query), followUpAt, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::toCompletedWAC - Failed to update walk around check record")
		return err
	}

	return nil
}
