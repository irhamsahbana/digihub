package repository

import (
	"codebase-app/internal/module/employee/entity"
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (r *employeeRepo) ImportEmployees(ctx context.Context, rows []entity.ImportEmployeeRow) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", rows).Msg("repo::ImportEmployees - Failed to begin transaction")
		return err
	}
	defer tx.Rollback()

	queryNonAdmin := `
		INSERT INTO
			users (id, branch_id, section_id, role_id, name, email, password)
		VALUES
			(?, ?, ?, ?, TRIM(UPPER(?)), TRIM(LOWER(?)), ?)
		`

	queryAdmin := `
		INSERT INTO
			users (id, role_id, name, email, password)
		VALUES
			(?, ?, TRIM(UPPER(?)), TRIM(LOWER(?)), ?)
		`

	for _, row := range rows {
		if row.RoleName != "ADMIN" {
			_, err = tx.ExecContext(ctx, r.db.Rebind(queryNonAdmin),
				ulid.Make().String(),
				row.BranchId,
				row.SectionId,
				row.RoleId,
				row.Name,
				row.Email,
				row.PasswordHashed,
			)
			if err != nil {
				log.Error().Err(err).Any("payload", row).Msg("repo::ImportEmployees - Failed to import employee")
				return err
			}
		} else {
			_, err = tx.ExecContext(ctx, r.db.Rebind(queryAdmin),
				ulid.Make().String(),
				row.RoleId,
				row.Name,
				row.Email,
				row.PasswordHashed,
			)
			if err != nil {
				log.Error().Err(err).Any("payload", row).Msg("repo::ImportEmployees - Failed to import employee")
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", rows).Msg("repo::ImportEmployees - Failed to commit transaction")
		return err
	}

	return nil
}

func (r *employeeRepo) GetBranches(ctx context.Context) ([]entity.Common, error) {
	query := `
		SELECT
			id, name
		FROM
			branches
		WHERE deleted_at IS NULL
		`

	var (
		res = []entity.Common{}
	)

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(query))
	if err != nil {
		log.Error().Err(err).Msg("repo::GetBranches - Failed to get branches")
		return res, err
	}

	return res, nil
}

func (r *employeeRepo) GetPotencies(ctx context.Context) ([]entity.Common, error) {
	query := `
		SELECT
			id, name
		FROM
			potencies
		WHERE deleted_at IS NULL
		`

	var (
		res = []entity.Common{}
	)

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(query))
	if err != nil {
		log.Error().Err(err).Msg("repo::GetSections - Failed to get sections")
		return res, err
	}

	return res, nil
}

func (r *employeeRepo) GetRoles(ctx context.Context) ([]entity.Common, error) {
	query := `
		SELECT
			id, name
		FROM
			roles
		WHERE deleted_at IS NULL
		`

	var (
		res = []entity.Common{}
	)

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(query))
	if err != nil {
		log.Error().Err(err).Msg("repo::GetRoles - Failed to get roles")
		return res, err
	}

	return res, nil
}

func (r *employeeRepo) IsEmailExist(ctx context.Context, email string) error {
	query := `
		SELECT EXISTS(
			SELECT
				1
			FROM
				users
			WHERE
				email = ? AND deleted_at IS NULL
		)
		`

	var (
		exist bool
	)

	err := r.db.GetContext(ctx, &exist, r.db.Rebind(query), email)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("repo::IsEmailExist - Failed to check email")
		return err
	}

	if exist {
		return errors.New("email sudah terdaftar")
	}

	return nil
}
