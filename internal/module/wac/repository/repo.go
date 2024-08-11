package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/wac/entity"
	"codebase-app/internal/module/wac/ports"
	"codebase-app/pkg/errmsg"
	storage "codebase-app/pkg/storage-manager"
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.WACRepository = &wacRepository{}

type wacRepository struct {
	db *sqlx.DB
}

func NewWACRepository() *wacRepository {
	return &wacRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *wacRepository) GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.WacItem
	}

	var (
		query = strings.Builder{}
		args  = make([]any, 0, 8)
		res   entity.GetWACsResponse
		data  = make([]dao, 0, req.Paginate)
	)
	res.Items = make(map[string][]entity.WacItem)

	query.WriteString(`
		SELECT
			COUNT(*) OVER() AS total_data,
			wac.id,
			c.name AS client_name,
			wac.status,
			wac.created_at
		FROM
			walk_around_checks wac
		LEFT JOIN
			walk_around_check_conditions wacc ON wacc.walk_around_check_id = wac.id
		LEFT JOIN
			clients c ON c.id = wac.client_id
		WHERE
			wac.deleted_at IS NULL
			AND (wac.user_id = ? OR wacc.assigned_user_id = ?)
	`)
	args = append(args, req.UserId, req.UserId)

	if req.Query != "" {
		query.WriteString(" AND (c.name ILIKE ? OR c.vehicle_license_number ILIKE ?)")
		args = append(args, "%"+req.Query+"%", "%"+req.Query+"%")
	}

	if req.Status != "" {
		query.WriteString(" AND wac.status = ?")
		args = append(args, req.Status)
	}

	query.WriteString(" ORDER BY wac.created_at DESC")

	query.WriteString(" LIMIT ? OFFSET ?")
	args = append(args, req.Paginate, (req.Page-1)*req.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query.String()), args...)
	if err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("repo::GetWACs - failed to get wacs")
		return res, err
	}

	// today := time.Now().In(time.FixedZone("WITA", 28800)).Format("2006-01-02") // WITA
	today := time.Now().Format("2006-01-02")
	length := len(data)
	for length > 0 {
		length--
		d := data[length]

		// date := d.CreatedAt.In(time.FixedZone("WITA", 28800)).Format("2006-01-02")
		date := d.CreatedAt.Format("2006-01-02")
		if today == date {
			date = "Hari ini"
		}

		res.Items[date] = append(res.Items[date], d.WacItem)
	}

	if len(data) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}

