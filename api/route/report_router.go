package route

import (
	"context"
	"fmt"
	"os"

	"api-chatbot/api/request"
	"api-chatbot/domain"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterReportRoutes registers all report-related routes
func RegisterReportRoutes(humaAPI huma.API, reportUseCase domain.ReportUseCase) {
	// =====================================================
	// Generate Monthly Report (PDF)
	// =====================================================

	type GenerateMonthlyReportResponse struct {
		ContentType string `header:"Content-Type"`
		ContentDisposition string `header:"Content-Disposition"`
		Body []byte
	}

	huma.Register(humaAPI, huma.Operation{
		OperationID: "generate-monthly-report",
		Method:      "POST",
		Path:        "/api/v1/admin/reports/generate-monthly",
		Summary:     "Generate and download monthly PDF report",
		Description: "Generates a comprehensive monthly analytics report in PDF format using Typst. The report includes cost analysis, user activity, conversation metrics, top queries, and system health. Returns the PDF file directly for download.",
		Tags:        []string{"Reports"},
	}, func(ctx context.Context, input *struct {
		Body request.GenerateMonthlyReportRequest
	}) (*GenerateMonthlyReportResponse, error) {
		result := reportUseCase.GenerateMonthlyReport(ctx, input.Body.Year, input.Body.Month)

		if !result.Success || result.Data == nil {
			return nil, huma.Error500InternalServerError(result.Info)
		}

		return &GenerateMonthlyReportResponse{
			ContentType: "application/pdf",
			ContentDisposition: fmt.Sprintf("attachment; filename=\"%s\"", result.Data.FileName),
			Body: result.Data.PDFData,
		}, nil
	})

	// =====================================================
	// Download Report (File Download)
	// =====================================================

	type DownloadReportResponse struct {
		Body []byte
	}

	huma.Register(humaAPI, huma.Operation{
		OperationID: "download-report",
		Method:      "GET",
		Path:        "/api/v1/admin/reports/download/{filename}",
		Summary:     "Download generated report",
		Description: "Downloads a previously generated report PDF file",
		Tags:        []string{"Reports"},
	}, func(ctx context.Context, input *struct {
		Filename string `path:"filename" doc:"Report filename"`
	}) (*DownloadReportResponse, error) {
		// Read the file
		filePath := "./reports/" + input.Filename
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, huma.Error404NotFound("Report file not found")
		}

		return &DownloadReportResponse{Body: data}, nil
	})

	// =====================================================
	// List Reports
	// =====================================================

	type ListReportsResponse struct {
		Body struct {
			Success bool     `json:"success"`
			Data    []string `json:"data"`
		}
	}

	huma.Register(humaAPI, huma.Operation{
		OperationID: "list-reports",
		Method:      "POST",
		Path:        "/api/v1/admin/reports/list",
		Summary:     "List all generated reports",
		Description: "Returns a list of all available report files",
		Tags:        []string{"Reports"},
	}, func(ctx context.Context, input *struct {
		Body request.AnalyticsBaseRequest
	}) (*ListReportsResponse, error) {
		// List all PDF files in reports directory
		files, err := os.ReadDir("./reports")
		if err != nil {
			return &ListReportsResponse{
				Body: struct {
					Success bool     `json:"success"`
					Data    []string `json:"data"`
				}{
					Success: false,
					Data:    []string{},
				},
			}, nil
		}

		var reportFiles []string
		for _, file := range files {
			if !file.IsDir() && len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".pdf" {
				reportFiles = append(reportFiles, file.Name())
			}
		}

		return &ListReportsResponse{
			Body: struct {
				Success bool     `json:"success"`
				Data    []string `json:"data"`
			}{
				Success: true,
				Data:    reportFiles,
			},
		}, nil
	})
}
