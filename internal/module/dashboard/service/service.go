package service

import (
	"bytes"
	"codebase-app/internal/module/dashboard/entity"
	"codebase-app/internal/module/dashboard/ports"
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

var _ ports.DashboardService = &DashbaordService{}

type DashbaordService struct {
	repo ports.DashboardRepository
}

func NewDashboardService(repo ports.DashboardRepository) *DashbaordService {
	return &DashbaordService{
		repo: repo,
	}
}

func (s *DashbaordService) GetLeadsTrends(ctx context.Context, request *entity.LeadTrendsRequest) ([]entity.LeadTrendsResponse, error) {
	return s.repo.GetLeadsTrends(ctx, request)
}

func (s *DashbaordService) GetWACSummary(ctx context.Context, request *entity.WACSummaryRequest) (entity.WACSummaryResponse, error) {
	return s.repo.GetWACSummary(ctx, request)
}

func (s *DashbaordService) GetWACSummaryTechnician(ctx context.Context, request *entity.WACSummaryRequest) (entity.TechWACSummaryResponse, error) {
	return s.repo.GetWACSummaryTechnician(ctx, request)
}

func (s *DashbaordService) GetWACLineChart(ctx context.Context, request *entity.GetWACLineChartRequest) (entity.GetWACLineChartResponse, error) {
	return s.repo.GetWACLineChart(ctx, request)
}

func (s *DashbaordService) GetActivities(ctx context.Context, request *entity.GetActivitiesRequest) (entity.GetActivitiesResponse, error) {
	return s.repo.GetActivities(ctx, request)
}

func (s *DashbaordService) GetActivitiesExported(ctx context.Context, request *entity.GetActivitiesRequest) (*entity.ExportMeta, error) {
	var (
		loc, _ = time.LoadLocation(request.Timezone)
		from   = request.FromTime.Format("02/01/2006")
		to     = request.ToTime.Format("02/01/2006")
	)

	resp, err := s.repo.GetActivities(ctx, request)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()

	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		log.Error().Err(err).Any("payload", request).Msg("service::GetActivitiesExported - failed to create style")
		return nil, err
	}

	f.SetCellValue("Sheet1", "A1", "Kalla Toyota")
	f.SetCellValue("Sheet1", "B1", fmt.Sprintf("Cabang %s", request.BranchName))
	f.SetCellValue("Sheet1", "A2", "Periode")
	f.SetCellValue("Sheet1", "B2", fmt.Sprintf("%s - %s", from, to))
	f.SetCellValue("Sheet1", "A3", "Report")
	f.SetCellValue("Sheet1", "B3", "Activity")

	// Header table
	headers := []string{"No", "Tanggal", "Nama Customer", "Cabang", "Penanggung Jawab", "Nomor WhatsApp", "Nomor Plat", "Jenis Mobil", "Status", "Potensi Leads", "Leads", "Total Revenue"}
	for i, h := range headers {
		col := fmt.Sprintf("%c5", 'A'+i)
		f.SetCellValue("Sheet1", col, h)
	}

	for i, a := range resp.Items {
		row := i + 6
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)
		c := a.CreatedAt.In(loc)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), c.Format("2006-01-02 15:04:05"))
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), a.CreatedAt.In(loc).Format("2006-01-02 15:04:05"))
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), a.ClientName)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), a.BranchName)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), a.EmployeeName)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), a.Phone)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), a.VehicleLicenseNumber)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), a.VehicleTypeName)

		s := a.Status
		if s == "wip" {
			s = "Pengerjaan"
		} else if s == "completed" {
			s = "Selesai"
		} else if s == "offered" {
			s = "Penawaran"
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", row), s)
		f.SetCellValue("Sheet1", fmt.Sprintf("J%d", row), a.TotalPotentialLeads)
		f.SetCellValue("Sheet1", fmt.Sprintf("K%d", row), a.TotalLeads)
		f.SetCellValue("Sheet1", fmt.Sprintf("L%d", row), a.TotalRevenue)
	}

	// Terapkan style center alignment pada title
	f.SetCellStyle("Sheet1", "A1", "L1", style)
	f.SetCellStyle("Sheet1", "A2", "L2", style)
	f.SetCellStyle("Sheet1", "A3", "L3", style)
	f.SetCellStyle("Sheet1", "A5", "L5", style)

	f.SetColWidth("Sheet1", "A", "A", 12) // No
	f.SetColWidth("Sheet1", "B", "B", 21) // Tanggal
	f.SetColWidth("Sheet1", "C", "C", 21) // Nama Customer
	f.SetColWidth("Sheet1", "D", "D", 21) // Cabang
	f.SetColWidth("Sheet1", "E", "E", 21) // Penanggung Jawab
	f.SetColWidth("Sheet1", "F", "F", 21) // Nomor WhatsApp
	f.SetColWidth("Sheet1", "G", "G", 15) // Nomor Plat
	f.SetColWidth("Sheet1", "H", "H", 21) // Jenis Mobil
	f.SetColWidth("Sheet1", "I", "I", 15) // Status
	f.SetColWidth("Sheet1", "J", "J", 15) // Potensi Leads
	f.SetColWidth("Sheet1", "K", "K", 15) // Leads
	f.SetColWidth("Sheet1", "L", "L", 15) // Total Revenue

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}

	meta := entity.ExportMeta{
		Buf:      buf,
		Filename: "aktivitas-wac-" + request.BranchName + "-" + time.Now().In(loc).Format("20060102150405") + ".xlsx",
	}

	return &meta, nil
}

func (s *DashbaordService) GetAdminSummary(ctx context.Context, request *entity.GetSummaryPerMonthRequest) (entity.GetSummaryPerMonthResponse, error) {
	return s.repo.GetAdminSummary(ctx, request)
}
