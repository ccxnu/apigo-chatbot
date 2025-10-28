package request

import (
	"time"

	"api-chatbot/domain"
)

// =====================================================
// Base Request (for endpoints with no specific params)
// =====================================================

type AnalyticsBaseRequest struct {
	domain.Base
}

// =====================================================
// Cost Analytics Request
// =====================================================

type GetCostAnalyticsRequest struct {
	domain.Base
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
}

// =====================================================
// Token Usage Request
// =====================================================

type GetTokenUsageRequest struct {
	domain.Base
	Period  string `json:"period" validate:"omitempty,oneof=day week month year all"` // Default: "month"
	GroupBy string `json:"groupBy" validate:"omitempty,oneof=hour day week"`          // Default: "day"
}

// =====================================================
// Active Users Request
// =====================================================

type GetActiveUsersRequest struct {
	domain.Base
	Period string `json:"period" validate:"omitempty,oneof=day week month all"` // Default: "month"
}

// =====================================================
// Conversation Metrics Request
// =====================================================

type GetConversationMetricsRequest struct {
	domain.Base
	Period string `json:"period" validate:"omitempty,oneof=day week month all"` // Default: "month"
}

// =====================================================
// Message Analytics Request
// =====================================================

type GetMessageAnalyticsRequest struct {
	domain.Base
	Period string `json:"period" validate:"omitempty,oneof=day week month all"` // Default: "month"
}

// =====================================================
// Top Queries Request
// =====================================================

type GetTopQueriesRequest struct {
	domain.Base
	Period        string  `json:"period" validate:"omitempty,oneof=day week month all"` // Default: "month"
	Limit         int     `json:"limit" validate:"omitempty,min=1,max=100"`             // Default: 20
	MinSimilarity float64 `json:"minSimilarity" validate:"omitempty,min=0,max=1"`       // Default: 0.5
}

// =====================================================
// Knowledge Usage Request
// =====================================================

type GetKnowledgeUsageRequest struct {
	domain.Base
	Period string `json:"period" validate:"omitempty,oneof=day week month all"` // Default: "month"
}

// =====================================================
// Report Generation Requests
// =====================================================

type GenerateMonthlyReportRequest struct {
	domain.Base
	Year  int `json:"year" validate:"required,min=2020,max=2100"` // Year (e.g., 2025)
	Month int `json:"month" validate:"required,min=1,max=12"`     // Month (1-12)
}

type GenerateCustomReportRequest struct {
	domain.Base
	StartDate string   `json:"startDate" validate:"required"` // Start date (YYYY-MM-DD)
	EndDate   string   `json:"endDate" validate:"required"`   // End date (YYYY-MM-DD)
	Metrics   []string `json:"metrics"`                       // Metrics to include (empty = all)
}
