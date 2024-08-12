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
			is_interested = :is_interested,
			notes = :notes,
			updated_at = NOW()
		WHERE
			id = :id
			AND walk_around_check_id = :walk_around_check_id
	`

	dataMap := make([]map[string]any, 0, len(req.VConditions))

	for _, v := range req.VConditions {
		dataMap = append(dataMap, map[string]any{
			"id":                   v.Id,
			"walk_around_check_id": req.Id,
			"is_interested":        v.IsInterested,
			"notes":                v.Notes,
		})
	}

	_, err = tx.NamedExecContext(ctx, query, dataMap)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWAC - Failed to update walk around check conditions")
		return res, err
	}

	query = `
		UPDATE
			walk_around_checks
		SET
			status = 'offered',
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

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), false, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::OfferWACUsedCard - Failed to update walk around check conditions")
		return res, err
	}

	res.Id = req.Id

	return res, nil
}
