package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type activity struct {
	Id                  string `db:"id"`
	UserId              string `db:"user_id"`
	Status              string `db:"status"`
	TotalPotentialLeads int    `db:"total_potential_leads"`
	TotalLeads          int    `db:"total_leads"`
	TotalCompletedLeads int    `db:"total_completed_leads"`
	TotalRevenue        int    `db:"total_revenue"`
}

func (r *wacRepository) addActivity(ctx context.Context, tx *sqlx.Tx, a activity) error {
	query := `
		INSERT INTO wac_activities (
			id,
			user_id,
			status,
			total_potential_leads,
			total_leads,
			total_completed_leads,
			total_revenue
		) VALUES (
			?, ?, ?, ?, ?, ?, ?
		)
	`

	_, err := tx.ExecContext(ctx, tx.Rebind(query),
		a.Id, a.UserId, a.Status, a.TotalPotentialLeads, a.TotalLeads, a.TotalCompletedLeads, a.TotalRevenue)
	if err != nil {
		log.Error().Err(err).Msg("repository::addActivity - An error occurred")
		return err
	}

	return nil
}
