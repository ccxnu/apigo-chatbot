package repository

import (
	"context"

	"api-chatbot/api/dal"
	d "api-chatbot/domain"
	"time"
)

type apiUsageRepository struct {
	dal *dal.DAL
}

func NewAPIUsageRepository(dalInstance *dal.DAL) d.APIUsageRepository {
	return &apiUsageRepository{dal: dalInstance}
}

func (r *apiUsageRepository) Track(ctx context.Context, params d.TrackAPIUsageParams) (*d.TrackAPIUsageResult, error) {
	return dal.ExecProc[d.TrackAPIUsageResult](
		r.dal,
		ctx,
		"sp_track_api_usage",
		params.APIKeyID,
		params.Endpoint,
		params.Method,
		params.StatusCode,
		params.TokensUsed,
		params.RequestTimeMs,
		params.IPAddress,
		params.UserAgent,
		params.ErrorMessage,
	)
}

func (r *apiUsageRepository) GetStats(ctx context.Context, keyID int, from, to *time.Time) (*d.APIUsageStats, error) {
	return dal.QueryRow[d.APIUsageStats](r.dal, ctx, "fn_get_api_usage_stats", keyID, from, to)
}
