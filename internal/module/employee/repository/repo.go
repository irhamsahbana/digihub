package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/employee/entity"
	"codebase-app/internal/module/employee/ports"
	"codebase-app/pkg"
	"codebase-app/pkg/errmsg"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var _ ports.EmployeeRepository = &employeeRepository{}

type employeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository() *employeeRepository {
	return &employeeRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *employeeRepository) GetEmployee(ctx context.Context, req *entity.GetEmployeeRequest) (entity.GetEmployeeResponse, error) {
	query := `
		SELECT
			u.id, u.name, u.email, u.whatsapp_number,
			b.id AS branch_id, b.name AS branch_name,
			p.id AS section_id, p.name AS section_name,
			r.id AS role_id, r.name AS role_name
		FROM
			users u
		JOIN
			roles r ON r.id = u.role_id
		JOIN
			branches b ON b.id = u.branch_id
		JOIN
			potencies p ON p.id = u.section_id
		WHERE
			u.id = ? AND u.deleted_at IS NULL
		`

	var (
		res = entity.GetEmployeeResponse{}
	)

	err := r.db.GetContext(ctx, &res, r.db.Rebind(query), req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Any("payload", req).Msg("repo::GetEmployee - Employee not found")
			return res, errmsg.NewCustomErrors(404).SetMessage("Karyawan tidak ditemukan")
		}
		log.Error().Err(err).Any("payload", req).Msg("repo::GetEmployee - Failed to get employee")
		return res, err
	}

	res.Branch = entity.Common{
		Id:   res.BranchId,
		Name: res.BranchName,
	}

	res.Section = entity.Common{
		Id:   res.SectionId,
		Name: res.SectionName,
	}

	res.Role = entity.Common{
		Id:   res.RoleId,
		Name: res.RoleName,
	}

	return res, nil
}

func (r *employeeRepository) UpdateEmployee(ctx context.Context, req *entity.UpdateEmployeeRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::UpdateEmployee - Failed to begin transaction")
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE
			users
		SET
			branch_id = ?, section_id = ?, role_id = ?, whatsapp_number = ?, name = ?, email = ?
		WHERE
			id = ?
		`

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.BranchId, req.SectionId, req.RoleId, req.WANumber, req.Name, req.Email, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::UpdateEmployee - Failed to update employee")
		return err
	}

	if req.Password != "" {
		query = `
			UPDATE
				users
			SET
				password = ?
			WHERE
				id = ?
			`

		hashed, err := pkg.HashPassword(req.Password)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::UpdateEmployee - Failed to hash password")
			return err
		}

		_, err = tx.ExecContext(ctx, r.db.Rebind(query), hashed, req.Id)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::UpdateEmployee - Failed to update password")
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::UpdateEmployee - Failed to commit transaction")
		return err
	}

	return nil
}

func (r *employeeRepository) CreateEmployee(ctx context.Context, req *entity.CreateEmployeeRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateEmployee - Failed to begin transaction")
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO
			users (id, branch_id, section_id, role_id, name, email, whatsapp_number, password)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?)
		`

	hashed, err := pkg.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateEmployee - Failed to hash password")
		return err
	}

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), ulid.Make().String(), req.BranchId, req.SectionId, req.RoleId, req.Name, req.Email, req.WANumber, hashed)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateEmployee - Failed to create employee")
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateEmployee - Failed to commit transaction")
		return err
	}

	return nil
}

func (r *employeeRepository) DeleteEmployee(ctx context.Context, req *entity.DeleteEmployeeRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::DeleteEmployee - Failed to begin transaction")
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE
			users
		SET
			deleted_at = NOW()
		WHERE
			id = ? AND deleted_at IS NULL
		`

	_, err = tx.ExecContext(ctx, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::DeleteEmployee - Failed to delete employee")
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::DeleteEmployee - Failed to commit transaction")
		return err
	}

	return nil
}
