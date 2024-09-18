package seeds

import (
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/xuri/excelize/v2"
)

type excelSeed struct {
	db   *sqlx.DB
	file *excelize.File
}

func newExcelSeed(db *sqlx.DB) (*excelSeed, error) {
	f, err := excelize.OpenFile("db/seeds/excel/seeds.xlsx")
	if err != nil {
		log.Error().Err(err).Msg("failed to open excel file")
		return nil, err
	}

	return &excelSeed{db: db, file: f}, nil
}

func SeedExcel(db *sqlx.DB, sheetName string) error {
	excelSeeder, err := newExcelSeed(db)
	if err != nil {
		log.Error().Err(err).Msg("failed to create excel seeder")
		return err
	}

	tx, err := excelSeeder.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
	}
	defer tx.Rollback()
	var errSeed error

	switch sheetName {
	case "roles":
		errSeed = excelSeeder.SeedRoles(tx)
		if errSeed != nil {
			return errSeed
		}
	case "branches":
		errSeed = excelSeeder.SeedBranches(tx)
		if errSeed != nil {
			return errSeed
		}
	case "potencies":
		errSeed = excelSeeder.SeedPotencies(tx)
		if errSeed != nil {
			return errSeed
		}
	case "areas":
		errSeed = excelSeeder.SeedAreas(tx)
		if errSeed != nil {
			return errSeed
		}
	case "vehicle_types":
		errSeed = excelSeeder.SeedVehicleTypes(tx)
		if errSeed != nil {
			return errSeed
		}
	case "trade_in_trends":
		errSeed = excelSeeder.SeedTradeInTrends(tx)
		if errSeed != nil {
			return errSeed
		}
	case "users":
		errSeed = excelSeeder.SeedUsers(tx)
		if errSeed != nil {
			return errSeed
		}
	case "fill_vehicle_types":
		errSeed = excelSeeder.FillVehicleTypes(tx)
		if errSeed != nil {
			return errSeed
		}
	}

	if err := excelSeeder.file.Save(); err != nil {
		log.Error().Err(err).Msg("failed to save excel file")
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction")
		return err
	}

	return nil
}

func (s *excelSeed) SeedRoles(tx *sqlx.Tx) error {
	rows, err := s.file.GetRows("roles")
	if err != nil {
		log.Error().Err(err).Msg("failed to get rows from excel")
		return err
	}

	idsInSheet := make([]string, len(rows)-1)
	lastRow := len(rows) - 1
	// insert into db
	for i, row := range rows {
		if i == 0 { // skip header
			continue
		}

		var (
			id   = row[0]
			name = row[1]
		)

		// if id is empty then add ULID to it, and when done it should be saved to db
		// and file
		if id == "" {
			id = ulid.Make().String()
			rowNumber := strconv.Itoa(i + 1)
			cell := "A" + rowNumber
			s.file.SetCellValue("roles", cell, id)
		}

		// check id is valid ULID
		if _, err := ulid.Parse(id); err != nil {
			log.Error().Err(err).Msg("invalid ULID")
			return err
		}

		idsInSheet[i-1] = id

		query := "INSERT INTO roles (id, name) VALUES (?, ?) ON CONFLICT DO NOTHING"
		_, err := tx.Exec(s.db.Rebind(query), id, name)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert role")
			return err
		}
	}

	// get all roles from db that are not in the sheet
	type role struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}
	var rolesNotInSheet []role

	query, args, err := sqlx.In("SELECT id, name FROM roles WHERE id NOT IN (?)", idsInSheet)
	if err != nil {
		log.Error().Err(err).Msg("failed to create query for roles not in sheet")
		return err
	}
	err = tx.Select(&rolesNotInSheet, s.db.Rebind(query), args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to get roles not in sheet")
		return err
	}

	// append roles not in sheet to the sheet
	for i, role := range rolesNotInSheet {
		rowNumber := strconv.Itoa(lastRow + i + 2)
		cellA := "A" + rowNumber
		cellB := "B" + rowNumber
		s.file.SetCellValue("roles", cellA, role.Id)
		s.file.SetCellValue("roles", cellB, role.Name)
	}

	log.Info().Msg("roles seeded successfully!")

	return nil
}

