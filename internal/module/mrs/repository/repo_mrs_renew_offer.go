package repository

import (
	"codebase-app/internal/module/mrs/entity"
	"codebase-app/pkg/errmsg"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (r *mrsRepository) RenewWAC(ctx context.Context, req *entity.RenewWACRequest) error {
	var (
		newWACId string = ulid.Make().String()
		waccs           = make([]daoWACC, 0)
	)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to begin transaction")
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to rollback transaction")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to commit transaction")
			}
		}
	}()

	// get WAC conditions
	err = r.getVehicleConditions(ctx, req, tx, &waccs)
	if err != nil {
		return err
	}

	isWACStillNeedFollowUp, WaccNotInterestedLeft, err := r.validateVehicleConditions(req, &waccs)
	if err != nil {
		return err
	}

	totalNewLeads := len(req.VehicleConditionIds)
	if totalNewLeads > 0 {
		err = r.createWACCopy(ctx, req, tx, req.WacId, newWACId, totalNewLeads)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create WAC copy")
			return err
		}

		err = r.moveVehicleConditionsToNewWAC(ctx, req, tx, newWACId, totalNewLeads)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to move WAC conditions")
			return err
		}

		err = r.updateTotalFollowUps(ctx, req.WacId, tx, totalNewLeads, WaccNotInterestedLeft)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to update total follow ups")
			return err
		}

		if isWACStillNeedFollowUp {
			err = r.ExtendFollowUpAndRecountingTotalFollowUps(ctx, tx, req.WacId)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to extend follow up and recounting total follow ups")
				return err
			}

			err = r.createFollowUpLog(ctx, tx, req.WacId, "perlu follow up lagi karena masih ada kondisi yang tidak tertarik")
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create follow up log")
				return err
			}
		} else {
			err := r.removeWACFromFollowUpList(ctx, tx, req.WacId)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to remove WAC from follow up list")
				return err
			}

			err = r.createFollowUpLog(ctx, tx, req.WacId, "semua kondisi berhasil menjadi leads")
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create log")
				return err
			}
		}
	} else { // if empty condition ids
		err = r.updateFollowUpDeadline(ctx, tx, req.WacId)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to update follow up deadline")
			return err
		}

		err = r.createFollowUpLog(ctx, tx, req.WacId, "follow up diperpanjang 7 hari ke depan")
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create log")
			return err
		}
		return nil
	}

	return nil
}

type daoWACC struct {
	Id           string `db:"id"`
	IsInterested bool   `db:"is_interested"`
}

func (r *mrsRepository) getVehicleConditions(ctx context.Context, req *entity.RenewWACRequest, tx *sqlx.Tx, waccs *[]daoWACC) error {
	query := `
		SELECT
			id,
			is_interested
		FROM
			walk_around_check_conditions
		WHERE
			walk_around_check_id = ?
	`

	err := tx.SelectContext(ctx, waccs, tx.Rebind(query), req.WacId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to get WAC conditions")
		return err
	}

	return nil
}

func (r *mrsRepository) validateVehicleConditions(req *entity.RenewWACRequest, waccs *[]daoWACC) (isWACStillNeedFollowUp bool, WaccNotInterestedLeft int, err error) {
	var IsWaccInterested = make(map[string]bool)
	var WaccNotInterested int

	for _, wacc := range *waccs {
		IsWaccInterested[wacc.Id] = wacc.IsInterested
		if !wacc.IsInterested {
			WaccNotInterested++
		}
	}

	lengthConditionIds := len(req.VehicleConditionIds)
	WaccNotInterestedLeft = WaccNotInterested - lengthConditionIds
	isWACStillNeedFollowUp = WaccNotInterestedLeft > 0

	if lengthConditionIds > 0 {
		// validate if wacc in request is interested
		for _, waccId := range req.VehicleConditionIds {
			if _, ok := IsWaccInterested[waccId]; !ok {
				log.Warn().Any("payload", req).Msg("repository::RenewWAC - WAC condition not found")
				return isWACStillNeedFollowUp, WaccNotInterestedLeft, errmsg.NewCustomErrors(404).SetMessage("Kondisi Kendaraan dengan id " + waccId + " tidak ditemukan")
			}

			if IsWaccInterested[waccId] {
				log.Warn().Any("payload", req).Msg("repository::RenewWAC - WAC condition already interested")
				return isWACStillNeedFollowUp, WaccNotInterestedLeft, errmsg.NewCustomErrors(403).SetMessage("Kondisi Kendaraan dengan id " + waccId + " sudah tertarik")
			}
		}
	}

	return isWACStillNeedFollowUp, WaccNotInterestedLeft, nil
}

