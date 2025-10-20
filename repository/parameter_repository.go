package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"api-chatbot/api/dal"
	"api-chatbot/domain"
)

const (
	// Functions (Read-only)
	fnGetAllParameters     = "fn_get_all_parameters"
	fnGetParameterByCode   = "fn_get_parameter_by_code"
	fnGetParametersByCodes = "fn_get_parameters_by_codes"
	fnGetParametersByName  = "fn_get_parameters_by_name"
	fnGetParameterValue    = "fn_get_parameter_value"
	// Stored Procedures (Writes)
	spCreateParameter = "sp_create_parameter"
	spUpdateParameter = "sp_update_parameter"
	spDeleteParameter = "sp_delete_parameter"
)

type parameterRepository struct {
	dal *dal.DAL
}

func NewParameterRepository(dal *dal.DAL) domain.ParameterRepository {
	return &parameterRepository{
		dal: dal,
	}
}

// GetAll retrieves all active parameters
func (r *parameterRepository) GetAll(ctx context.Context) ([]domain.Parameter, error) {
	params, err := dal.QueryRows[domain.Parameter](r.dal, ctx, fnGetAllParameters)
	if err != nil {
		return nil, fmt.Errorf("failed to get all parameters via %s: %w", fnGetAllParameters, err)
	}
	return params, nil
}

// GetByCode retrieves a single parameter by code
func (r *parameterRepository) GetByCode(ctx context.Context, code string) (*domain.Parameter, error) {
	params, err := dal.QueryRows[domain.Parameter](r.dal, ctx, fnGetParameterByCode, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameter by code via %s: %w", fnGetParameterByCode, err)
	}

	if len(params) == 0 {
		return nil, nil
	}

	return &params[0], nil
}

// GetByCodes retrieves multiple parameters by their codes (bulk operation)
func (r *parameterRepository) GetByCodes(ctx context.Context, codes []string) ([]domain.Parameter, error) {
	params, err := dal.QueryRows[domain.Parameter](r.dal, ctx, fnGetParametersByCodes, codes)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters by codes via %s: %w", fnGetParametersByCodes, err)
	}
	return params, nil
}

// GetByName retrieves parameters matching a name pattern
func (r *parameterRepository) GetByName(ctx context.Context, namePattern string) ([]domain.Parameter, error) {
	params, err := dal.QueryRows[domain.Parameter](r.dal, ctx, fnGetParametersByName, namePattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters by name via %s: %w", fnGetParametersByName, err)
	}
	return params, nil
}

type ValueResult struct {
	Value json.RawMessage `db:"fn_get_parameter_value"`
}

// GetValue retrieves only the data value of a parameter
func (r *parameterRepository) GetValue(ctx context.Context, code string) (map[string]any, error) {
	rows, err := dal.QueryRows[ValueResult](r.dal, ctx, fnGetParameterValue, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameter value via %s: %w", fnGetParameterValue, err)
	}

	if len(rows) == 0 {
		return make(map[string]any), nil
	}

	var data map[string]any
	if err := json.Unmarshal(rows[0].Value, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameter value: %w", err)
	}

	return data, nil
}

// AddParameter creates a new parameter
func (r *parameterRepository) AddParameter(ctx context.Context, params domain.AddParameterParams) (*domain.AddParameterResult, error) {
	result, err := dal.ExecProc[domain.AddParameterResult](
		r.dal,
		ctx,
		spCreateParameter,
		params.Name,
		params.Code,
		params.Data,
		params.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spCreateParameter, err)
	}

	return result, nil
}

// UpdParameter updates an existing parameter
func (r *parameterRepository) UpdParameter(ctx context.Context, params domain.UpdParameterParams) (*domain.UpdParameterResult, error) {
	result, err := dal.ExecProc[domain.UpdParameterResult](
		r.dal,
		ctx,
		spUpdateParameter,
		params.Code,
		params.Name,
		params.Data,
		params.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spUpdateParameter, err)
	}

	return result, nil
}

// DelParameter soft deletes a parameter
func (r *parameterRepository) DelParameter(ctx context.Context, code string) (*domain.DelParameterResult, error) {
	result, err := dal.ExecProc[domain.DelParameterResult](
		r.dal,
		ctx,
		spDeleteParameter,
		code,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spDeleteParameter, err)
	}

	return result, nil
}
