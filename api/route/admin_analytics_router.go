package route

import (
	"context"

	"api-chatbot/api/request"
	d "api-chatbot/domain"

	"github.com/danielgtaylor/huma/v2"
)

// Response types for analytics endpoints
type CostAnalyticsResponse struct {
	Body d.Result[*d.CostAnalytics]
}

type TokenUsageResponse struct {
	Body d.Result[[]d.TokenUsage]
}

type ActiveUsersResponse struct {
	Body d.Result[*d.ActiveUsers]
}

type ConversationMetricsResponse struct {
	Body d.Result[*d.ConversationMetrics]
}

type MessageAnalyticsResponse struct {
	Body d.Result[*d.MessageAnalytics]
}

type TopQueriesResponse struct {
	Body d.Result[[]d.TopQuery]
}

type KnowledgeUsageResponse struct {
	Body d.Result[[]d.KnowledgeUsage]
}

type SystemHealthResponse struct {
	Body d.Result[[]d.SystemHealthMetric]
}

type AnalyticsOverviewResponse struct {
	Body d.Result[*d.AnalyticsOverview]
}

// RegisterAnalyticsRoutes registers all analytics endpoints
func RegisterAnalyticsRoutes(humaAPI huma.API, analyticsUC d.AnalyticsUseCase) {

	// =====================================================
	// Dashboard Overview
	// =====================================================
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-analytics-overview",
		Method:      "POST",
		Path:        "/api/v1/admin/analytics/overview",
		Summary:     "Get dashboard overview with key metrics",
		Description: "Returns all key metrics for the main admin dashboard including costs, tokens, users, and conversations.",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.AnalyticsBaseRequest
	}) (*AnalyticsOverviewResponse, error) {
		result := analyticsUC.GetAnalyticsOverview(ctx)
		return &AnalyticsOverviewResponse{Body: result}, nil
	})

	// =====================================================
	// Cost Analytics
	// =====================================================
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-cost-analytics",
		Method:      "POST",
		Path:        "/api/v1/admin/analytics/costs",
		Summary:     "Get cost analytics",
		Description: "Returns detailed cost breakdown including LLM costs, embedding costs, and per-conversation costs. Optionally filter by date range.",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.GetCostAnalyticsRequest
	}) (*CostAnalyticsResponse, error) {
		result := analyticsUC.GetCostAnalytics(ctx, input.Body.StartDate, input.Body.EndDate)
		return &CostAnalyticsResponse{Body: result}, nil
	})

	// =====================================================
	// Token Usage
	// =====================================================
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-token-usage",
		Method:      "POST",
		Path:        "/api/v1/admin/analytics/tokens",
		Summary:     "Get token usage statistics",
		Description: "Returns token usage trends grouped by time period. Supports hourly, daily, and weekly grouping.",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.GetTokenUsageRequest
	}) (*TokenUsageResponse, error) {
		period := input.Body.Period
		if period == "" {
			period = "month"
		}
		groupBy := input.Body.GroupBy
		if groupBy == "" {
			groupBy = "day"
		}
		result := analyticsUC.GetTokenUsage(ctx, period, groupBy)
		return &TokenUsageResponse{Body: result}, nil
	})

	// =====================================================
	// Active Users
	// =====================================================
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-active-users",
		Method:      "POST",
		Path:        "/api/v1/admin/analytics/users",
		Summary:     "Get active user statistics",
		Description: "Returns user activity metrics including DAU/WAU/MAU, new users, returning users, and breakdown by role.",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.GetActiveUsersRequest
	}) (*ActiveUsersResponse, error) {
		period := input.Body.Period
		if period == "" {
			period = "month"
		}
		result := analyticsUC.GetActiveUsers(ctx, period)
		return &ActiveUsersResponse{Body: result}, nil
	})

	// =====================================================
	// Conversation Metrics
	// =====================================================
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-conversation-metrics",
		Method:      "POST",
		Path:        "/api/v1/admin/analytics/conversations",
		Summary:     "Get conversation metrics",
		Description: "Returns conversation statistics including total conversations, admin intervention rate, and conversation quality metrics.",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.GetConversationMetricsRequest
	}) (*ConversationMetricsResponse, error) {
		period := input.Body.Period
		if period == "" {
			period = "month"
		}
		result := analyticsUC.GetConversationMetrics(ctx, period)
		return &ConversationMetricsResponse{Body: result}, nil
	})

	// =====================================================
	// Message Analytics
	// =====================================================
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-message-analytics",
		Method:      "POST",
		Path:        "/api/v1/admin/analytics/messages",
		Summary:     "Get message analytics",
		Description: "Returns message volume statistics including user messages, bot responses, admin messages, and peak hour analysis.",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.GetMessageAnalyticsRequest
	}) (*MessageAnalyticsResponse, error) {
		period := input.Body.Period
		if period == "" {
			period = "month"
		}
		result := analyticsUC.GetMessageAnalytics(ctx, period)
		return &MessageAnalyticsResponse{Body: result}, nil
	})

	// =====================================================
	// Top Queries
	// =====================================================
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-top-queries",
		Method:      "POST",
		Path:        "/api/v1/admin/analytics/top-queries",
		Summary:     "Get most asked questions",
		Description: "Returns the most frequently asked questions with their average similarity scores and answer quality indicators.",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.GetTopQueriesRequest
	}) (*TopQueriesResponse, error) {
		period := input.Body.Period
		if period == "" {
			period = "month"
		}
		limit := input.Body.Limit
		if limit == 0 {
			limit = 20
		}
		minSimilarity := input.Body.MinSimilarity
		if minSimilarity == 0 {
			minSimilarity = 0.5
		}
		result := analyticsUC.GetTopQueries(ctx, period, limit, minSimilarity)
		return &TopQueriesResponse{Body: result}, nil
	})

	// =====================================================
	// Knowledge Base Usage
	// =====================================================
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-knowledge-usage",
		Method:      "POST",
		Path:        "/api/v1/admin/analytics/knowledge-base",
		Summary:     "Get knowledge base usage statistics",
		Description: "Returns chunk usage statistics showing which document chunks are most frequently retrieved and their average similarity scores.",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.GetKnowledgeUsageRequest
	}) (*KnowledgeUsageResponse, error) {
		period := input.Body.Period
		if period == "" {
			period = "month"
		}
		result := analyticsUC.GetKnowledgeUsage(ctx, period)
		return &KnowledgeUsageResponse{Body: result}, nil
	})

	// =====================================================
	// System Health
	// =====================================================
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-system-health",
		Method:      "POST",
		Path:        "/api/v1/admin/analytics/health",
		Summary:     "Get system health metrics",
		Description: "Returns system performance metrics including response times, error counts, and uptime statistics.",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.AnalyticsBaseRequest
	}) (*SystemHealthResponse, error) {
		result := analyticsUC.GetSystemHealth(ctx)
		return &SystemHealthResponse{Body: result}, nil
	})
}
