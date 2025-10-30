package repository

import (
	"context"
	"encoding/json"

	"api-chatbot/api/dal"
	d "api-chatbot/domain"
)

type apiKeyRepository struct {
	dal *dal.DAL
}

func NewAPIKeyRepository(dalInstance *dal.DAL) d.APIKeyRepository {
	return &apiKeyRepository{dal: dalInstance}
}

func (r *apiKeyRepository) Create(ctx context.Context, params d.CreateAPIKeyParams) (*d.CreateAPIKeyResult, error) {
	// Convert slices to JSONB
	allowedIPsJSON, err := json.Marshal(params.AllowedIPs)
	if err != nil {
		return nil, err
	}

	permissionsJSON, err := json.Marshal(params.Permissions)
	if err != nil {
		return nil, err
	}

	claimsJSON, err := json.Marshal(params.Claims)
	if err != nil {
		return nil, err
	}

	return dal.ExecProc[d.CreateAPIKeyResult](
		r.dal,
		ctx,
		"sp_create_api_key",
		params.Name,
		params.Value,
		params.Type,
		claimsJSON,
		params.RateLimit,
		allowedIPsJSON,
		permissionsJSON,
		params.ExpiresAt,
		params.CreatedBy,
	)
}

func (r *apiKeyRepository) GetByValue(ctx context.Context, keyValue string) (*d.APIKey, error) {
	return dal.QueryRow[d.APIKey](r.dal, ctx, "fn_get_api_key_by_value", keyValue)
}

func (r *apiKeyRepository) GetByID(ctx context.Context, keyID int) (*d.APIKey, error) {
	return dal.QueryRow[d.APIKey](r.dal, ctx, "fn_get_api_key_by_id", keyID)
}

func (r *apiKeyRepository) GetAll(ctx context.Context) ([]d.APIKey, error) {
	return dal.QueryRows[d.APIKey](r.dal, ctx, "fn_get_all_api_keys")
}

func (r *apiKeyRepository) Update(ctx context.Context, params d.UpdateAPIKeyParams) (*d.UpdateAPIKeyResult, error) {
	var allowedIPsJSON, permissionsJSON interface{}

	if params.AllowedIPs != nil {
		data, err := json.Marshal(*params.AllowedIPs)
		if err != nil {
			return nil, err
		}
		allowedIPsJSON = data
	}

	if params.Permissions != nil {
		data, err := json.Marshal(*params.Permissions)
		if err != nil {
			return nil, err
		}
		permissionsJSON = data
	}

	return dal.ExecProc[d.UpdateAPIKeyResult](
		r.dal,
		ctx,
		"sp_update_api_key",
		params.KeyID,
		params.Name,
		params.RateLimit,
		allowedIPsJSON,
		permissionsJSON,
		params.IsActive,
		params.ExpiresAt,
	)
}

func (r *apiKeyRepository) UpdateLastUsed(ctx context.Context, keyID int) (*d.UpdateAPIKeyLastUsedResult, error) {
	return dal.ExecProc[d.UpdateAPIKeyLastUsedResult](
		r.dal,
		ctx,
		"sp_update_api_key_last_used",
		keyID,
	)
}

func (r *apiKeyRepository) Delete(ctx context.Context, keyID int) (*d.DeleteAPIKeyResult, error) {
	return dal.ExecProc[d.DeleteAPIKeyResult](
		r.dal,
		ctx,
		"sp_delete_api_key",
		keyID,
	)
}
