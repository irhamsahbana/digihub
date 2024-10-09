package repository

import (
	"codebase-app/internal/module/mrs/entity"
	"codebase-app/pkg/errmsg"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (r *mrsRepository) getWACOwnerId(ctx context.Context, tx *sqlx.Tx, wacId string) (string, error) {
	var userId string

	query := `
		SELECT user_id
		FROM Walk_around_checks
		WHERE id = ?
	`

	err := tx.GetContext(ctx, &userId, tx.Rebind(query), wacId)
	if err != nil {
		log.Error().Err(err).Str("wac_id", wacId).Msg("repo::getWACOwnerId - An error occurred")
		return "", err
	}

	return userId, nil
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

func (r *mrsRepository) createWACCopy(ctx context.Context, tx *sqlx.Tx, req *entity.RenewWACRequest, oldWACId, newWACId string, totalLeads int) error {
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
	tx *sqlx.Tx,
	req *entity.RenewWACRequest,
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
	tx *sqlx.Tx,
	oldWACId string,
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
				COUNT(*)
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
				follow_up_at = follow_up_at + INTERVAL '7 day',
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
			is_followed_up = TRUE,
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