func (r *mrsRepository) createWACCopy(ctx context.Context, req *entity.RenewWACRequest, tx *sqlx.Tx, oldWACId, newWACId string, totalLeads int) error {
	queryCopy := `
			WITH wac_copy AS (
				SELECT
					*
				FROM
					walk_around_checks
				WHERE
					id = ?
			)
			INSERT INTO walk_around_checks (
				id,
				follow_up_wac_id,
				branch_id,
				section_id,
				client_id,
				user_id,
				status,

				total_potential_leads,
				total_leads
			) VALUES (
				?,
				?,
				(SELECT branch_id FROM wac_copy),
				(SELECT section_id FROM wac_copy),
				(SELECT client_id FROM wac_copy),
				(SELECT user_id FROM wac_copy),
				'wip',
				?,
				?
			)
		`

	_, err := tx.ExecContext(ctx, tx.Rebind(queryCopy),
		oldWACId,
		newWACId,
		req.WacId,
		totalLeads,
		totalLeads,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *mrsRepository) moveVehicleConditionsToNewWAC(
	ctx context.Context,
	req *entity.RenewWACRequest,
	tx *sqlx.Tx,
	newWACId string,
	totalNewLeads int,
) error {
	query := `
			UPDATE
				walk_around_check_conditions
			SET
				walk_around_check_id = ?,
				is_interested = TRUE
			WHERE
				is_interested = FALSE
				AND walk_around_check_id = ?
				AND id IN (
		`
	args := make([]interface{}, 0)
	args = append(args, newWACId, req.WacId)

	for i := 0; i < totalNewLeads; i++ {
		if i == totalNewLeads-1 {
			query += "?)"
		} else {
			query += "?, "
		}

		args = append(args, req.VehicleConditionIds[i])
	}

	query = tx.Rebind(query)
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *mrsRepository) updateTotalFollowUps(
	ctx context.Context,
	oldWACId string,
	tx *sqlx.Tx,
	totalNewLeads,
	WaccNotInterestedLeft int,
) error {
	query := `
		UPDATE
			walk_around_checks
		SET
			total_potential_leads = total_potential_leads - ?,
			total_follow_ups = ?
		WHERE
			id = ?
	`

	_, err := tx.ExecContext(ctx, tx.Rebind(query), totalNewLeads, WaccNotInterestedLeft, oldWACId)
	if err != nil {
		return err
	}

	return nil
}

func (r *mrsRepository) ExtendFollowUpAndRecountingTotalFollowUps(
	ctx context.Context,
	tx *sqlx.Tx,
	oldWACId string,
) error {
	query := `
	UPDATE
		walk_around_checks
	SET
		is_needs_follow_up = TRUE,
		follow_up_at = NOW() + INTERVAL '7 day',
		updated_at = NOW(),
		total_follow_ups = (
			SELECT
				COUNT(wacc.id)
			FROM
				walk_around_check_conditions wacc
			WHERE
				wacc.walk_around_check_id = ?
				AND wacc.is_interested = FALSE
		)
	WHERE
		id = ?
`

	_, err := tx.ExecContext(ctx, tx.Rebind(query), oldWACId, oldWACId)
	if err != nil {
		return err
	}

	return nil
}

func (r *mrsRepository) updateFollowUpDeadline(
	ctx context.Context,
	tx *sqlx.Tx,
	WACId string,
) error {
	query := `
			UPDATE
				walk_around_checks
			SET
				is_needs_follow_up = TRUE,
				follow_up_at = NOW() + INTERVAL '7 day',
				updated_at = NOW()
			WHERE
				id = ?
		`

	_, err := tx.ExecContext(ctx, tx.Rebind(query), WACId)
	if err != nil {
		return err
	}
	return nil
}

func (r *mrsRepository) removeWACFromFollowUpList(
	ctx context.Context,
	tx *sqlx.Tx,
	WACId string,
) error {
	query := `
				UPDATE
					walk_around_checks
				SET
					is_needs_follow_up = FALSE,
					follow_up_at = NULL,
					updated_at = NOW()
				WHERE
					id = ?
			`

	_, err := tx.ExecContext(ctx, tx.Rebind(query), WACId)
	if err != nil {
		return err
	}

	return nil
}

func (r *mrsRepository) createFollowUpLog(
	ctx context.Context,
	tx *sqlx.Tx,
	WACId string,
	notes string,
) error {
	query := `
		INSERT INTO
			wac_follow_up_logs (id, walk_around_check_id, notes)
		VALUES
			(?, ?, ?)
	`

	_, err := tx.ExecContext(ctx, tx.Rebind(query),
		ulid.Make().String(),
		WACId,
		notes,
	)
	if err != nil {
		return err
	}

	return nil
}