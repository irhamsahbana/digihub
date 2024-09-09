package entity

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
