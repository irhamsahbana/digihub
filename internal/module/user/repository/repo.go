package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/jwthandler"
	"codebase-app/pkg/storage-manager"
	"context"
	"database/sql"
	"strings"
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

func (r *userRepository) GetProfile(ctx context.Context, req *entity.GetProfileRequest) (entity.GetProfileResponse, error) {
	type dao struct {
		Id          string  `db:"id"`
		Name        string  `db:"name"`
		Email       string  `db:"email"`
		WANum       string  `db:"whatsapp_number"`
		Path        *string `db:"path"`
		BranchId    string  `db:"branch_id"`
		BranchName  string  `db:"branch_name"`
		SectionId   string  `db:"section_id"`
		SectionName string  `db:"section_name"`
		RoleId      string  `db:"role_id"`
		RoleName    string  `db:"role_name"`
	}

	var (
		res  entity.GetProfileResponse
		data = new(dao)
	)

	query := `
		SELECT
			u.id,
			u.name,
			u.email,
			u.whatsapp_number,
			u.path,
			b.id as branch_id,
			b.name as branch_name,
			s.id as section_id,
			s.name as section_name,
			r.id as role_id,
			r.name as role_name
		FROM
			users u
		LEFT JOIN
			branches b ON u.branch_id = b.id
		LEFT JOIN
			sections s ON u.section_id = s.id
		LEFT JOIN
			roles r ON u.role_id = r.id
		WHERE u.id = ?
	`

	err := r.db.GetContext(ctx, data, r.db.Rebind(query), req.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Any("payload", req).Msg("repo::GetProfile - user not found")
			return res, errmsg.NewCustomErrors(404, errmsg.WithMessage("User not found"))
		}
		log.Error().Err(err).Any("payload", req).Msg("repo::GetProfile - failed to get user data")
		return res, err
	}

	res.Branch.Id = data.BranchId
	res.Branch.Name = data.BranchName
	res.Section.Id = data.SectionId
	res.Section.Name = data.SectionName
	res.Role.Id = data.RoleId
	res.Role.Name = data.RoleName
	res.Id = data.Id
	res.Name = data.Name
	res.Email = data.Email
	res.WANum = data.WANum

	// check if user has profile picture
	// if yes, generate public url
	if data.Path != nil {
		var (
			filePath = strings.Split(*data.Path, "/")
			filename string
		)

		if len(filePath) > 0 {
			filename = filePath[len(filePath)-1]
		}

		url := storage.GeneratePublicURL(filename)
		res.Image = &url
	}

	return res, nil
}
