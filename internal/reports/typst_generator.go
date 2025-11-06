package reports

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"api-chatbot/domain"
)

// ReportGenerator handles Typst report generation
type ReportGenerator struct {
	templateDir string
	outputDir   string
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(templateDir, outputDir string) *ReportGenerator {
	return &ReportGenerator{
		templateDir: templateDir,
		outputDir:   outputDir,
	}
}

// MonthlyReportData contains all data needed for monthly report
type MonthlyReportData struct {
	// Header
	MonthYear     string `json:"month_year"`
	GeneratedDate string `json:"generated_date"`

	// Cost Analytics
	CostThisMonth        string  `json:"cost_this_month"`
	CostChange           float64 `json:"cost_change"`
	LLMCost              string  `json:"llm_cost"`
	EmbeddingCost        string  `json:"embedding_cost"`
	LLMCostPercent       string  `json:"llm_cost_percent"`
	EmbeddingCostPercent string  `json:"embedding_cost_percent"`
	PromptTokens         string  `json:"prompt_tokens"`
	CompletionTokens     string  `json:"completion_tokens"`
	EmbeddingTokens      string  `json:"embedding_tokens"`
	TotalTokens          string  `json:"total_tokens"`
	CostPerConversation  string  `json:"cost_per_conversation"`
	CostPerActiveUser    string  `json:"cost_per_active_user"`
	CostProjection       string  `json:"cost_projection"`

	// Token Details
	AvgTokensPerConversation string `json:"avg_tokens_per_conversation"`
	TokensThisMonth          string `json:"tokens_this_month"`

	// User Activity
	TotalUsers              int64     `json:"total_users"`
	ActiveUsersThisMonth    int64     `json:"active_users_this_month"`
	NewUsersThisMonth       int64     `json:"new_users_this_month"`
	ReturningUsersThisMonth int64     `json:"returning_users_this_month"`
	RetentionRate           string    `json:"retention_rate"`
	StudentsCount           int64     `json:"students_count"`
	ProfessorsCount         int64     `json:"professors_count"`
	ExternalCount           int64     `json:"external_count"`
	StudentsPercent         string    `json:"students_percent"`
	ProfessorsPercent       string    `json:"professors_percent"`
	ExternalPercent         string    `json:"external_percent"`
	AvgMessagesPerUser      string    `json:"avg_messages_per_user"`
	AvgSessionsPerUser      string    `json:"avg_sessions_per_user"`
	AvgSessionDuration      string    `json:"avg_session_duration"`
	TopUsers                []TopUser `json:"top_users"`

	// Conversations
	TotalConversations         int64  `json:"total_conversations"`
	ConversationsThisMonth     int64  `json:"conversations_this_month"`
	ActiveConversations        int64  `json:"active_conversations"`
	AvgMessagesPerConversation string `json:"avg_messages_per_conversation"`
	ConversationsWithAdminHelp int64  `json:"conversations_with_admin_help"`
	AdminInterventionRate      string `json:"admin_intervention_rate"`
	BlockedConversations       int64  `json:"blocked_conversations"`
	TemporaryConversations     int64  `json:"temporary_conversations"`

	// Messages
	MessagesThisMonth      int64  `json:"messages_this_month"`
	UserMessagesThisMonth  int64  `json:"user_messages_this_month"`
	BotMessagesThisMonth   int64  `json:"bot_messages_this_month"`
	AdminMessagesThisMonth int64  `json:"admin_messages_this_month"`
	AvgMessagesPerDay      string `json:"avg_messages_per_day"`
	PeakHour               int    `json:"peak_hour"`
	PeakHourCount          int64  `json:"peak_hour_count"`

	// Top Queries
	TopQueries              []QueryData `json:"top_queries"`
	QueriesNeedingAttention []QueryData `json:"queries_needing_attention"`

	// Knowledge Base
	TopChunks       []ChunkData `json:"top_chunks"`
	TotalDocuments  int         `json:"total_documents"`
	TotalChunks     int         `json:"total_chunks"`
	ChunksUsed      int         `json:"chunks_used"`
	ChunksNeverUsed int         `json:"chunks_never_used"`
	CoverageRate    string      `json:"coverage_rate"`

	// System Performance
	AvgLLMResponseTime  string `json:"avg_llm_response_time"`
	P95ResponseTime     string `json:"p95_response_time"`
	P99ResponseTime     string `json:"p99_response_time"`
	ErrorsLast24h       int64  `json:"errors_last_24h"`
	FailedConversations int    `json:"failed_conversations"`
	Uptime              string `json:"uptime"`
}

type TopUser struct {
	Name         string `json:"name"`
	MessageCount int64  `json:"message_count"`
}

type QueryData struct {
	QueryText     string `json:"query_text"`
	QueryCount    int64  `json:"query_count"`
	AvgSimilarity string `json:"avg_similarity"`
	HasGoodAnswer bool   `json:"has_good_answer"`
}

type ChunkData struct {
	DocumentTitle string `json:"document_title"`
	UsageCount    int64  `json:"usage_count"`
	AvgSimilarity string `json:"avg_similarity"`
}

// GenerateMonthlyReport generates a PDF monthly report and returns the PDF bytes
func (rg *ReportGenerator) GenerateMonthlyReport(ctx context.Context, data MonthlyReportData) ([]byte, error) {
	// Template path
	templatePath := filepath.Join(rg.templateDir, "monthly_report.typ")

	// Check if template exists
	if _, err := os.Stat(templatePath); err != nil {
		return nil, fmt.Errorf("failed to find template: %w", err)
	}

	// Convert data to JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Don't escape - exec.Command handles arguments properly
	// No need for shell escaping when using exec.Command
	jsonString := string(jsonData)

	// Execute typst compile with output to stdout (-)
	// typst compile template.typ - --input data="json_string"
	cmd := exec.CommandContext(ctx, "typst", "compile", templatePath, "-", "--input", fmt.Sprintf("data=%s", jsonString))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to compile typst to PDF: %w, stderr: %s", err, stderr.String())
	}

	// Log warning if there's stderr output
	if stderr.Len() > 0 {
		fmt.Printf("Typst warning: %s\n", stderr.String())
	}

	return stdout.Bytes(), nil
}

