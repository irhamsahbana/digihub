package repository

import (
	"codebase-app/internal/module/wac/entity"
	"codebase-app/pkg/errmsg"
	storage "codebase-app/pkg/storage-manager"
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func (r *wacRepository) GetWAC(ctx context.Context, req *entity.GetWACRequest) (entity.GetWACResponse, error) {
	type dao struct {
		Id          string  `db:"id"`
		ClientName  string  `db:"client_name"`
		VLicenseNum string  `db:"vehicle_license_number"`
		VTypeId     string  `db:"vehicle_type_id"`
		VTypeName   string  `db:"vehicle_type_name"`
		ClientWANum string  `db:"whatsapp_number"`
		IsUsedCar   bool    `db:"is_used_car"`
		IsOffered   bool    `db:"is_offered"`
		InvoiceNum  *string `db:"invoice_number"`
		Revenue     float64 `db:"revenue"`
		Status      string  `db:"status"`
	}

	type daoVCondition struct {
		Id              string  `db:"id"`
		InvoiceNum      *string `db:"invoice_number"`
		Revenue         float64 `db:"revenue"`
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
			wac.is_used_car,
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
			wacc.revenue,
			wacc.invoice_number,
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
	res.IsUsedCar = data.IsUsedCar
	res.IsOffered = data.IsOffered
	res.InvoiceNumber = data.InvoiceNum
	res.Revenue = data.Revenue
	res.Status = data.Status
	res.VehicleType.Id = data.VTypeId
	res.VehicleType.Name = data.VTypeName

	for _, vc := range datavc {
		var (
			filePath = strings.Split(vc.Path, "/")
			filename string
		)

		if len(filePath) > 0 {
			filename = filePath[len(filePath)-1]
		}

		res.VehicleConditions = append(res.VehicleConditions, entity.VCondition{
			Id:         vc.Id,
			InvoiceNum: vc.InvoiceNum,
			Revenue:    vc.Revenue,
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
			Image:        storage.GenerateSignedURL(filename, 1*time.Minute),
			IsInterested: vc.IsInterested,
			Notes:        vc.Notes,
		})
	}

	return res, nil
}
