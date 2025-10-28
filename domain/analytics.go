package domain

import (
	"context"
	"time"
)

// =====================================================
// Cost Analytics
// =====================================================

type CostAnalytics struct {
	PeriodStart              time.Time `json:"periodStart" db:"period_start"`
	PeriodEnd                time.Time `json:"periodEnd" db:"period_end"`
	TotalCost                float64   `json:"totalCost" db:"total_cost"`
	LLMCost                  float64   `json:"llmCost" db:"llm_cost"`
	EmbeddingCost            float64   `json:"embeddingCost" db:"embedding_cost"`
	PromptTokens             int64     `json:"promptTokens" db:"prompt_tokens"`
	CompletionTokens         int64     `json:"completionTokens" db:"completion_tokens"`
	TotalTokens              int64     `json:"totalTokens" db:"total_tokens"`
	EmbeddingTokens          int64     `json:"embeddingTokens" db:"embedding_tokens"`
	ConversationCount        int64     `json:"conversationCount" db:"conversation_count"`
	CostPerConversation      float64   `json:"costPerConversation" db:"cost_per_conversation"`
	AvgTokensPerConversation float64   `json:"avgTokensPerConversation" db:"avg_tokens_per_conversation"`
}

// =====================================================
// Token Usage
// =====================================================

type TokenUsage struct {
	PeriodLabel              string    `json:"periodLabel" db:"period_label"`
	PeriodStart              time.Time `json:"periodStart" db:"period_start"`
	PeriodEnd                time.Time `json:"periodEnd" db:"period_end"`
	PromptTokens             int64     `json:"promptTokens" db:"prompt_tokens"`
	CompletionTokens         int64     `json:"completionTokens" db:"completion_tokens"`
	TotalTokens              int64     `json:"totalTokens" db:"total_tokens"`
	MessageCount             int64     `json:"messageCount" db:"message_count"`
	ConversationCount        int64     `json:"conversationCount" db:"conversation_count"`
	AvgTokensPerMessage      float64   `json:"avgTokensPerMessage" db:"avg_tokens_per_message"`
}

// =====================================================
// Active Users
// =====================================================

type ActiveUsers struct {
	Period                 string  `json:"period" db:"period"`
	TotalUsers             int64   `json:"totalUsers" db:"total_users"`
	ActiveUsers            int64   `json:"activeUsers" db:"active_users"`
	NewUsers               int64   `json:"newUsers" db:"new_users"`
	ReturningUsers         int64   `json:"returningUsers" db:"returning_users"`
	Students               int64   `json:"students" db:"students"`
	Professors             int64   `json:"professors" db:"professors"`
	External               int64   `json:"external" db:"external"`
	AvgMessagesPerUser     float64 `json:"avgMessagesPerUser" db:"avg_messages_per_user"`
	AvgSessionsPerUser     float64 `json:"avgSessionsPerUser" db:"avg_sessions_per_user"`
}

// =====================================================
// Conversation Metrics
// =====================================================

type ConversationMetrics struct {
	Period                      string  `json:"period" db:"period"`
	TotalConversations          int64   `json:"totalConversations" db:"total_conversations"`
	ActiveConversations         int64   `json:"activeConversations" db:"active_conversations"`
	NewConversations            int64   `json:"newConversations" db:"new_conversations"`
	AvgMessagesPerConversation  float64 `json:"avgMessagesPerConversation" db:"avg_messages_per_conversation"`
	ConversationsWithAdminHelp  int64   `json:"conversationsWithAdminHelp" db:"conversations_with_admin_help"`
	AdminInterventionRate       float64 `json:"adminInterventionRate" db:"admin_intervention_rate"`
	BlockedConversations        int64   `json:"blockedConversations" db:"blocked_conversations"`
	TemporaryConversations      int64   `json:"temporaryConversations" db:"temporary_conversations"`
}

// =====================================================
// Message Analytics
// =====================================================

type MessageAnalytics struct {
	Period              string  `json:"period" db:"period"`
	TotalMessages       int64   `json:"totalMessages" db:"total_messages"`
	UserMessages        int64   `json:"userMessages" db:"user_messages"`
	BotMessages         int64   `json:"botMessages" db:"bot_messages"`
	AdminMessages       int64   `json:"adminMessages" db:"admin_messages"`
	AvgMessagesPerDay   float64 `json:"avgMessagesPerDay" db:"avg_messages_per_day"`
	PeakHour            int     `json:"peakHour" db:"peak_hour"`
	PeakHourCount       int64   `json:"peakHourCount" db:"peak_hour_count"`
}

// =====================================================
// Top Queries
// =====================================================

type TopQuery struct {
	QueryText      string    `json:"queryText" db:"query_text"`
	QueryCount     int64     `json:"queryCount" db:"query_count"`
	AvgSimilarity  float64   `json:"avgSimilarity" db:"avg_similarity"`
	LastAsked      time.Time `json:"lastAsked" db:"last_asked"`
	HasGoodAnswer  bool      `json:"hasGoodAnswer" db:"has_good_answer"`
}

