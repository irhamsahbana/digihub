package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/common/entity"
	"codebase-app/internal/module/common/ports"
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.CommonRepository = &commonRepository{}

type commonRepository struct {
	db *sqlx.DB
}

func NewCommonRepository() *commonRepository {
	return &commonRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *commonRepository) GetAreas(ctx context.Context) ([]entity.CommonResponse, error) {
	var (
		result = make([]entity.CommonResponse, 0)
	)

	query := `
		SELECT
			id, name
		FROM
			areas
	`

	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetAreas - Failed to get areas")
		return nil, err
	}

	return result, nil
}

func (r *commonRepository) GetPotencies(ctx context.Context) ([]entity.CommonResponse, error) {
	var (
		result = make([]entity.CommonResponse, 0)
	)

	query := `
		SELECT
			id, name
		FROM
			potencies
	`

	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetPotencies - Failed to get potencies")
		return nil, err
	}

	return result, nil
}

func (r *commonRepository) GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error) {
	var (
		result = make([]entity.CommonResponse, 0)
	)

	query := `
		SELECT
			id, name
		FROM
			vehicle_types
	`

	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetVehicleTypes - Failed to get vehicle types")
		return nil, err
	}

	return result, nil
}

func (r *commonRepository) GetEmployees(ctx context.Context, req *entity.GetEmployeesRequest) (entity.GetEmployeesResult, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.EmployeeResult
	}

	var (
		result = entity.GetEmployeesResult{}
		data   = make([]dao, 0, req.Paginate)
		query  = strings.Builder{}
		args   = make([]any, 0, 3)
	)

	query.WriteString(`
		SELECT
			COUNT(*) OVER() AS total_data,
			usr.id AS user_id,
			usr.branch_id,
			usr.section_id,
			usr.role_id,
			usr.name,
			usr.email,
			usr.whatsapp_number,
			br.name AS branch_name,
			sc.name AS section_name,
			rl.name AS role_name
		FROM
			users usr
		LEFT JOIN
			branches br ON usr.branch_id = br.id
		LEFT JOIN
			sections sc ON usr.section_id = sc.id
		LEFT JOIN
			roles rl ON usr.role_id = rl.id
		WHERE
			1 = 1
	`)

	if req.Role != "" {
		if req.Role == "service_advisor" {
			query.WriteString("AND usr.role_id = (SELECT id FROM roles WHERE name = 'service_advisor')")
		} else if req.Role == "technician" {
			query.WriteString("AND usr.role_id = (SELECT id FROM roles WHERE name = 'technician')")
		}
	} else {
		query.WriteString("AND usr.role_id NOT IN (SELECT id FROM roles WHERE name = 'admin')")
	}

	if req.BranchId != "" {
		query.WriteString(" AND usr.branch_id = ?")
		args = append(args, req.BranchId)
	}

	query.WriteString(`LIMIT ? OFFSET ?`)
	args = append(args, req.Paginate, (req.Page-1)*req.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query.String()), args...)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetEmployees - Failed to get employees")
		return result, err
	}

	if len(data) > 0 {
		result.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		result.Items = append(result.Items, d.EmployeeResult)
	}

	result.Meta.CountTotalPage(req.Page, req.Paginate, result.Meta.TotalData)

	return result, nil
}
