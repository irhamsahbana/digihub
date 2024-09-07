package repository

import (
	"codebase-app/internal/module/wac/entity"
	"context"

	"github.com/rs/zerolog/log"
)

func (r *wacRepository) OfferWAC(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error) {
	var res entity.OfferWACResponse

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

	for _, c := range req.VConditions {
		_, err = tx.ExecContext(ctx, r.db.Rebind(query), c.IsInterested, c.Notes, c.Id, req.Id)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::OfferWAC - Failed to update walk around check conditions")
			return res, err
		}
	}

	query = `
		UPDATE
			walk_around_checks
		SET
			status = 'wip',
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWAC - Failed to update walk around check record")
		return res, err
	}

	res.Id = req.Id

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
			status = 'offered',
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.Id)
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

	return res, nil
}
