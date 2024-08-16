package repository

import (
	"codebase-app/internal/module/wac/entity"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type user struct {
	Id        string `db:"id"`
	BranchId  string `db:"branch_id"`
	SectionId string `db:"section_id"`
}

func (r *wacRepository) CreateWAC(ctx context.Context, req *entity.CreateWACRequest) (entity.CreateWACResponse, error) {
	var result entity.CreateWACResponse

	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateWAC - Failed to begin transaction")
		return result, err
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				req.RemoveBase64()
				log.Error().Err(err).Any("payload", req).Msg("repo::CreateWAC - Failed to rollback transaction")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				req.RemoveBase64()
				log.Error().Err(err).Any("payload", req).Msg("repo::CreateWAC - Failed to commit transaction")
			}
		}
	}()

	// Generate WAC Id
	wacId := ulid.Make().String()

	// Get client Id or create new client if not exists
	clientId, err := r.getClientId(ctx, tx, req)
	if err != nil {
		req.RemoveBase64()
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateWAC - Failed to get client id")
		return result, err
	}

	// Get user data
	userData, err := r.getUserData(ctx, tx, req.UserId)
	if err != nil {
		req.RemoveBase64()
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateWAC - Failed to get user data")
		return result, err
	}

	// Create walk around check record
	err = r.createWACRecord(ctx, tx, wacId, userData, clientId, req.UserId)
	if err != nil {
		req.RemoveBase64()
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateWAC - Failed to create walk around check record")
		return result, err
	}

	// Create walk around check conditions
	err = r.createWACConditions(ctx, tx, wacId, req.VehicleConditions)
	if err != nil {
		req.RemoveBase64()
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateWAC - Failed to create walk around check conditions")
		return result, err
	}

	result.Id = wacId
	return result, nil
}

func (r *wacRepository) getClientId(ctx context.Context, tx *sqlx.Tx, req *entity.CreateWACRequest) (string, error) {
	var clientId string
	query := `SELECT id FROM clients WHERE vehicle_license_number = ?`
	err := tx.GetContext(ctx, &clientId, r.db.Rebind(query), req.VehicleRegistrationNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			clientId = ulid.Make().String()
			query = `
			INSERT INTO clients (id, name, vehicle_type_id, vehicle_license_number, phone)
			VALUES (?, ?, ?, ?, ?)`
			_, err = tx.ExecContext(ctx, r.db.Rebind(query), clientId, req.Name, req.VehicleTypeId, req.VehicleRegistrationNumber, req.WhatsAppNumber)
			if err != nil {
				req.RemoveBase64()
				log.Error().Err(err).Any("payload", req).Msg("repo::CreateWAC - Failed to create new client")
				return "", err
			}
		} else {
			req.RemoveBase64()
			log.Error().Err(err).Any("payload", req).Msg("repo::CreateWAC - Failed to get client id")
			return "", err
		}
	}
	return clientId, nil
}

func (r *wacRepository) getUserData(ctx context.Context, tx *sqlx.Tx, userId string) (user, error) {
	var u user
	query := `SELECT id, branch_id, section_id FROM users WHERE id = ?`
	err := tx.GetContext(ctx, &u, r.db.Rebind(query), userId)
	if err != nil {
		log.Error().Err(err).Any("user_id", userId).Msg("repo::CreateWAC - Failed to get user data")
		return u, err
	}
	return u, nil
}

func (r *wacRepository) createWACRecord(ctx context.Context, tx *sqlx.Tx, wacId string, u user, clientId, userId string) error {
	query := `
	INSERT INTO walk_around_checks (id, branch_id, section_id, user_id, client_id)
	VALUES (?, ?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, r.db.Rebind(query), wacId, u.BranchId, u.SectionId, userId, clientId)
	if err != nil {
		log.Error().Err(err).Any("wac_id", wacId).Any("user_id", userId).
			Msg("repo::CreateWAC - Failed to create walk around check")
	}
	return err
}

func (r *wacRepository) createWACConditions(ctx context.Context, tx *sqlx.Tx, wacId string, conditions []entity.VehicleCondition) error {
	for _, co := range conditions {
		ua := new(user)
		if co.ServiceAdvisorId != nil {
			query := `SELECT id, branch_id, section_id FROM users WHERE id = ?`
			err := tx.GetContext(ctx, ua, r.db.Rebind(query), *co.ServiceAdvisorId)
			if err != nil {
				log.Error().Err(err).Any("payload", co.ServiceAdvisorId).Msg("repo::CreateWAC - Failed to get user data")
				return err
			}
		}

		query := `
		INSERT INTO walk_around_check_conditions (
			id, walk_around_check_id, area_id, potency_id, notes, path,
			assigned_user_id, assigned_branch_id, assigned_section_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := tx.ExecContext(ctx, r.db.Rebind(query), ulid.Make().String(), wacId, co.AreaId, co.PotencyId, co.Notes, co.Path, ua.Id, ua.BranchId, ua.SectionId)
		if err != nil {
			log.Error().Err(err).Any("payload", co).Msg("repo::CreateWAC - Failed to create walk around check conditions")
			return err
		}
	}
	return nil
}