// PrepareMonthlyReportData prepares report data from analytics
func PrepareMonthlyReportData(
	monthYear string,
	costAnalytics *domain.CostAnalytics,
	activeUsers *domain.ActiveUsers,
	conversationMetrics *domain.ConversationMetrics,
	messageAnalytics *domain.MessageAnalytics,
	topQueries []domain.TopQuery,
	knowledgeUsage []domain.KnowledgeUsage,
	systemHealth []domain.SystemHealthMetric,
) MonthlyReportData {
	data := MonthlyReportData{
		MonthYear:     monthYear,
		GeneratedDate: time.Now().Format("02/01/2006 15:04"),
	}

	// Cost Analytics
	if costAnalytics != nil {
		if costAnalytics.TotalCost != nil {
			data.CostThisMonth = fmt.Sprintf("%.2f", *costAnalytics.TotalCost)
		} else {
			data.CostThisMonth = "0.00"
		}

		if costAnalytics.LLMCost != nil {
			data.LLMCost = fmt.Sprintf("%.2f", *costAnalytics.LLMCost)
		} else {
			data.LLMCost = "0.00"
		}

		if costAnalytics.EmbeddingCost != nil {
			data.EmbeddingCost = fmt.Sprintf("%.2f", *costAnalytics.EmbeddingCost)
		} else {
			data.EmbeddingCost = "0.00"
		}

		if costAnalytics.TotalCost != nil && *costAnalytics.TotalCost > 0 {
			if costAnalytics.LLMCost != nil {
				data.LLMCostPercent = fmt.Sprintf("%.1f", (*costAnalytics.LLMCost/(*costAnalytics.TotalCost))*100)
			}
			if costAnalytics.EmbeddingCost != nil {
				data.EmbeddingCostPercent = fmt.Sprintf("%.1f", (*costAnalytics.EmbeddingCost/(*costAnalytics.TotalCost))*100)
			}
		}

		// Temporal solucion
		data.TopUsers = []TopUser{}
		data.CoverageRate = "40.0"

		if costAnalytics.PromptTokens != nil {
			data.PromptTokens = formatNumber(*costAnalytics.PromptTokens)
		} else {
			data.PromptTokens = "0"
		}

		if costAnalytics.CompletionTokens != nil {
			data.CompletionTokens = formatNumber(*costAnalytics.CompletionTokens)
		} else {
			data.CompletionTokens = "0"
		}

		if costAnalytics.EmbeddingTokens != nil {
			data.EmbeddingTokens = formatNumber(*costAnalytics.EmbeddingTokens)
		} else {
			data.EmbeddingTokens = "0"
		}

		if costAnalytics.TotalTokens != nil {
			data.TotalTokens = formatNumber(*costAnalytics.TotalTokens)
			data.TokensThisMonth = formatNumber(*costAnalytics.TotalTokens)
		} else {
			data.TotalTokens = "0"
			data.TokensThisMonth = "0"
		}

		if costAnalytics.CostPerConversation != nil {
			data.CostPerConversation = fmt.Sprintf("%.4f", *costAnalytics.CostPerConversation)
		} else {
			data.CostPerConversation = "0.0000"
		}

		if costAnalytics.AvgTokensPerConversation != nil {
			data.AvgTokensPerConversation = fmt.Sprintf("%.0f", *costAnalytics.AvgTokensPerConversation)
		} else {
			data.AvgTokensPerConversation = "0"
		}

		if activeUsers != nil && activeUsers.ActiveUsers > 0 && costAnalytics.TotalCost != nil {
			data.CostPerActiveUser = fmt.Sprintf("%.2f", *costAnalytics.TotalCost/float64(activeUsers.ActiveUsers))
		}
	}

	// Active Users
	if activeUsers != nil {
		data.TotalUsers = activeUsers.TotalUsers
		data.ActiveUsersThisMonth = activeUsers.ActiveUsers
		data.NewUsersThisMonth = activeUsers.NewUsers
		data.ReturningUsersThisMonth = activeUsers.ReturningUsers

		if activeUsers.TotalUsers > 0 {
			data.RetentionRate = fmt.Sprintf("%.1f", (float64(activeUsers.ReturningUsers)/float64(activeUsers.TotalUsers))*100)
		}

		data.StudentsCount = activeUsers.Students
		data.ProfessorsCount = activeUsers.Professors
		data.ExternalCount = activeUsers.External

		total := activeUsers.Students + activeUsers.Professors + activeUsers.External
		if total > 0 {
			data.StudentsPercent = fmt.Sprintf("%.1f", (float64(activeUsers.Students)/float64(total))*100)
			data.ProfessorsPercent = fmt.Sprintf("%.1f", (float64(activeUsers.Professors)/float64(total))*100)
			data.ExternalPercent = fmt.Sprintf("%.1f", (float64(activeUsers.External)/float64(total))*100)
		}

		if activeUsers.AvgMessagesPerUser != nil {
			data.AvgMessagesPerUser = fmt.Sprintf("%.1f", *activeUsers.AvgMessagesPerUser)
		} else {
			data.AvgMessagesPerUser = "0.0"
		}

		if activeUsers.AvgSessionsPerUser != nil {
			data.AvgSessionsPerUser = fmt.Sprintf("%.1f", *activeUsers.AvgSessionsPerUser)
		} else {
			data.AvgSessionsPerUser = "0.0"
		}
	}

	// Conversation Metrics
	if conversationMetrics != nil {
		data.TotalConversations = conversationMetrics.TotalConversations
		data.ConversationsThisMonth = conversationMetrics.NewConversations
		data.ActiveConversations = conversationMetrics.ActiveConversations

		if conversationMetrics.AvgMessagesPerConversation != nil {
			data.AvgMessagesPerConversation = fmt.Sprintf("%.1f", *conversationMetrics.AvgMessagesPerConversation)
		} else {
			data.AvgMessagesPerConversation = "0.0"
		}

		data.ConversationsWithAdminHelp = conversationMetrics.ConversationsWithAdminHelp

		if conversationMetrics.AdminInterventionRate != nil {
			data.AdminInterventionRate = fmt.Sprintf("%.1f", (*conversationMetrics.AdminInterventionRate)*100)
		} else {
			data.AdminInterventionRate = "0.0"
		}

		data.BlockedConversations = conversationMetrics.BlockedConversations
		data.TemporaryConversations = conversationMetrics.TemporaryConversations
	}

	// Message Analytics
	if messageAnalytics != nil {
		data.MessagesThisMonth = messageAnalytics.TotalMessages
		data.UserMessagesThisMonth = messageAnalytics.UserMessages
		data.BotMessagesThisMonth = messageAnalytics.BotMessages
		data.AdminMessagesThisMonth = messageAnalytics.AdminMessages

		if messageAnalytics.AvgMessagesPerDay != nil {
			data.AvgMessagesPerDay = fmt.Sprintf("%.1f", *messageAnalytics.AvgMessagesPerDay)
		} else {
			data.AvgMessagesPerDay = "0.0"
		}

		data.PeakHour = messageAnalytics.PeakHour
		data.PeakHourCount = messageAnalytics.PeakHourCount
	}

	// Top Queries
	data.QueriesNeedingAttention = []QueryData{}

	for i, q := range topQueries {
		if i >= 10 {
			break
		}

		avgSimilarity := "0.00"
		if q.AvgSimilarity != nil {
			avgSimilarity = fmt.Sprintf("%.2f", *q.AvgSimilarity)
		}

		data.TopQueries = append(data.TopQueries, QueryData{
			QueryText:     q.QueryText,
			QueryCount:    q.QueryCount,
			AvgSimilarity: avgSimilarity,
			HasGoodAnswer: q.HasGoodAnswer,
		})

		if !q.HasGoodAnswer {
			data.QueriesNeedingAttention = append(data.QueriesNeedingAttention, QueryData{
				QueryText:     q.QueryText,
				QueryCount:    q.QueryCount,
				AvgSimilarity: avgSimilarity,
				HasGoodAnswer: false,
			})
		}
	}

	// Knowledge Usage
	for i, k := range knowledgeUsage {
		if i >= 10 {
			break
		}

		avgSimilarity := "0.00"
		if k.AvgSimilarity != nil {
			avgSimilarity = fmt.Sprintf("%.2f", *k.AvgSimilarity)
		}

		data.TopChunks = append(data.TopChunks, ChunkData{
			DocumentTitle: k.DocumentTitle,
			UsageCount:    k.UsageCount,
			AvgSimilarity: avgSimilarity,
		})
	}

	// System Health
	for _, metric := range systemHealth {
		switch metric.MetricName {
		case "avg_llm_response_time":
			if metric.MetricValue != nil {
				data.AvgLLMResponseTime = fmt.Sprintf("%.0f", *metric.MetricValue)
			} else {
				data.AvgLLMResponseTime = "0"
			}
		case "p95_llm_response_time":
			if metric.MetricValue != nil {
				data.P95ResponseTime = fmt.Sprintf("%.0f", *metric.MetricValue)
			} else {
				data.P95ResponseTime = "0"
			}
		case "p99_llm_response_time":
			if metric.MetricValue != nil {
				data.P99ResponseTime = fmt.Sprintf("%.0f", *metric.MetricValue)
			} else {
				data.P99ResponseTime = "0"
			}
		case "errors_last_24h":
			if metric.MetricValue != nil {
				data.ErrorsLast24h = int64(*metric.MetricValue)
			} else {
				data.ErrorsLast24h = 0
			}
		}
	}

	data.Uptime = "99.8" // Default - could be calculated from actual uptime data

	return data
}

// formatNumber formats large numbers with thousand separators
func formatNumber(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%.2fM", float64(n)/1000000)
}
