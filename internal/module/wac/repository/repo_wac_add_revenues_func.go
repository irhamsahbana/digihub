package repository

import (
	"codebase-app/internal/module/wac/entity"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func (r *wacRepository) IsStillNeedAddRevenue(ctx context.Context, tx *sqlx.Tx, wacId string) (bool, error) {
	var IsStillNeedRevenue bool

	// this query is to check if all revenue added where invoice_number is null and is_interested is true
	query := `
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

	err := tx.GetContext(ctx, &IsStillNeedRevenue, r.db.Rebind(query), wacId)
	if err != nil {
		log.Error().Err(err).Any("wac_id", wacId).Msg("repo::IsStillNeedAddRevenue - failed to check all revenue added")
		return false, err
	}

	return IsStillNeedRevenue, nil
}

func (r *wacRepository) updateConditions(ctx context.Context, tx *sqlx.Tx, req *entity.AddWACRevenuesRequest) error {
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
		_, err := tx.ExecContext(ctx, query,
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

	return nil
}

func (r *wacRepository) countTotalLeadsCompleted(ctx context.Context, tx *sqlx.Tx, wacId string) error {
	query := `
		UPDATE
			walk_around_checks
		SET
			total_leads_completed = COALESCE(
				(
					SELECT
						SUM(1)
					FROM
						walk_around_check_conditions
					WHERE
						walk_around_check_id = ?
						AND is_interested = TRUE
						AND invoice_number IS NOT NULL
				),
				0
			),
			updated_at = NOW()
		WHERE
			id = ?
	`

	_, err := tx.ExecContext(ctx, r.db.Rebind(query), wacId, wacId)
	if err != nil {
		log.Error().Err(err).Any("wac_id", wacId).Msg("repo::countTotalLeadsCompleted - failed to count total leads completed")
		return err
	}

	return nil
}

func (r *wacRepository) completingWAC(ctx context.Context, tx *sqlx.Tx, req *entity.AddWACRevenuesRequest) error {
	query := `
			UPDATE
				walk_around_checks
			SET
				revenue = COALESCE(
					(
						SELECT
							SUM(revenue)
						FROM
							walk_around_check_conditions
						WHERE
							walk_around_check_id = ?
							AND is_interested = TRUE
							AND invoice_number IS NOT NULL
					),
					0
				),
				status = 'completed',
				updated_at = NOW()
			WHERE
				id = ?
		`

	_, err := tx.ExecContext(ctx, r.db.Rebind(query), req.Id, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::completingWAC - failed to update status")
		return err
	}

	return nil
}

func (r *wacRepository) isNeedFollowUp(ctx context.Context, tx *sqlx.Tx, wacId string) (bool, error) {
	var isNeedFollowUp bool
	query := `
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

	err := tx.GetContext(ctx, &isNeedFollowUp, tx.Rebind(query), wacId)
	if err != nil {
		log.Error().Err(err).Any("wac_id", wacId).Msg("repo::isNeedFollowUp - failed to check need follow up")
		return false, err
	}

	return isNeedFollowUp, nil
}

func (r *wacRepository) updateNextFollowUpAt(ctx context.Context, tx *sqlx.Tx, req *entity.AddWACRevenuesRequest) error {
	// add 7 days from
	followUpAt := time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02 15:04:05")
	query := `
				UPDATE
					walk_around_checks
				SET
					is_needs_follow_up = TRUE,
					total_follow_ups = COALESCE(
						(
							SELECT
								SUM(1)
							FROM
								walk_around_check_conditions
							WHERE
								walk_around_check_id = ?
								AND is_interested = FALSE
						),
						0
					),
					updated_at = NOW(),
					follow_up_at = ?
				WHERE
					id = ?
			`

	_, err := tx.ExecContext(ctx, tx.Rebind(query), req.Id, followUpAt, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to follow up")
		return err
	}

	return nil
}

func (r *wacRepository) createFollowUpLogs(ctx context.Context, tx *sqlx.Tx, req *entity.AddWACRevenuesRequest) error {
	query := `
		INSERT INTO
			wac_follow_up_logs (id, walk_around_check_id, notes)
		VALUES (?, ?, ?)
	`

	_, err := tx.ExecContext(ctx, tx.Rebind(query),
		ulid.Make().String(),
		req.Id,
		"perlu follow up",
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to create log")
		return err
	}

	return nil
}

func (r *wacRepository) addCompletedActivity(ctx context.Context, tx *sqlx.Tx, req *entity.AddWACRevenuesRequest) error {
	query := `
	SELECT
		id AS wac_id,
		user_id,
		total_potential_leads,
		total_leads,
		(
			SELECT
				COALESCE(SUM(1), 0)
			FROM
				walk_around_check_conditions
			WHERE
				walk_around_check_id = ?
				AND is_interested = TRUE
				AND invoice_number IS NOT NULL
		) AS total_completed_leads
	FROM
		walk_around_checks
	WHERE
		id = ?
	`

	var a activity
	a.Status = "completed"
	err := tx.GetContext(ctx, &a, r.db.Rebind(query), req.Id, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::AddRevenues - failed to get walk around check record")
		return err
	}

	totalRevenue := decimal.NewFromFloat(0.00)
	for _, revenue := range req.Revenues {
		totalRevenue = totalRevenue.Add(decimal.NewFromFloat(revenue.Revenue))
	}
	a.TotalRevenue, _ = totalRevenue.Float64()

	err = r.addActivity(ctx, tx, a)
	if err != nil {
		return err
	}
	return nil
}
