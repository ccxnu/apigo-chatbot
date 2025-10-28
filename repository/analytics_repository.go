package repository

import (
	"context"
	"fmt"
	"time"

	"api-chatbot/api/dal"
	d "api-chatbot/domain"
)

const (
	// Analytics functions
	fnGetCostAnalytics      = "fn_get_cost_analytics"
	fnGetTokenUsage         = "fn_get_token_usage"
	fnGetActiveUsers        = "fn_get_active_users"
	fnGetConversationMetrics = "fn_get_conversation_metrics"
	fnGetMessageAnalytics   = "fn_get_message_analytics"
	fnGetTopQueries         = "fn_get_top_queries"
	fnGetKnowledgeUsage     = "fn_get_knowledge_usage"
	fnGetSystemHealth       = "fn_get_system_health"
	fnGetAnalyticsOverview  = "fn_get_analytics_overview"
)

type analyticsRepository struct {
	dal *dal.DAL
}

func NewAnalyticsRepository(dal *dal.DAL) d.AnalyticsRepository {
	return &analyticsRepository{
		dal: dal,
	}
}

// GetCostAnalytics retrieves cost analytics for a date range
func (r *analyticsRepository) GetCostAnalytics(ctx context.Context, startDate, endDate *time.Time) (*d.CostAnalytics, error) {
	// Handle nil dates (will use defaults in stored procedure)
	var start, end interface{}
	if startDate != nil {
		start = *startDate
	} else {
		start = nil
	}
	if endDate != nil {
		end = *endDate
	} else {
		end = nil
	}

	results, err := dal.QueryRows[d.CostAnalytics](
		r.dal,
		ctx,
		fnGetCostAnalytics,
		start,
		end,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost analytics via %s: %w", fnGetCostAnalytics, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no cost analytics data returned")
	}

	return &results[0], nil
}

// GetTokenUsage retrieves token usage statistics
func (r *analyticsRepository) GetTokenUsage(ctx context.Context, period, groupBy string) ([]d.TokenUsage, error) {
	results, err := dal.QueryRows[d.TokenUsage](
		r.dal,
		ctx,
		fnGetTokenUsage,
		period,
		groupBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get token usage via %s: %w", fnGetTokenUsage, err)
	}

	return results, nil
}

// GetActiveUsers retrieves active user statistics
func (r *analyticsRepository) GetActiveUsers(ctx context.Context, period string) (*d.ActiveUsers, error) {
	results, err := dal.QueryRows[d.ActiveUsers](
		r.dal,
		ctx,
		fnGetActiveUsers,
		period,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users via %s: %w", fnGetActiveUsers, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no active users data returned")
	}

	return &results[0], nil
}

// GetConversationMetrics retrieves conversation statistics
func (r *analyticsRepository) GetConversationMetrics(ctx context.Context, period string) (*d.ConversationMetrics, error) {
	results, err := dal.QueryRows[d.ConversationMetrics](
		r.dal,
		ctx,
		fnGetConversationMetrics,
		period,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation metrics via %s: %w", fnGetConversationMetrics, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no conversation metrics data returned")
	}

	return &results[0], nil
}

// GetMessageAnalytics retrieves message statistics
func (r *analyticsRepository) GetMessageAnalytics(ctx context.Context, period string) (*d.MessageAnalytics, error) {
	results, err := dal.QueryRows[d.MessageAnalytics](
		r.dal,
		ctx,
		fnGetMessageAnalytics,
		period,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get message analytics via %s: %w", fnGetMessageAnalytics, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no message analytics data returned")
	}

	return &results[0], nil
}

// GetTopQueries retrieves most asked questions
func (r *analyticsRepository) GetTopQueries(ctx context.Context, period string, limit int, minSimilarity float64) ([]d.TopQuery, error) {
	results, err := dal.QueryRows[d.TopQuery](
		r.dal,
		ctx,
		fnGetTopQueries,
		period,
		limit,
		minSimilarity,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get top queries via %s: %w", fnGetTopQueries, err)
	}

	return results, nil
}

// GetKnowledgeUsage retrieves knowledge base usage statistics
func (r *analyticsRepository) GetKnowledgeUsage(ctx context.Context, period string) ([]d.KnowledgeUsage, error) {
	results, err := dal.QueryRows[d.KnowledgeUsage](
		r.dal,
		ctx,
		fnGetKnowledgeUsage,
		period,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge usage via %s: %w", fnGetKnowledgeUsage, err)
	}

	return results, nil
}

// GetSystemHealth retrieves system health metrics
func (r *analyticsRepository) GetSystemHealth(ctx context.Context) ([]d.SystemHealthMetric, error) {
	results, err := dal.QueryRows[d.SystemHealthMetric](
		r.dal,
		ctx,
		fnGetSystemHealth,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get system health via %s: %w", fnGetSystemHealth, err)
	}

	return results, nil
}

// GetAnalyticsOverview retrieves dashboard overview
func (r *analyticsRepository) GetAnalyticsOverview(ctx context.Context) (*d.AnalyticsOverview, error) {
	results, err := dal.QueryRows[d.AnalyticsOverview](
		r.dal,
		ctx,
		fnGetAnalyticsOverview,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics overview via %s: %w", fnGetAnalyticsOverview, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no analytics overview data returned")
	}

	return &results[0], nil
}
