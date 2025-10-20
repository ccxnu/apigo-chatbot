package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

type parameterUseCase struct {
	repo           d.ParameterRepository
	cache          d.ParameterCache
	contextTimeout time.Duration
}

func NewParameterUseCase(repo d.ParameterRepository, cache d.ParameterCache, timeout time.Duration) d.ParameterUseCase {
	return &parameterUseCase{
		repo:           repo,
		cache:          cache,
		contextTimeout: timeout,
	}
}

func (u *parameterUseCase) GetAll(c context.Context) d.Result[[]d.Parameter] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	cached := u.cache.GetAll()
	if len(cached) > 0 {
		return d.Success(cached)
	}

	params, err := u.repo.GetAll(ctx)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch all parameters from database", err, "operation", "GetAll")
		return d.Error[[]d.Parameter](u.cache, "ERR_INTERNAL_DB")
	}

	u.cache.LoadAll(params)

	logger.LogInfo(ctx, "Failed to fetch all parameters from database", params)
	return d.Success(params)
}

func (u *parameterUseCase) GetByCode(c context.Context, code string) d.Result[*d.Parameter] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if param, exists := u.cache.Get(code); exists {
		return d.Success(param)
	}

	param, err := u.repo.GetByCode(ctx, code)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch parameter by code", err,
			"operation", "GetByCode",
			"code", code,
		)
		return d.Error[*d.Parameter](u.cache, "ERR_INTERNAL_DB")
	}

	if param == nil {
		logger.LogWarn(ctx, "Parameter not found",
			"operation", "GetByCode",
			"code", code,
		)
		return d.Error[*d.Parameter](u.cache, "ERR_PARAM_NOT_FOUND")
	}

	u.cache.Set(code, param)

	return d.Success(param)
}

func (u *parameterUseCase) GetByCodes(c context.Context, codes []string) d.Result[[]d.Parameter] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	var result []d.Parameter
	var missingCodes []string

	for _, code := range codes {
		if param, exists := u.cache.Get(code); exists {
			result = append(result, *param)
		} else {
			missingCodes = append(missingCodes, code)
		}
	}

	if len(missingCodes) > 0 {
		params, err := u.repo.GetByCodes(ctx, missingCodes)
		if err != nil {
			return d.Error[[]d.Parameter](u.cache, "ERR_INTERNAL_DB")
		}

		for i := range params {
			u.cache.Set(params[i].Code, &params[i])
			result = append(result, params[i])
		}
	}

	return d.Success(result)
}

func (u *parameterUseCase) GetByName(c context.Context, namePattern string) d.Result[[]d.Parameter] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	params, err := u.repo.GetByName(ctx, namePattern)
	if err != nil {
		return d.Error[[]d.Parameter](u.cache, "ERR_INTERNAL_DB")
	}

	return d.Success(params)
}

func (u *parameterUseCase) GetValue(c context.Context, code string) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if value, exists := u.cache.GetValue(code); exists {
		return d.Success(value)
	}

	param, err := u.repo.GetByCode(ctx, code)
	if err != nil {
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if param == nil {
		return d.Error[d.Data](u.cache, "ERR_PARAM_NOT_FOUND")
	}

	u.cache.Set(code, param)
	data, _ := param.GetDataAsMap()

	return d.Success(data)
}

func (u *parameterUseCase) AddParameter(c context.Context, params d.AddParameterParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.repo.AddParameter(ctx, params)
	if err != nil {
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.cache, result.Code)
	}

	// Reload into cache
	param, _ := u.repo.GetByCode(ctx, params.Code)
	if param != nil {
		u.cache.Set(params.Code, param)
	}

	return d.Success(d.Data{"code": params.Code})
}

func (u *parameterUseCase) UpdParameter(c context.Context, params d.UpdParameterParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.repo.UpdParameter(ctx, params)
	if err != nil {
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.cache, result.Code)
	}

	// Reload into cache
	param, _ := u.repo.GetByCode(ctx, params.Code)
	if param != nil {
		u.cache.Set(params.Code, param)
	} else {
		u.cache.Delete(params.Code)
	}

	return d.Success(d.Data{"code": params.Code})
}

func (u *parameterUseCase) DelParameter(c context.Context, code string) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.repo.DelParameter(ctx, code)
	if err != nil {
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.cache, result.Code)
	}

	// Remove from cache
	u.cache.Delete(code)

	return d.Success(d.Data{"code": code})
}

func (u *parameterUseCase) ReloadCache(c context.Context) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	params, err := u.repo.GetAll(ctx)
	if err != nil {
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	u.cache.LoadAll(params)

	return d.Success(d.Data{"count": len(params)})
}
