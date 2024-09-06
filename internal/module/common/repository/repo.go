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

func (r *commonRepository) GetBranches(ctx context.Context, req *entity.GetBranchesRequest) (entity.GetBranchesResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.CommonResponse
	}

	var (
		result entity.GetBranchesResponse
		data   = make([]dao, 0, req.Paginate)
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id, name
		FROM
			branches
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.Paginate, (req.Page-1)*req.Paginate)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetBranches - Failed to get branches")
		return result, err
	}

	for _, d := range data {
		result.Items = append(result.Items, d.CommonResponse)
	}

	if len(data) > 0 {
		result.Meta.TotalData = data[0].TotalData
	}

	result.Meta.CountTotalPage(req.Page, req.Paginate, result.Meta.TotalData)

	return result, nil
}

func (r *commonRepository) GetAreas(ctx context.Context) ([]entity.AreaResponse, error) {
	var (
		result = make([]entity.AreaResponse, 0)
	)

	query := `
		SELECT
			id, name, type
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

func (r *commonRepository) GetPotencies(ctx context.Context, req *entity.GetPotenciesRequest) ([]entity.GetPotencyResponse, error) {
	type userDao struct {
		Id          string `db:"id"`
		Name        string `db:"name"`
		BranchId    string `db:"branch_id"`
		BranchName  string `db:"branch_name"`
		SectionName string `db:"section_name"`
	}
	var (
		result    = make([]entity.GetPotencyResponse, 0)
		potencies = make([]entity.CommonResponse, 0)
		user      = userDao{}
	)

	query := `
		SELECT
			id, name
		FROM
			potencies
	`

	err := r.db.SelectContext(ctx, &potencies, query)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetPotencies - Failed to get potencies")
		return nil, err
	}

	query = `
		SELECT
			usr.id, usr.name, usr.branch_id, br.name AS branch_name, sc.name AS section_name
		FROM
			users usr
		LEFT JOIN
			branches br ON usr.branch_id = br.id
		LEFT JOIN
			sections sc ON usr.section_id = sc.id
		WHERE
			usr.id = ?
	`

	err = r.db.GetContext(ctx, &user, r.db.Rebind(query), req.UserId)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetPotencies - Failed to get potencies")
		return nil, err
	}

	for _, p := range potencies {
		potency := entity.GetPotencyResponse{
			Id:   p.Id,
			Name: p.Name,
		}

		if user.SectionName == p.Name {
			u := entity.CommonResponse{
				Id:   user.Id,
				Name: user.Name,
			}

			b := entity.CommonResponse{
				Id:   user.BranchId,
				Name: user.BranchName,
			}

			potency.User = &u
			potency.Branch = &b
		}

		result = append(result, potency)
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

	if !req.IncludeMe {
		query.WriteString(" AND usr.id != ?")
		args = append(args, req.UserId)
	}

	if req.Role != "" {
		if req.Role == "service_advisor" {
			query.WriteString(" AND usr.role_id = (SELECT id FROM roles WHERE name = 'service_advisor')")
		} else if req.Role == "technician" {
			query.WriteString(" AND usr.role_id = (SELECT id FROM roles WHERE name = 'technician')")
		}
	} else {
		query.WriteString(" AND usr.role_id NOT IN (SELECT id FROM roles WHERE name = 'admin')")
	}

	if req.BranchId != "" {
		query.WriteString(" AND usr.branch_id = ?")
		args = append(args, req.BranchId)
	}

	query.WriteString(` LIMIT ? OFFSET ?`)
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