func (s *excelSeed) SeedBranches(tx *sqlx.Tx) error {
	rows, err := s.file.GetRows("branches")
	if err != nil {
		log.Error().Err(err).Msg("failed to get rows from excel")
		return err
	}

	idsInSheet := make([]string, len(rows)-1)
	lastRow := len(rows) - 1
	// insert into db
	for i, row := range rows {
		if i == 0 { // skip header
			continue
		}

		var (
			id   = row[0]
			name = row[1]
		)

		// if id is empty then add ULID to it, and when done it should be saved to db
		// and file
		if id == "" {
			id = ulid.Make().String()
			rowNumber := strconv.Itoa(i + 1)
			cell := "A" + rowNumber
			s.file.SetCellValue("branches", cell, id)
		}

		// check id is valid ULID
		if _, err := ulid.Parse(id); err != nil {
			log.Error().Err(err).Msg("invalid ULID")
			return err
		}

		idsInSheet[i-1] = id

		query := "INSERT INTO branches (id, name) VALUES (?, ?) ON CONFLICT DO NOTHING"
		_, err := tx.Exec(s.db.Rebind(query), id, name)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert branch")
			return err
		}
	}

	// get all branches from db that are not in the sheet
	type branch struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	var branchesNotInSheet []branch

	query, args, err := sqlx.In("SELECT id, name FROM branches WHERE id NOT IN (?)", idsInSheet)
	if err != nil {
		log.Error().Err(err).Msg("failed to create query for branches not in sheet")
		return err
	}

	err = tx.Select(&branchesNotInSheet, s.db.Rebind(query), args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to get branches not in sheet")
		return err
	}

	// append branches not in sheet to the sheet
	for i, branch := range branchesNotInSheet {
		rowNumber := strconv.Itoa(lastRow + i + 2) // +2 because 1 based index and header
		cellA := "A" + rowNumber
		cellB := "B" + rowNumber
		s.file.SetCellValue("branches", cellA, branch.Id)
		s.file.SetCellValue("branches", cellB, branch.Name)
	}

	log.Info().Msg("branches seeded successfully!")

	return nil
}

func (s *excelSeed) SeedPotencies(tx *sqlx.Tx) error {
	rows, err := s.file.GetRows("potencies_sections")
	if err != nil {
		log.Error().Err(err).Msg("failed to get rows from excel")
		return err
	}

	idsInSheet := make([]string, len(rows)-1)
	lastRow := len(rows) - 1
	// insert into db
	for i, row := range rows {
		if i == 0 { // skip header
			continue
		}

		var (
			id   = row[0]
			name = row[1]
		)

		// if id is empty then add ULID to it, and when done it should be saved to db
		// and file
		if id == "" {
			id = ulid.Make().String()
			rowNumber := strconv.Itoa(i + 1)
			cell := "A" + rowNumber
			s.file.SetCellValue("potencies_sections", cell, id)
		}

		// check id is valid ULID
		if _, err := ulid.Parse(id); err != nil {
			log.Error().Err(err).Msg("invalid ULID")
			return err
		}

		idsInSheet[i-1] = id

		query := "INSERT INTO potencies (id, name) VALUES (?, ?) ON CONFLICT DO NOTHING"
		_, err := tx.Exec(s.db.Rebind(query), id, name)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert potency")
			return err
		}
	}

	// get all potencies from db that are not in the sheet
	type potency struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}
	var potenciesNotInSheet []potency

	query, args, err := sqlx.In("SELECT id, name FROM potencies WHERE id NOT IN (?)", idsInSheet)
	if err != nil {
		log.Error().Err(err).Msg("failed to create query for potencies not in sheet")
		return err
	}
	err = tx.Select(&potenciesNotInSheet, s.db.Rebind(query), args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to get potencies not in sheet")
		return err
	}

	// append potencies not in sheet to the sheet
	for i, potency := range potenciesNotInSheet {
		rowNumber := strconv.Itoa(lastRow + i + 2)
		cellA := "A" + rowNumber
		cellB := "B" + rowNumber
		s.file.SetCellValue("potencies_sections", cellA, potency.Id)
		s.file.SetCellValue("potencies_sections", cellB, potency.Name)
	}

	log.Info().Msg("potencies seeded successfully!")

	return nil
}

