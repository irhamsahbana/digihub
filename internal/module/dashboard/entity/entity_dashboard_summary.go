package entity

type WACSummaryRequest struct {
	UserId string
	// month format 2021-01
	Month string `query:"month" validate:"required,datetime=2006-01"`
}

type WACSummaryResponse struct {
	Month                  string         `json:"month"`
	WACCounts              int            `json:"wac_counts" db:"wac_counts"`
	TotalLeadDistributions int            `json:"total_lead_distributions" db:"total_lead_distributions"`
	Summaries              []Summary      `json:"summaries"`
	DistributionOfLeads    []Distribution `json:"distribution_of_leads"`
}

type Summary struct {
	Title               string `json:"title" db:"title"`
	TotalPotencialLeads int    `json:"total_potencial_leads" db:"total_potencial_leads"`
	TotalLeads          int    `json:"total_leads" db:"total_leads"`
	TotalWoDo           int    `json:"total_wo_do" db:"total_wo_do"`
}

type Distribution struct {
	Title      string  `json:"title"`
	Percentage float64 `json:"percentage"`
}
