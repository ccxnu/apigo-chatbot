package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

type analyticsUseCase struct {
	repo           d.AnalyticsRepository
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

func NewAnalyticsUseCase(
	repo d.AnalyticsRepository,
	paramCache d.ParameterCache,
	timeout time.Duration,
) d.AnalyticsUseCase {
	return &analyticsUseCase{
		repo:           repo,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

// GetCostAnalytics retrieves cost analytics for a date range
func (uc *analyticsUseCase) GetCostAnalytics(
	c context.Context,
	startDate, endDate *time.Time,
) d.Result[*d.CostAnalytics] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	analytics, err := uc.repo.GetCostAnalytics(ctx, startDate, endDate)
	if err != nil {
		logger.LogError(ctx, "Failed to get cost analytics", err,
			"operation", "GetCostAnalytics",
		)
		return d.Error[*d.CostAnalytics](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(analytics)
}

// GetTokenUsage retrieves token usage statistics
func (uc *analyticsUseCase) GetTokenUsage(
	c context.Context,
	period, groupBy string,
) d.Result[[]d.TokenUsage] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Validate period
	validPeriods := map[string]bool{"day": true, "week": true, "month": true, "year": true, "all": true}
	if !validPeriods[period] {
		period = "month"
	}

	// Validate groupBy
	validGroupBy := map[string]bool{"hour": true, "day": true, "week": true}
	if !validGroupBy[groupBy] {
		groupBy = "day"
	}

	usage, err := uc.repo.GetTokenUsage(ctx, period, groupBy)
	if err != nil {
		logger.LogError(ctx, "Failed to get token usage", err,
			"operation", "GetTokenUsage",
			"period", period,
			"groupBy", groupBy,
		)
		return d.Error[[]d.TokenUsage](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(usage)
}

// GetActiveUsers retrieves active user statistics
func (uc *analyticsUseCase) GetActiveUsers(
	c context.Context,
	period string,
) d.Result[*d.ActiveUsers] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Validate period
	validPeriods := map[string]bool{"day": true, "week": true, "month": true, "all": true}
	if !validPeriods[period] {
		period = "month"
	}

	users, err := uc.repo.GetActiveUsers(ctx, period)
	if err != nil {
		logger.LogError(ctx, "Failed to get active users", err,
			"operation", "GetActiveUsers",
			"period", period,
		)
		return d.Error[*d.ActiveUsers](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(users)
}

// GetConversationMetrics retrieves conversation statistics
func (uc *analyticsUseCase) GetConversationMetrics(
	c context.Context,
	period string,
) d.Result[*d.ConversationMetrics] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Validate period
	validPeriods := map[string]bool{"day": true, "week": true, "month": true, "all": true}
	if !validPeriods[period] {
		period = "month"
	}

	metrics, err := uc.repo.GetConversationMetrics(ctx, period)
	if err != nil {
		logger.LogError(ctx, "Failed to get conversation metrics", err,
			"operation", "GetConversationMetrics",
			"period", period,
		)
		return d.Error[*d.ConversationMetrics](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(metrics)
}

// GetMessageAnalytics retrieves message statistics
func (uc *analyticsUseCase) GetMessageAnalytics(
	c context.Context,
	period string,
) d.Result[*d.MessageAnalytics] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Validate period
	validPeriods := map[string]bool{"day": true, "week": true, "month": true, "all": true}
	if !validPeriods[period] {
		period = "month"
	}

	analytics, err := uc.repo.GetMessageAnalytics(ctx, period)
	if err != nil {
		logger.LogError(ctx, "Failed to get message analytics", err,
			"operation", "GetMessageAnalytics",
			"period", period,
		)
		return d.Error[*d.MessageAnalytics](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(analytics)
}

// GetTopQueries retrieves most asked questions
func (uc *analyticsUseCase) GetTopQueries(
	c context.Context,
	period string,
	limit int,
	minSimilarity float64,
) d.Result[[]d.TopQuery] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Validate period
	validPeriods := map[string]bool{"day": true, "week": true, "month": true, "all": true}
	if !validPeriods[period] {
		period = "month"
	}

	// Validate limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// Validate similarity
	if minSimilarity < 0 || minSimilarity > 1 {
		minSimilarity = 0.5
	}

	queries, err := uc.repo.GetTopQueries(ctx, period, limit, minSimilarity)
	if err != nil {
		logger.LogError(ctx, "Failed to get top queries", err,
			"operation", "GetTopQueries",
			"period", period,
			"limit", limit,
		)
		return d.Error[[]d.TopQuery](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(queries)
}

// GetKnowledgeUsage retrieves knowledge base usage statistics
func (uc *analyticsUseCase) GetKnowledgeUsage(
	c context.Context,
	period string,
) d.Result[[]d.KnowledgeUsage] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Validate period
	validPeriods := map[string]bool{"day": true, "week": true, "month": true, "all": true}
	if !validPeriods[period] {
		period = "month"
	}

	usage, err := uc.repo.GetKnowledgeUsage(ctx, period)
	if err != nil {
		logger.LogError(ctx, "Failed to get knowledge usage", err,
			"operation", "GetKnowledgeUsage",
			"period", period,
		)
		return d.Error[[]d.KnowledgeUsage](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(usage)
}

// GetSystemHealth retrieves system health metrics
func (uc *analyticsUseCase) GetSystemHealth(
	c context.Context,
) d.Result[[]d.SystemHealthMetric] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	health, err := uc.repo.GetSystemHealth(ctx)
	if err != nil {
		logger.LogError(ctx, "Failed to get system health", err,
			"operation", "GetSystemHealth",
		)
		return d.Error[[]d.SystemHealthMetric](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(health)
}

// GetAnalyticsOverview retrieves dashboard overview
func (uc *analyticsUseCase) GetAnalyticsOverview(
	c context.Context,
) d.Result[*d.AnalyticsOverview] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	overview, err := uc.repo.GetAnalyticsOverview(ctx)
	if err != nil {
		logger.LogError(ctx, "Failed to get analytics overview", err,
			"operation", "GetAnalyticsOverview",
		)
		return d.Error[*d.AnalyticsOverview](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(overview)
}