func (s *excelSeed) SeedAreas(tx *sqlx.Tx) error {
	rows, err := s.file.GetRows("areas")
	if err != nil {
		log.Error().Err(err).Msg("failed to get rows from excel")
		return err
	}

	idsInSheet := make([]string, len(rows)-1)
	lastRow := len(rows) - 1
	// insert into db
	for i, row := range rows {
		if i == 0 { // skip header
			continue
		}

		var (
			id    = row[0]
			types = row[1]
			name  = row[2]
		)

		// if id is empty then add ULID to it, and when done it should be saved to db
		// and file
		if id == "" {
			id = ulid.Make().String()
			rowNumber := strconv.Itoa(i + 1)
			cell := "A" + rowNumber
			s.file.SetCellValue("areas", cell, id)
		}

		// check id is valid ULID
		if _, err := ulid.Parse(id); err != nil {
			log.Error().Err(err).Msg("invalid ULID")
			return err
		}

		idsInSheet[i-1] = id

		query := "INSERT INTO areas (id, type, name) VALUES (?, ?, ?) ON CONFLICT DO NOTHING"
		_, err := tx.Exec(s.db.Rebind(query), id, types, name)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert area")
			return err
		}
	}

	// get all areas from db that are not in the sheet
	type area struct {
		Id   string `db:"id"`
		Type string `db:"type"`
		Name string `db:"name"`
	}
	var areasNotInSheet []area

	query, args, err := sqlx.In("SELECT id, type, name FROM areas WHERE id NOT IN (?)", idsInSheet)
	if err != nil {
		log.Error().Err(err).Msg("failed to create query for areas not in sheet")
		return err
	}
	err = tx.Select(&areasNotInSheet, s.db.Rebind(query), args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to get areas not in sheet")
		return err
	}

	// append areas not in sheet to the sheet
	for i, area := range areasNotInSheet {
		rowNumber := strconv.Itoa(lastRow + i + 2)
		cellA := "A" + rowNumber
		cellB := "B" + rowNumber
		cellC := "C" + rowNumber
		s.file.SetCellValue("areas", cellA, area.Id)
		s.file.SetCellValue("areas", cellB, area.Type)
		s.file.SetCellValue("areas", cellC, area.Name)
	}

	log.Info().Msg("areas seeded successfully!")

	return nil
}

func (s *excelSeed) SeedVehicleTypes(tx *sqlx.Tx) error {
	rows, err := s.file.GetRows("vehicle_types")
	if err != nil {
		log.Error().Err(err).Msg("failed to get rows from excel")
		return err
	}

	idsInSheet := make([]string, len(rows)-1)
	lastRow := len(rows) - 1
	// insert into db
	for i, row := range rows {
		if i == 0 { // skip header
			continue
		}

		var (
			id   = row[0]
			name = row[1]
		)

		// if id is empty then add ULID to it, and when done it should be saved to db
		// and file
		if id == "" {
			id = ulid.Make().String()
			rowNumber := strconv.Itoa(i + 1)
			cell := "A" + rowNumber
			s.file.SetCellValue("vehicle_types", cell, id)
		}

		// check id is valid ULID
		if _, err := ulid.Parse(id); err != nil {
			log.Error().Err(err).Msg("invalid ULID")
			return err
		}

		idsInSheet[i-1] = id

		query := "INSERT INTO vehicle_types (id, name) VALUES (?, ?) ON CONFLICT DO NOTHING"
		_, err := tx.Exec(s.db.Rebind(query), id, name)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert vehicle type")
			return err
		}
	}

	// get all vehicle types from db that are not in the sheet
	type vehicleType struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	var vehicleTypesNotInSheet []vehicleType

	query, args, err := sqlx.In("SELECT id, name FROM vehicle_types WHERE id NOT IN (?)", idsInSheet)
	if err != nil {
		log.Error().Err(err).Msg("failed to create query for vehicle types not in sheet")
		return err
	}

	err = tx.Select(&vehicleTypesNotInSheet, s.db.Rebind(query), args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to get vehicle types not in sheet")
		return err
	}

	// append vehicle types not in sheet to the sheet
	for i, vehicleType := range vehicleTypesNotInSheet {
		rowNumber := strconv.Itoa(lastRow + i + 2)
		cellA := "A" + rowNumber
		cellB := "B" + rowNumber
		s.file.SetCellValue("vehicle_types", cellA, vehicleType.Id)
		s.file.SetCellValue("vehicle_types", cellB, vehicleType.Name)
	}

	log.Info().Msg("vehicle types seeded successfully!")

	return nil
}

func (s *excelSeed) FillVehicleTypes(tx *sqlx.Tx) error {
	rowsHTI, err := s.file.GetRows("hi_trade_in")
	if err != nil {
		log.Error().Err(err).Msg("failed to get rows from excel")
		return err
	}

	rowsVT, err := s.file.GetRows("vehicle_types")
	if err != nil {
		log.Error().Err(err).Msg("failed to get rows from excel")
		return err
	}

	var (
		vehicleModel    = make(map[string]struct{})
		vehicleTypeName = make(map[string]struct{})
	)

	for i, row := range rowsVT {
		if i == 0 { // skip header
			continue
		}

		var (
			id   = row[0]
			name = row[1]
		)

		if id == "" {
			id = ulid.Make().String()
			rowNumber := strconv.Itoa(i + 1)
			cell := "A" + rowNumber
			s.file.SetCellValue("vehicle_types", cell, id)
		}

		name = strings.ToUpper(strings.Trim(name, " "))
		vehicleTypeName[name] = struct{}{}
	}

	for i, row := range rowsHTI {
		if i == 0 { // skip header
			continue
		}

		var (
			// brand = row[1]
			model = row[2]
			// type_ = row[3]
			// year  = row[4]
		)

		model = strings.ToUpper(strings.Trim(model, " "))
		vehicleModel[model] = struct{}{}
	}

	// loop through vehicleModel and check if it exists in vehicleTypeName
	// if exists then continue, if not then add it to vehicle_types sheet

	lastRowInVT := len(rowsVT) - 1 // 1 based index

	for model := range vehicleModel {
		if _, ok := vehicleTypeName[model]; ok {
			continue
		}

		// add to vehicle_types sheet
		rowNumber := strconv.Itoa(lastRowInVT + 2)
		cellA := "A" + rowNumber
		cellB := "B" + rowNumber
		s.file.SetCellValue("vehicle_types", cellA, ulid.Make().String())
		s.file.SetCellValue("vehicle_types", cellB, model)

		lastRowInVT++
	}

	log.Info().Msg("vehicle types filled successfully!")

	return nil
}

