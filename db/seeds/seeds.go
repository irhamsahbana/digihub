package seeds

import (
	"codebase-app/internal/adapter"
	"context"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

// Seed struct.
type Seed struct {
	db *sqlx.DB
}

// NewSeed return a Seed with a pool of connection to a dabase.
func newSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

func Execute(db *sqlx.DB, table string, total int) {
	seed := newSeed(db)
	seed.run(table, total)
}

// Run seeds.
func (s *Seed) run(table string, total int) {

	switch table {
	case "roles":
		s.rolesSeed()
	case "potencies":
		s.potenciesSeed()
	case "areas":
		s.areasSeed()
	case "vehicle_types":
		s.vehicleTypesSeed()
	case "branches":
		s.branchesSeed(total)
	case "users":
		s.usersSeed(total)
	case "all":
		s.rolesSeed()
		s.vehicleTypesSeed()
		s.potenciesSeed()
		s.areasSeed()
		s.branchesSeed(total)
		s.usersSeed(total)
	case "delete-all":
		s.deleteAll()
	default:
		log.Warn().Msg("No seed to run")
	}

	if table != "" {
		log.Info().Msg("Seed ran successfully")
		log.Info().Msg("Exiting ...")
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
		os.Exit(0)
	}
}

func (s *Seed) deleteAll() {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.Exec(`DELETE FROM users`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting users")
		return
	}
	log.Info().Msg("users table deleted successfully")

	_, err = tx.Exec(`DELETE FROM branches`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting branches")
		return
	}
	log.Info().Msg("branches table deleted successfully")

	_, err = tx.Exec(`DELETE FROM roles`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting roles")
		return
	}
	log.Info().Msg("roles table deleted successfully")

	_, err = tx.Exec(`DELETE FROM potencies`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting potencies")
		return
	}
	log.Info().Msg("potencies table deleted successfully")

	_, err = tx.Exec(`DELETE FROM vehicle_types`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting vehicle types")
		return
	}
	log.Info().Msg("vehicle types table deleted successfully")

	_, err = tx.Exec(`DELETE FROM areas`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting areas")
		return
	}
	log.Info().Msg("areas table deleted successfully")

	log.Info().Msg("=== All tables deleted successfully ===")
}

func (s *Seed) areasSeed() {
	areaMaps := []map[string]any{
		{"id": ulid.Make().String(), "name": "Kap Mesin", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Bumper Depan", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Pintu kanan depan", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Pintu kanan tengah", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Pintu kanan belakang", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Pintu kiri depan", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Pintu kiri tengah", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Pintu kiri belakang", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Pintu belakang", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Bumper Belakang", "type": "exterior"},
		{"id": ulid.Make().String(), "name": "Atap Ruang Mesin", "type": "exterior"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO areas (id, name, type)
		VALUES (:id, :name, :type)
	`, areaMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating areas")
		return
	}

	log.Info().Msg("areas table seeded successfully")
}

func (s *Seed) potenciesSeed() {
	potencyMaps := []map[string]any{
		{"id": "01J3VHA25VY2WPBNSNP9RA61ST", "name": "Leads to General Review"},
		{"id": "01J3VHA25VY2WPBNSNP9XZDWAG", "name": "Leads to Body Paint"},
		{"id": "01J3VHA25VY2WPBNSNPB9P9FQ3", "name": "Leads to Used-car"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO potencies (id, name)
		VALUES (:id, :name)
	`, potencyMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating potencies")
		return
	}

	log.Info().Msg("potencies table seeded successfully")
}

func (s *Seed) vehicleTypesSeed() {
	vehicleTypeMaps := []map[string]any{
		{"id": ulid.Make().String(), "name": "AGYA"},
		{"id": ulid.Make().String(), "name": "ALPHARD"},
		{"id": ulid.Make().String(), "name": "AVANZA"},
		{"id": ulid.Make().String(), "name": "CALYA"},
		{"id": ulid.Make().String(), "name": "CAMRY"},
		{"id": ulid.Make().String(), "name": "COROLLA"},
		{"id": ulid.Make().String(), "name": "DYNA"},
		{"id": ulid.Make().String(), "name": "FORTUNER"},
		{"id": ulid.Make().String(), "name": "HIACE"},
		{"id": ulid.Make().String(), "name": "HILUX"},
		{"id": ulid.Make().String(), "name": "INNOVA"},
		{"id": ulid.Make().String(), "name": "KIJANG"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO vehicle_types (id, name)
		VALUES (:id, :name)
	`, vehicleTypeMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating vehicle types")
		return
	}

	log.Info().Msg("vehicle types table seeded successfully")
}

// rolesSeed seeds the roles table.
func (s *Seed) rolesSeed() {
	roleMaps := []map[string]any{
		{"id": "01J3VHA25R8KTG9MQX43KBZ9MW", "name": "admin"},
		{"id": "01J3VHA25R8KTG9MQX45R8F3V7", "name": "service_advisor"},
		{"id": "01J3VHA25R8KTG9MQX47GRF4KW", "name": "technician"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO roles (id, name)
		VALUES (:id, :name)
	`, roleMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating roles")
		return
	}

	log.Info().Msg("roles table seeded successfully")
}

func (s *Seed) branchesSeed(total int) {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Error().Err(err).Msg("Error rolling back transaction")
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	branchMaps := make([]map[string]any, 0)

	for i := 0; i < total; i++ {
		dataBranchToInsert := make(map[string]any)
		dataBranchToInsert["id"] = ulid.Make().String()
		dataBranchToInsert["name"] = gofakeit.City()
		dataBranchToInsert["address"] = gofakeit.Address().Address

		branchMaps = append(branchMaps, dataBranchToInsert)
	}

	_, err = tx.NamedExec(`
		INSERT INTO branches (id, name, address)
		VALUES (:id, :name, :address)
	`, branchMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating branches")
		return
	}

	log.Info().Msg("branches table seeded successfully")
}

// users
func (s *Seed) usersSeed(total int) {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	type roleData struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	type branchData struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	var (
		roles    = make([]roleData, 0)
		branches = make([]branchData, 0)
		userMaps = make([]map[string]any, 0)
	)

	err = s.db.Select(&roles, `SELECT id, name FROM roles`)
	if err != nil {
		log.Error().Err(err).Msg("Error selecting roles")
		return
	}

	err = s.db.Select(&branches, `SELECT id, name FROM branches`)
	if err != nil {
		log.Error().Err(err).Msg("Error selecting branches")
		return
	}

	for i := 0; i < total; i++ {
		selectedRole := roles[gofakeit.Number(0, len(roles)-1)]
		selectedBranch := branches[gofakeit.Number(0, len(branches)-1)]

		dataUserToInsert := make(map[string]any)
		dataUserToInsert["id"] = ulid.Make().String()
		dataUserToInsert["role_id"] = selectedRole.Id
		dataUserToInsert["branch_id"] = selectedBranch.Id
		dataUserToInsert["name"] = gofakeit.Name()
		dataUserToInsert["email"] = gofakeit.Email()
		dataUserToInsert["password"] = "$2y$10$mVf4BKsfPSh/pjgHjvk.JOlGdkIYgBGyhaU9WQNMWpYskK9MZlb0G" // password

		userMaps = append(userMaps, dataUserToInsert)
	}

	_, err = tx.NamedExec(`
		INSERT INTO users (id, role_id, branch_id, name, email, password)
		VALUES (:id, :role_id, :branch_id, :name, :email, :password)
	`, userMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating users")
		return
	}

	log.Info().Msg("users table seeded successfully")
}
