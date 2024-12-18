package entity

import (
	"bytes"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/types"
	"time"
)

type LeadTrendsRequest struct {
	UserId string
}

type LeadTrendsResponse struct {
	Month            string `json:"month" db:"month"`
	ReviewConditions int    `json:"review_conditions" db:"review_conditions"`
	Leads            int    `json:"leads" db:"leads"`
}

type TechWACSummaryResponse struct {
	Month                string         `json:"month"`
	TotalWACNeedFollowUp int            `json:"total_wac_need_follow_up" db:"total_wac_need_follow_up"`
	TotalWACFollowedUp   int            `json:"total_wac_followed_up" db:"total_wac_followed_up"`
	TotalLeads           int            `json:"total_leads" db:"total_leads"`
	DistributionOfLeads  []Distribution `json:"distribution_of_leads"`
}

type GetActivitiesRequest struct {
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
	Search   string `query:"search" validate:"omitempty,min=3"`
	From     string `query:"from" validate:"omitempty,datetime=2006-01-02"`
	To       string `query:"to" validate:"omitempty,datetime=2006-01-02"`
	Timezone string `query:"timezone" validate:"omitempty,timezone"`
	BranchId string `query:"branch_id" validate:"omitempty,ulid,exist=branches.id"`
	Export   uint   `query:"export" validate:"omitempty,min=0,max=1"`

	FromTime   time.Time
	ToTime     time.Time
	BranchName string
}

func (r *GetActivitiesRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}

	if r.Timezone == "" {
		r.Timezone = "Asia/Makassar"
	}
}

func (r *GetActivitiesRequest) Validate() error {
	errs := errmsg.NewCustomErrors(400)

	if (r.From == "" && r.To != "") || (r.From != "" && r.To == "") {
		errs.Add("from", "from dan to harus diisi bersamaan")
		errs.Add("to", "from dan to harus diisi bersamaan")
	}

	if r.From != "" && r.To != "" {
		from, err := time.Parse("2006-01-02", r.From)
		if err != nil {
			errs.Add("from", "from merupakan tanggal yang tidak valid")
		}
		r.FromTime = from

		to, err := time.Parse("2006-01-02", r.To)
		if err != nil {
			errs.Add("to", "to merupakan tanggal yang tidak valid")
		}
		r.ToTime = to

		if from.After(to) {
			errs.Add("from", "from seharusnya tidak lebih besar dari to")
		}

		if to.Before(from) {
			errs.Add("to", "to seharusnya tidak lebih kecil dari from")
		}
	}

	if r.Export == 1 {
		if r.From == "" || r.To == "" {
			errs.Add("from", "from harus diisi ketika export")
			errs.Add("to", "to harus diisi ketika export")
		}

		if r.BranchId == "" {
			errs.Add("branch_id", "branch_id harus diisi ketika export")
		}
	}

	if errs.HasErrors() {
		return errs
	} else {
		return nil
	}
}

type GetActivitiesResponse struct {
	Items []Activity `json:"items"`
	Meta  types.Meta `json:"meta"`
}

type Activity struct {
	Id                   string    `json:"id" db:"id"`
	EmployeeName         string    `json:"employee_name" db:"employee_name"`
	ClientName           string    `json:"client_name" db:"client_name"`
	BranchName           string    `json:"branch_name" db:"branch_name"`
	VehicleLicenseNumber string    `json:"vehicle_license_number" db:"vehicle_license_number"`
	Phone                string    `json:"phone" db:"phone"`
	VehicleTypeName      string    `json:"vehicle_type_name" db:"vehicle_type_name"`
	Status               string    `json:"status" db:"status"`
	TotalPotentialLeads  int       `json:"total_potential_leads" db:"total_potential_leads"`
	TotalLeads           int       `json:"total_leads" db:"total_leads"`
	TotalRevenue         float64   `json:"total_revenue" db:"total_revenue"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

type ExportMeta struct {
	Filename string
	Buf      *bytes.Buffer
}
