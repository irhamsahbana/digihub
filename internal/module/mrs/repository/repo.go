package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/mrs/entity"
	"codebase-app/internal/module/mrs/ports"
	"codebase-app/pkg/errmsg"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var _ ports.MRSRepository = &mrsRepository{}

type mrsRepository struct {
	db *sqlx.DB
}

func NewMRSRepository() *mrsRepository {
	return &mrsRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *mrsRepository) GetMRSs(ctx context.Context, req *entity.GetMRSsRequest) (entity.GetMRSsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.MRSItem
	}

	var (
		res  entity.GetMRSsResponse
		data = make([]dao, 0)
	)
	res.Items = make([]entity.MRSItem, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			wac.id,
			c.name AS client,
			sa.name AS service_advisor,
			wac.follow_up_at
		FROM
			walk_around_checks wac
		LEFT JOIN
			clients c ON c.id = wac.client_id
		LEFT JOIN
			users sa ON sa.id = wac.user_id
		WHERE
			wac.branch_id = (SELECT branch_id FROM users WHERE id = ?)
			AND wac.is_needs_follow_up = TRUE
		ORDER BY
			wac.follow_up_at DESC
	`

	if err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.UserId); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetMRSs - Failed to get MRSs")
		return res, err
	}

	for _, d := range data {
		res.Items = append(res.Items, d.MRSItem)
	}

	if len(res.Items) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}

func (r *mrsRepository) RenewWAC(ctx context.Context, req *entity.RenewWACRequest) error {
	type daoWACC struct {
		Id           string `db:"id"`
		IsInterested bool   `db:"is_interested"`
	}

	var (
		newWACId         string = ulid.Make().String()
		waccs                   = make([]daoWACC, 0)
		IsWaccInterested        = make(map[string]bool)
		queryLog                = `
		INSERT INTO
			wac_follow_up_logs (id, walk_around_check_id, notes)
		VALUES
			(?, ?, ?)
	`
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

	query := `
		SELECT
			id,
			is_interested
		FROM
			walk_around_check_conditions
		WHERE
			walk_around_check_id = ?
	`

	err = tx.SelectContext(ctx, &waccs, tx.Rebind(query), req.WacId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to get WAC conditions")
		return err
	}

	var WaccNotInterested int

	for _, wacc := range waccs {
		IsWaccInterested[wacc.Id] = wacc.IsInterested
		if !wacc.IsInterested {
			WaccNotInterested++
		}
	}

	// validate if wacc in request is interested
	for _, waccId := range req.VehicleConditionIds {
		if _, ok := IsWaccInterested[waccId]; !ok {
			log.Warn().Any("payload", req).Msg("repository::RenewWAC - WAC condition not found")
			return errmsg.NewCustomErrors(404).SetMessage("Kondisi Kendaraan dengan id " + waccId + " tidak ditemukan")
		}

		if IsWaccInterested[waccId] {
			log.Warn().Any("payload", req).Msg("repository::RenewWAC - WAC condition already interested")
			return errmsg.NewCustomErrors(403).SetMessage("Kondisi Kendaraan dengan id " + waccId + " sudah tertarik")
		}
	}

	lengthConditionIds := len(req.VehicleConditionIds)
	WaccNotInterestedLeft := WaccNotInterested - lengthConditionIds
	isWACStillNeedFollowUp := WaccNotInterestedLeft > 0

	if lengthConditionIds > 0 {
		// create WAC copy from the previous WAC
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
			req.WacId,
			newWACId,
			req.WacId,
			lengthConditionIds,
			lengthConditionIds,
		)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create WAC copy")
			return err
		}

		query = `
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

		for i := 0; i < lengthConditionIds; i++ {
			if i == lengthConditionIds-1 {
				query += "?)"
			} else {
				query += "?, "
			}

			args = append(args, req.VehicleConditionIds[i])
		}

		query = tx.Rebind(query)
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to delete WAC conditions")
			return err
		}

		// update total follow ups in old WAC
		query = `
			UPDATE
				walk_around_checks
			SET
				total_potential_leads = total_potential_leads - ?,
				total_follow_ups = ?
			WHERE
				id = ?
		`

		_, err = tx.ExecContext(ctx, tx.Rebind(query), lengthConditionIds, WaccNotInterestedLeft, req.WacId)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to update total follow ups")
			return err
		}

		if isWACStillNeedFollowUp {
			query = `
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

			_, err = tx.ExecContext(ctx, tx.Rebind(query), req.WacId)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to update WAC")
				return err
			}

			// create log
			_, err = tx.ExecContext(ctx, tx.Rebind(queryLog),
				ulid.Make().String(),
				req.WacId,
				"perlu follow up lagi karena masih ada kondisi yang tidak tertarik",
			)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create log")
				return err
			}
		} else {
			query = `
				UPDATE
					walk_around_checks
				SET
					is_needs_follow_up = FALSE,
					follow_up_at = NULL,
					updated_at = NOW()
				WHERE
					id = ?
			`

			_, err = tx.ExecContext(ctx, tx.Rebind(query), req.WacId)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to update WAC")
				return err
			}

			// create log
			_, err = tx.ExecContext(ctx, tx.Rebind(queryLog),
				ulid.Make().String(),
				req.WacId,
				"semua kondisi berhasil menjadi leads",
			)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create log")
				return err
			}
		}
	} else { // if empty condition ids
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

		_, err = tx.ExecContext(ctx, tx.Rebind(query), req.WacId)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to update WAC")
			return err
		}

		// create log
		_, err = tx.ExecContext(ctx, tx.Rebind(queryLog),
			ulid.Make().String(),
			req.WacId,
			"follow up diperpanjang 7 hari ke depan",
		)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create log")
			return err
		}
		return nil
	}

	return nil
}
