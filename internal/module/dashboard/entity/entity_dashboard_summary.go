package entity

type WACSummaryRequest struct {
	UserId   string
	UserRole string
	// month format 2021-01
	Month string `query:"month" validate:"required,datetime=2006-01"`
}

type WACSummaryResponse struct {
	Month                  string         `json:"month"`
	WACCounts              int            `json:"wac_counts" db:"wac_counts"`
	TotalWACOnOffered      int            `json:"total_wac_on_offered" db:"total_wac_on_offered"`
	TotalLeadDistributions int            `json:"total_lead_distributions" db:"total_lead_distributions"`
	Summaries              []Summary      `json:"summaries"`
	DistributionOfLeads    []Distribution `json:"distribution_of_leads"`
	ServiceTrends          []Trend        `json:"service_trends"`
	Tiers                  Tier           `json:"tiers"`
	Promotions             []Promotion    `json:"promotions"`
}

type Tier struct {
	Current string  `json:"current_tier"`
	Next    *string `json:"next_tier"`
	Revenue float64 `json:"revenue"`
}

type Promotion struct {
	Id    string `json:"id"`
	Image string `json:"image"`
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
	Total      int     `json:"total"`
}

type Trend struct {
	Types string `json:"type" db:"type"`
	Area  string `json:"area" db:"area"`
	Leads any    `json:"leads" db:"leads"` // actual type is int
}
