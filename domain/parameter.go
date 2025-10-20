package domain

import (
	"context"
	"encoding/json"
	"time"

	"api-chatbot/api/dal"
)

// Parameter Domain
type Parameter struct {
	ID          int             `json:"id" db:"prm_id"`
	Name        string          `json:"name" db:"prm_name"`
	Code        string          `json:"code" db:"prm_code"`
	Data        json.RawMessage `json:"data" db:"prm_data"`
	Description string          `json:"description" db:"prm_description"`
	Active      bool            `json:"active" db:"prm_active"`
	CreatedAt   time.Time       `json:"createdAt" db:"prm_created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" db:"prm_updated_at"`
}

func (p *Parameter) GetDataAsMap() (Data, error) {
	var data Data
	if err := json.Unmarshal(p.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// =====================================================
// Parameter Repository Params & Results
// =====================================================

type AddParameterParams struct {
	Name        string
	Code        string
	Data        json.RawMessage
	Description string
}

type AddParameterResult struct {
	dal.DbResult
}

type UpdParameterParams struct {
	Code        string
	Name        string
	Data        json.RawMessage
	Description string
}

type UpdParameterResult struct {
	dal.DbResult
}

type DelParameterResult struct {
	dal.DbResult
}

// Repository Interface
type ParameterRepository interface {
	GetAll(ctx context.Context) ([]Parameter, error)
	GetByCode(ctx context.Context, code string) (*Parameter, error)
	GetByCodes(ctx context.Context, codes []string) ([]Parameter, error)
	GetByName(ctx context.Context, namePattern string) ([]Parameter, error)
	AddParameter(ctx context.Context, params AddParameterParams) (*AddParameterResult, error)
	UpdParameter(ctx context.Context, params UpdParameterParams) (*UpdParameterResult, error)
	DelParameter(ctx context.Context, code string) (*DelParameterResult, error)
}

// UseCase Interface
type ParameterUseCase interface {
	GetAll(ctx context.Context) Result[[]Parameter]
	GetByCode(ctx context.Context, code string) Result[*Parameter]
	GetByCodes(ctx context.Context, codes []string) Result[[]Parameter]
	GetByName(ctx context.Context, namePattern string) Result[[]Parameter]
	GetValue(ctx context.Context, code string) Result[Data]
	AddParameter(ctx context.Context, params AddParameterParams) Result[Data]
	UpdParameter(ctx context.Context, params UpdParameterParams) Result[Data]
	DelParameter(ctx context.Context, code string) Result[Data]
	ReloadCache(ctx context.Context) Result[Data]
}

// Parameter Cache Interface
type ParameterCache interface {
	Get(code string) (*Parameter, bool)
	GetValue(code string) (Data, bool)
	Set(code string, param *Parameter)
	Delete(code string)
	LoadAll(params []Parameter)
	GetAll() []Parameter
	Clear()
}
