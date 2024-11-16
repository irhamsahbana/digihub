package service

import (
	"codebase-app/internal/module/employee/entity"
	"codebase-app/internal/module/employee/ports"
	"codebase-app/pkg"
	"codebase-app/pkg/errmsg"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

var _ ports.EmployeeService = &employeeService{}

type employeeService struct {
	repo ports.EmployeeRepository
}

func NewEmployeeService(repo ports.EmployeeRepository) *employeeService {
	return &employeeService{
		repo: repo,
	}
}

func (s *employeeService) GetEmployee(ctx context.Context, req *entity.GetEmployeeRequest) (entity.GetEmployeeResponse, error) {
	res, err := s.repo.GetEmployee(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, req *entity.UpdateEmployeeRequest) error {
	err := s.repo.UpdateEmployee(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *employeeService) CreateEmployee(ctx context.Context, req *entity.CreateEmployeeRequest) error {
	err := s.repo.CreateEmployee(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *employeeService) DeleteEmployee(ctx context.Context, req *entity.DeleteEmployeeRequest) error {
	err := s.repo.DeleteEmployee(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *employeeService) ImportEmployees(ctx context.Context, req *entity.ImportEmployeesRequest) error {
	var dataImport []entity.ImportEmployeeRow
	var mapBranches = make(map[string]string)
	var mapSections = make(map[string]string)
	var mapRoles = make(map[string]string)
	var mapEmails = make(map[string]bool)

	dataBranches, err := s.repo.GetBranches(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - Failed to get branches")
		return errmsg.NewCustomErrors(500).SetMessage("Gagal mendapatkan data cabang")
	}

	dataPotencies, err := s.repo.GetPotencies(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - Failed to get potencies")
		return errmsg.NewCustomErrors(500).SetMessage("Gagal mendapatkan data section")
	}

	dataRoles, err := s.repo.GetRoles(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - Failed to get roles")
		return errmsg.NewCustomErrors(500).SetMessage("Gagal mendapatkan data role")
	}

	for _, branch := range dataBranches {
		mapBranches[strings.ToUpper(branch.Name)] = branch.Id
	}

	for _, potencies := range dataPotencies {
		mapSections[strings.ToUpper(potencies.Name)] = potencies.Id
	}

	for _, role := range dataRoles {
		mapRoles[strings.ToUpper(role.Name)] = role.Id
	}

	// Decode base64 to xlsx file
	data, err := base64.StdEncoding.DecodeString(req.File)
	if err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - Failed to decode base64 string")
		return errmsg.NewCustomErrors(400).SetMessage("Gagal mendecode file base64")
	}

	// validate file is xlsx
	err = isXLSXFile(data)
	if err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - File is not a valid xlsx file")
		return errmsg.NewCustomErrors(400).SetMessage("File bukan file xlsx yang valid")
	}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "employees-*.xlsx")
	if err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - Failed to create temporary file")
		return errmsg.NewCustomErrors(500).SetMessage("Gagal membuat file temporary")
	}
	defer os.Remove(tmpFile.Name())

	// Write the decoded data to the temporary file
	if _, err := tmpFile.Write(data); err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - Failed to write to temporary file")
		return errmsg.NewCustomErrors(500).SetMessage("Gagal menulis ke file temporary")
	}
	if err := tmpFile.Close(); err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - Failed to close temporary file")
		return errmsg.NewCustomErrors(500).SetMessage("Gagal menutup file temporary")
	}

	// Read with excelize
	xlsx, err := excelize.OpenFile(tmpFile.Name())
	if err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - Failed to open xlsx file")
		return errmsg.NewCustomErrors(500).SetMessage("Gagal membuka file xlsx")
	}

	// Get all the rows in the import_users sheet
	rows, err := xlsx.GetRows("import_users")
	if err != nil {
		log.Error().Err(err).Msg("service::ImportEmployees - Failed to get rows from xlsx file")
		if errSheetNotExist, ok := err.(*excelize.ErrSheetNotExist); ok {
			return errmsg.NewCustomErrors(400).SetMessage("Sheet " + errSheetNotExist.SheetName + " tidak ditemukan")
		}
	}

	errExcelValidation := errmsg.NewCustomErrors(400).SetMessage("Beberapa data pegawai tidak valid")

	// Iterate over the rows
	for i, row := range rows {
		var r entity.ImportEmployeeRow

		// Skip the header
		if i == 0 {
			continue
		}

		rowname := fmt.Sprintf("file[%d]", i)

		var (
			name        = row[1]
			branchName  = strings.ToUpper(row[2])
			sectionName = strings.ToUpper(row[3])
			roleName    = strings.ToUpper(row[4])
			email       = strings.ToLower(row[5])
			password    = row[6]
		)

		r.BranchId = mapBranches[branchName]
		r.SectionId = mapSections[sectionName]
		r.BranchName = branchName
		r.SectionName = sectionName
		r.RoleName = roleName

		if roleName != "ADMIN" { // if not admin, checking branch and section
			// checking branch and section
			if _, ok := mapBranches[branchName]; !ok {
				errExcelValidation.Add(rowname+".branch_name", "cabang tidak ditemukan")
			}

			if _, ok := mapSections[sectionName]; !ok {
				errExcelValidation.Add(rowname+".section_name", "section tidak ditemukan")
			}
		}

		if name == "" {
			errExcelValidation.Add(rowname+".name", "nama tidak boleh kosong")
		} else {
			r.Name = name
		}

		if roleName == "" {
			errExcelValidation.Add(rowname+".role_name", "role tidak boleh kosong")
		} else {
			if roleName != "ADMIN" && roleName != "SERVICE ADVISOR" && roleName != "MRA" {
				errExcelValidation.Add(rowname+".role_name", "role tidak valid")
			} else {
				if roleName == "ADMIN" {
					r.RoleId = mapRoles["ADMIN"]
				}
				if roleName == "SERVICE ADVISOR" {
					r.RoleId = mapRoles["SERVICE_ADVISOR"]
				}
				if roleName == "MRA" {
					r.RoleId = mapRoles["TECHNICIAN"]
				}
			}
		}

		if email == "" {
			errExcelValidation.Add(rowname+".email", "email tidak boleh kosong")
		} else {
			if _, ok := mapEmails[email]; ok {
				errExcelValidation.Add(rowname+".email", "email duplikat ditemukan")
			} else {
				mapEmails[email] = true
			}

			err = s.repo.IsEmailExist(ctx, email)
			if err != nil {
				if err.Error() == "email sudah terdaftar" {
					errExcelValidation.Add(rowname+".email", "email sudah terdaftar")
				} else {
					log.Error().Err(err).Msg("service::ImportEmployees - Failed to check email exist")
					return errExcelValidation.SetMessage("Gagal memeriksa email")
				}
			}
		}

		if password == "" {
			errExcelValidation.Add(rowname+".password", "password tidak boleh kosong")
		} else {
			if !isPasswordValid(password) {
				errExcelValidation.Add(rowname+".password", "password minimal 8 karakter dan mengandung karakter spesial")
			} else {
				hashed, err := pkg.HashPassword(password)
				if err != nil {
					log.Error().Err(err).Msg("service::ImportEmployees - Failed to hash password")
					errExcelValidation.Add(rowname+".password", "gagal menghash password")
				} else {
					r.PasswordHashed = hashed
				}
			}
		}

		dataImport = append(dataImport, r)
	}

	if errExcelValidation.HasErrors() {
		// msg := PopulateMsg(errExcelValidation.Errors)
		// errExcelValidation.SetMessage(msg)
		return errExcelValidation
	}

	err = s.repo.ImportEmployees(ctx, dataImport)
	if err != nil {
		return err
	}

	return nil
}

func isXLSXFile(data []byte) error {
	// Check the file header
	if len(data) < 4 {
		return fmt.Errorf("file is too short to determine type")
	}

	if data[0] != 0x50 || data[1] != 0x4b || data[2] != 0x03 || data[3] != 0x04 {
		return fmt.Errorf("file is not a valid xlsx file")
	}

	return nil
}

func isPasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}

	// contains at least one special character
	if !strings.ContainsAny(password, "!@#$%^&*()_+{}|:<>?") {
		return false
	}

	return true
}

func PopulateMsg(errors map[string][]string) string {
	var allErrors []string
	for field, errorList := range errors {
		for _, errMsg := range errorList {
			allErrors = append(allErrors, fmt.Sprintf("Baris %s: %s", extractRowAndField(field), errMsg))
		}
	}
	return strings.Join(allErrors, "\n") // Menggabungkan pesan dengan newline
}

func extractRowAndField(field string) string {
	var row int
	var column string
	fmt.Sscanf(field, "file[%d].%s", &row, &column)
	return fmt.Sprintf("%d", row)
}