func (s *excelSeed) SeedUsers(tx *sqlx.Tx) error {
	rows, err := s.file.GetRows("users")
	if err != nil {
		log.Error().Err(err).Msg("failed to get rows from excel")
		return err
	}

	// insert into db
	for i, row := range rows {
		if i == 0 { // skip header
			continue
		}

		var (
			id          = row[0]
			name        = row[1]
			branchName  = row[2]
			sectionName = row[3]
			RoleName    = row[4]
			email       = row[5]
			password    = row[6]
		)

		// manipulate data
		switch strings.ToUpper(RoleName) {
		case "SERVICE ADVISOR":
			RoleName = "service_advisor"
		case "MRA":
			RoleName = "technician"
		}

		// bcrypt password
		passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Error().Err(err).Msg("failed to hash password")
			return err
		}

		// if id is empty then add ULID to it, and when done it should be saved to db
		// and file
		if id == "" {
			id = ulid.Make().String()
			rowNumber := strconv.Itoa(i + 1)
			cell := "A" + rowNumber
			s.file.SetCellValue("users", cell, id)
		}

		// make sure email is in lowercase
		email = strings.ToLower(email)

		// check id is valid ULID
		if _, err := ulid.Parse(id); err != nil {
			log.Error().Err(err).Msg("invalid ULID")
			return err
		}

		query := `
			INSERT INTO users (
				id, name, branch_id, section_id, role_id, email, password
			) VALUES (
				?,
				?,
				(SELECT id FROM branches WHERE UPPER(name) = UPPER(?)),
				(SELECT id FROM potencies WHERE UPPER(name) = UPPER(?)),
				(SELECT id FROM roles WHERE name = ?),
				?,
				?
			) ON CONFLICT DO NOTHING
		`

		_, err = tx.Exec(s.db.Rebind(query), id, name, branchName, sectionName, RoleName, email, string(passwordHashed))
		if err != nil {
			log.Error().Err(err).Msg("failed to insert user")
			return err
		}
	}

	log.Info().Msg("users seeded successfully!")

	return nil
}

func (s *excelSeed) SeedTradeInTrends(tx *sqlx.Tx) error {
	rows, err := s.file.GetRows("hi_trade_in")
	if err != nil {
		log.Error().Err(err).Msg("failed to get rows from excel")
		return err
	}

	for i, row := range rows {
		if i == 0 { // skip header
			continue
		}

		var (
			brand = strings.ToUpper(strings.Trim(row[1], " "))
			model = strings.ToUpper(strings.Trim(row[2], " "))
			type_ = strings.ToUpper(strings.Trim(row[3], " "))
		)

		// try to convert string to int
		year, err := strconv.Atoi(row[4])
		if err != nil {
			log.Error().Err(err).Msg("failed to convert year to int")
			return err
		}

		minPurchase, err := strconv.Atoi(row[5])
		if err != nil {
			log.Error().Err(err).Msg("failed to convert min purchase to int")
			return err
		}

		maxPurchase, err := strconv.Atoi(row[6])
		if err != nil {
			log.Error().Err(err).Msg("failed to convert max purchase to int")
			return err
		}

		// on conflict do update min_purchase and max_purchase
		query := `
			INSERT INTO trade_in_trends (brand, model, type, year, min_purchase, max_purchase)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT (brand, model, type, year)
			DO UPDATE SET min_purchase = ?, max_purchase = ?, created_at = NOW()
			`
		_, err = tx.Exec(s.db.Rebind(query),
			brand, model, type_, year, minPurchase, maxPurchase,
			minPurchase, maxPurchase,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert hi tread in trend")
			return err
		}
	}

	log.Info().Msg("hi tread in trends seeded successfully!")

	return nil
}
