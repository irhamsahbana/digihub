package repository

import (
	"codebase-app/internal/module/common/entity"
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (r *commonRepo) CreateBranch(ctx context.Context, req *entity.CreateBranchRequest) error {
	query := `
		INSERT INTO branches (id, name, address)
		VALUES (?, TRIM(UPPER(?)), TRIM(UPPER(?)))
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		ulid.Make().String(),
		req.Name,
		req.Address,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateBranch - Failed to create branch")
		return err
	}

	return nil
}

// update
func (r *commonRepo) UpdateBranch(ctx context.Context, req *entity.UpdateBranchRequest) error {
	query := `
		UPDATE branches
		SET name = TRIM(UPPER(?)), address = TRIM(UPPER(?))
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		req.Name,
		req.Address,
		req.Id,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::UpdateBranch - Failed to update branch")
		return err
	}

	return nil
}
