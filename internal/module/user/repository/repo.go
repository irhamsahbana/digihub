package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/jwthandler"
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var _ ports.UserRepository = &userRepository{}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository() *userRepository {
	return &userRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *userRepository) Login(ctx context.Context, req *entity.LoginRequest) (entity.LoginResponse, error) {
	type dao struct {
		Id       string `db:"id"`
		Email    string `db:"email"`
		Name     string `db:"name"`
		Password string `db:"password"`
		Role     string `db:"role"`
	}

	var (
		res  entity.LoginResponse
		data dao
	)

	query := `
		SELECT
			u.id, u.email, u.name, u.password, r.name as role
		FROM
			users u
		LEFT JOIN
			roles r ON u.role_id = r.id
		WHERE email = ?
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("repo::Login - user not found")
			return res, errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid Kredensial"))
		}

		log.Error().Err(err).Msg("repo::Login - failed to get user data")
		return res, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		log.Error().Err(err).Msg("repo::Login - invalid password")
		return res, errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid Kredensial"))
	}

	var (
		expiredAt = time.Now().Add(time.Hour * 12)
		payload   = jwthandler.CostumClaimsPayload{
			UserId:          data.Id,
			Role:            data.Role,
			TokenExpiration: expiredAt,
		}
	)

	token, err := jwthandler.GenerateTokenString(payload)
	if err != nil {
		log.Error().Err(err).Msg("repo::Login - failed to generate token")
		return res, errmsg.NewCustomErrors(500, errmsg.WithMessage("Internal Server Error"))
	}

	res.Email = data.Email
	res.Name = data.Name
	res.Role = data.Role
	res.AccessToken = token
	res.ExpiredAt = expiredAt.UTC()

	return res, nil
}