func (r *wacRepository) GetWAC(ctx context.Context, req *entity.GetWACRequest) (entity.GetWACResponse, error) {
	type dao struct {
		Id          string  `db:"id"`
		ClientName  string  `db:"client_name"`
		VLicenseNum string  `db:"vehicle_license_number"`
		VTypeId     string  `db:"vehicle_type_id"`
		VTypeName   string  `db:"vehicle_type_name"`
		ClientWANum string  `db:"whatsapp_number"`
		IsOffered   bool    `db:"is_offered"`
		InvoiceNum  *string `db:"invoice_number"`
		Revenue     float64 `db:"revenue"`
		Status      string  `db:"status"`
	}

	type daoVCondition struct {
		Id              string  `db:"id"`
		PotencyId       string  `db:"potency_id"`
		PotencyName     string  `db:"potency_name"`
		AreaId          string  `db:"area_id"`
		AreaName        string  `db:"area_name"`
		AreaType        string  `db:"area_type"`
		AUserId         string  `db:"assigned_user_id"`
		AUserName       string  `db:"assigned_user_name"`
		AUserBranchId   string  `db:"assigned_user_branch_id"`
		AUserBranchName string  `db:"assigned_user_branch_name"`
		Path            string  `db:"path"`
		IsInterested    bool    `db:"is_interested"`
		Notes           *string `db:"notes"`
	}

	var (
		res    entity.GetWACResponse
		data   dao
		datavc = make([]daoVCondition, 0)
		query  = strings.Builder{}
	)

	query.WriteString(`
		SELECT
			wac.id,
			c.name AS client_name,
			c.vehicle_license_number,
			vt.id AS vehicle_type_id,
			vt.name AS vehicle_type_name,
			c.phone AS whatsapp_number,
			CASE
				WHEN wac.status = 'offered' THEN TRUE
				ELSE FALSE
			END AS is_offered,
			wac.invoice_number,
			wac.revenue,
			wac.status
		FROM
			walk_around_checks wac
		LEFT JOIN
			clients c ON c.id = wac.client_id
		LEFT JOIN
			vehicle_types vt ON vt.id = c.vehicle_type_id
		WHERE
			wac.id = ?
			AND wac.deleted_at IS NULL
	`)

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query.String()), req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Any("payload", req).Msg("repo::GetWAC - wac not found")
			return res, errmsg.NewCustomErrors(404, errmsg.WithMessage("WAC tidak ditemukan"))
		}
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWAC - failed to get wac")
		return res, err
	}

	query.Reset()
	query.WriteString(`
		SELECT
			wacc.id,
			p.id AS potency_id,
			p.name AS potency_name,
			a.id AS area_id,
			a.name AS area_name,
			a.type AS area_type,
			au.id AS assigned_user_id,
			au.name AS assigned_user_name,
			b.id AS assigned_user_branch_id,
			b.name AS assigned_user_branch_name,
			wacc.is_interested,
			wacc.path,
			wacc.notes
		FROM
			walk_around_check_conditions wacc
		LEFT JOIN
			potencies p ON p.id = wacc.potency_id
		LEFT JOIN
			areas a ON a.id = wacc.area_id
		LEFT JOIN
			users au ON au.id = wacc.assigned_user_id
		LEFT JOIN
			branches b ON b.id = au.branch_id
		WHERE
			wacc.walk_around_check_id = ?
			AND wacc.deleted_at IS NULL
	`)

	err = r.db.SelectContext(ctx, &datavc, r.db.Rebind(query.String()), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWAC - failed to get wac vehicle conditions")
		return res, err
	}

	res.Id = data.Id
	res.ClientName = data.ClientName
	res.VehicleLicenseNumber = data.VLicenseNum
	res.WhatsappNumber = data.ClientWANum
	res.IsOffered = data.IsOffered
	res.InvoiceNumber = data.InvoiceNum
	res.Revenue = data.Revenue
	res.Status = data.Status
	res.VehicleType.Id = data.VTypeId
	res.VehicleType.Name = data.VTypeName

	for _, vc := range datavc {
		var (
			filepath = strings.Split(vc.Path, "/")
			filename string
		)

		if len(filepath) > 0 {
			filename = filepath[len(filepath)-1]
		}

		res.VehicleConditions = append(res.VehicleConditions, entity.VCondition{
			Id: vc.Id,
			Potency: entity.Common{
				Id:   vc.PotencyId,
				Name: vc.PotencyName,
			},
			Area: entity.Area{
				Common: entity.Common{
					Id:   vc.AreaId,
					Name: vc.AreaName,
				},
				Type: vc.AreaType,
			},
			Assginee: entity.AssigneedUser{
				Branch: entity.Common{
					Id:   vc.AUserBranchId,
					Name: vc.AUserBranchName,
				},
				User: entity.Common{
					Id:   vc.AUserId,
					Name: vc.AUserName,
				},
			},
			Image:        storage.GenerateSignedURL(filename, 5*time.Minute),
			IsInterested: vc.IsInterested,
			Notes:        vc.Notes,
		})
	}

	return res, nil
}

// func PathToURL(path *string) *string {
// 	// devide path by /, then get the last index
// 	// then join it with the base url
// 	if path == nil {
// 		return nil
// 	}

// 	file := strings.Split(*path, "/")
// 	if len(file) == 0 {
// 		return nil
// 	}
// 	base := config.Envs.App.BaseURL
// 	url := base + "/storage/" + file[len(file)-1]

// 	return &url
// }