// =====================================================
// Knowledge Base Usage
// =====================================================

type KnowledgeUsage struct {
	ChunkID       int       `json:"chunkId" db:"chunk_id"`
	DocumentTitle string    `json:"documentTitle" db:"document_title"`
	UsageCount    int64     `json:"usageCount" db:"usage_count"`
	AvgSimilarity float64   `json:"avgSimilarity" db:"avg_similarity"`
	LastUsed      time.Time `json:"lastUsed" db:"last_used"`
}

// =====================================================
// System Health
// =====================================================

type SystemHealthMetric struct {
	MetricName  string  `json:"metricName" db:"metric_name"`
	MetricValue float64 `json:"metricValue" db:"metric_value"`
	MetricUnit  string  `json:"metricUnit" db:"metric_unit"`
}

// =====================================================
// Dashboard Overview
// =====================================================

type AnalyticsOverview struct {
	CostThisMonth          float64   `json:"costThisMonth" db:"cost_this_month"`
	TokensThisMonth        int64     `json:"tokensThisMonth" db:"tokens_this_month"`
	ActiveUsersToday       int64     `json:"activeUsersToday" db:"active_users_today"`
	ConversationsThisMonth int64     `json:"conversationsThisMonth" db:"conversations_this_month"`
	MessagesToday          int64     `json:"messagesToday" db:"messages_today"`
	AvgResponseTimeMs      float64   `json:"avgResponseTimeMs" db:"avg_response_time_ms"`
	AdminInterventionRate  float64   `json:"adminInterventionRate" db:"admin_intervention_rate"`
	LastUpdated            time.Time `json:"lastUpdated" db:"last_updated"`
}

// =====================================================
// Repository Interface
// =====================================================

type AnalyticsRepository interface {
	// Cost analytics
	GetCostAnalytics(ctx context.Context, startDate, endDate *time.Time) (*CostAnalytics, error)

	// Token usage
	GetTokenUsage(ctx context.Context, period, groupBy string) ([]TokenUsage, error)

	// Active users
	GetActiveUsers(ctx context.Context, period string) (*ActiveUsers, error)

	// Conversation metrics
	GetConversationMetrics(ctx context.Context, period string) (*ConversationMetrics, error)

	// Message analytics
	GetMessageAnalytics(ctx context.Context, period string) (*MessageAnalytics, error)

	// Top queries
	GetTopQueries(ctx context.Context, period string, limit int, minSimilarity float64) ([]TopQuery, error)

	// Knowledge base usage
	GetKnowledgeUsage(ctx context.Context, period string) ([]KnowledgeUsage, error)

	// System health
	GetSystemHealth(ctx context.Context) ([]SystemHealthMetric, error)

	// Dashboard overview
	GetAnalyticsOverview(ctx context.Context) (*AnalyticsOverview, error)
}

// =====================================================
// Use Case Interface
// =====================================================

type AnalyticsUseCase interface {
	// Cost analytics
	GetCostAnalytics(ctx context.Context, startDate, endDate *time.Time) Result[*CostAnalytics]

	// Token usage
	GetTokenUsage(ctx context.Context, period, groupBy string) Result[[]TokenUsage]

	// Active users
	GetActiveUsers(ctx context.Context, period string) Result[*ActiveUsers]

	// Conversation metrics
	GetConversationMetrics(ctx context.Context, period string) Result[*ConversationMetrics]

	// Message analytics
	GetMessageAnalytics(ctx context.Context, period string) Result[*MessageAnalytics]

	// Top queries
	GetTopQueries(ctx context.Context, period string, limit int, minSimilarity float64) Result[[]TopQuery]

	// Knowledge base usage
	GetKnowledgeUsage(ctx context.Context, period string) Result[[]KnowledgeUsage]

	// System health
	GetSystemHealth(ctx context.Context) Result[[]SystemHealthMetric]

	// Dashboard overview
	GetAnalyticsOverview(ctx context.Context) Result[*AnalyticsOverview]
}

// =====================================================
// Report Types
// =====================================================

type GeneratedReport struct {
	FilePath      string    `json:"filePath"`
	FileName      string    `json:"fileName"`
	ReportType    string    `json:"reportType"` // "monthly", "quarterly", "custom"
	Period        string    `json:"period"`
	GeneratedAt   time.Time `json:"generatedAt"`
	FileSizeBytes int64     `json:"fileSizeBytes"`
}

// =====================================================
// Report Use Case Interface
// =====================================================

type ReportUseCase interface {
	// Generate monthly report
	GenerateMonthlyReport(ctx context.Context, year int, month int) Result[*GeneratedReport]

	// Generate custom report
	GenerateCustomReport(ctx context.Context, startDate time.Time, endDate time.Time, metrics []string) Result[*GeneratedReport]
}
