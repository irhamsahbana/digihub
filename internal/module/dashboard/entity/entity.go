package entity

type LeadTrendsRequest struct {
	UserId string
}

type LeadTrendsResponse struct {
	Month            string `json:"month" db:"month"`
	ReviewConditions int    `json:"review_conditions" db:"review_conditions"`
	Leads            int    `json:"leads" db:"leads"`
}

type WACSummaryRequest struct {
	UserId string
}

type WACSummaryResponse struct {
	WACCounts           int            `json:"wac_counts"`
	Summaries           []Summary      `json:"summaries"`
	DistributionOfLeads []Distribution `json:"distribution_of_leads"`
}

type Summary struct {
	Title               string `json:"title"`
	TotalPotencialLeads int    `json:"total_potencial_leads"`
	TotalLeads          int    `json:"total_leads"`
	TotalWoDo           int    `json:"total_wo_do"`
}

type Distribution struct {
	Title      string  `json:"title"`
	Percentage float64 `json:"percentage"`
}
