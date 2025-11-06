package usecase

import (
	"context"
	"fmt"
	"time"

	"api-chatbot/domain"
	"api-chatbot/internal/reports"
)

type reportUseCase struct {
	analyticsRepo   domain.AnalyticsRepository
	reportGenerator *reports.ReportGenerator
	timeout         time.Duration
}

func NewReportUseCase(
	analyticsRepo domain.AnalyticsRepository,
	reportGenerator *reports.ReportGenerator,
	timeout time.Duration,
) domain.ReportUseCase {
	return &reportUseCase{
		analyticsRepo:   analyticsRepo,
		reportGenerator: reportGenerator,
		timeout:         timeout,
	}
}

func (u *reportUseCase) GenerateMonthlyReport(ctx context.Context, year int, month int) domain.Result[*domain.GeneratedReport] {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	// Calculate date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	monthYear := startDate.Format("January 2006")
	if startDate.Month() >= 1 && startDate.Month() <= 12 {
		spanishMonths := []string{
			"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
			"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
		}
		monthYear = fmt.Sprintf("%s %d", spanishMonths[startDate.Month()-1], year)
	}

	// Gather all analytics data
	costAnalytics, err := u.analyticsRepo.GetCostAnalytics(ctx, &startDate, &endDate)
	if err != nil {
		return domain.Result[*domain.GeneratedReport]{
			Success: false,
			Info:    "Error al obtener analíticas de costos",
			Code:    "ERR_COST_ANALYTICS",
		}
	}

	activeUsers, err := u.analyticsRepo.GetActiveUsers(ctx, "month")
	if err != nil {
		return domain.Result[*domain.GeneratedReport]{
			Success: false,
			Info:    "Error al obtener usuarios activos",
			Code:    "ERR_ACTIVE_USERS",
		}
	}

	conversationMetrics, err := u.analyticsRepo.GetConversationMetrics(ctx, "month")
	if err != nil {
		return domain.Result[*domain.GeneratedReport]{
			Success: false,
			Info:    "Error al obtener métricas de conversaciones",
			Code:    "ERR_CONVERSATION_METRICS",
		}
	}

	messageAnalytics, err := u.analyticsRepo.GetMessageAnalytics(ctx, "month")
	if err != nil {
		return domain.Result[*domain.GeneratedReport]{
			Success: false,
			Info:    "Error al obtener analíticas de mensajes",
			Code:    "ERR_MESSAGE_ANALYTICS",
		}
	}

	topQueries, err := u.analyticsRepo.GetTopQueries(ctx, "month", 20, 0.5)
	if err != nil {
		return domain.Result[*domain.GeneratedReport]{
			Success: false,
			Info:    "Error al obtener consultas principales",
			Code:    "ERR_TOP_QUERIES",
		}
	}

	knowledgeUsage, err := u.analyticsRepo.GetKnowledgeUsage(ctx, "month")
	if err != nil {
		return domain.Result[*domain.GeneratedReport]{
			Success: false,
			Info:    "Error al obtener uso de base de conocimientos",
			Code:    "ERR_KNOWLEDGE_USAGE",
		}
	}

	systemHealth, err := u.analyticsRepo.GetSystemHealth(ctx)
	if err != nil {
		return domain.Result[*domain.GeneratedReport]{
			Success: false,
			Info:    "Error al obtener salud del sistema",
			Code:    "ERR_SYSTEM_HEALTH",
		}
	}

	// Prepare report data
	reportData := reports.PrepareMonthlyReportData(
		monthYear,
		costAnalytics,
		activeUsers,
		conversationMetrics,
		messageAnalytics,
		topQueries,
		knowledgeUsage,
		systemHealth,
	)

	// Generate PDF report
	pdfBytes, err := u.reportGenerator.GenerateMonthlyReport(ctx, reportData)
	if err != nil {
		return domain.Result[*domain.GeneratedReport]{
			Success: false,
			Info:    fmt.Sprintf("Error al generar reporte: %v", err),
			Code:    "ERR_GENERATE_REPORT",
		}
	}

	generatedReport := &domain.GeneratedReport{
		FileName:      fmt.Sprintf("reporte_mensual_%s.pdf", time.Now().Format("2006-01-02")),
		ReportType:    "monthly",
		Period:        monthYear,
		GeneratedAt:   time.Now(),
		FileSizeBytes: int64(len(pdfBytes)),
		PDFData:       pdfBytes,
	}

	// return domain.Success(generatedReport)
	return domain.Result[*domain.GeneratedReport]{
		Success: true,
		Info:    "Reporte generado exitosamente",
		Code:    "SUCCESS",
		Data:    generatedReport,
	}
}

func (u *reportUseCase) GenerateCustomReport(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
	metrics []string,
) domain.Result[*domain.GeneratedReport] {
	// TODO: Implement custom report generation
	return domain.Result[*domain.GeneratedReport]{
		Success: false,
		Info:    "Reportes personalizados no implementados aún",
		Code:    "ERR_NOT_IMPLEMENTED",
	}
}
