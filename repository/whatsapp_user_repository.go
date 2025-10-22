package repository

import (
	"context"
	"fmt"

	"api-chatbot/api/dal"
	d "api-chatbot/domain"
)

const (
	// Functions (Read-only)
	fnGetUserByIdentity = "fn_get_user_by_identity"
	fnGetUserByWhatsApp = "fn_get_user_by_whatsapp"
	// Stored Procedures (Writes)
	spCreateUser         = "sp_create_user"
	spUpdateUserWhatsApp = "sp_update_user_whatsapp"
)

type whatsappUserRepository struct {
	dal *dal.DAL
}

func NewWhatsAppUserRepository(dal *dal.DAL) d.WhatsAppUserRepository {
	return &whatsappUserRepository{
		dal: dal,
	}
}

// GetByIdentity retrieves a user by their identity number
func (r *whatsappUserRepository) GetByIdentity(ctx context.Context, identityNumber string) (*d.WhatsAppUser, error) {
	users, err := dal.QueryRows[d.WhatsAppUser](r.dal, ctx, fnGetUserByIdentity, identityNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get user via %s: %w", fnGetUserByIdentity, err)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

// GetByWhatsApp retrieves a user by their WhatsApp number
func (r *whatsappUserRepository) GetByWhatsApp(ctx context.Context, whatsapp string) (*d.WhatsAppUser, error) {
	users, err := dal.QueryRows[d.WhatsAppUser](r.dal, ctx, fnGetUserByWhatsApp, whatsapp)
	if err != nil {
		return nil, fmt.Errorf("failed to get user via %s: %w", fnGetUserByWhatsApp, err)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

// Create creates a new user
func (r *whatsappUserRepository) Create(ctx context.Context, params d.CreateUserParams) (*d.CreateUserResult, error) {
	result, err := dal.ExecProc[d.CreateUserResult](
		r.dal,
		ctx,
		spCreateUser,
		params.IdentityNumber,
		params.Name,
		params.Email,
		params.Phone,
		params.Role,
		params.WhatsApp,
		params.Details,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spCreateUser, err)
	}

	return result, nil
}

// UpdateWhatsApp updates a user's WhatsApp number
func (r *whatsappUserRepository) UpdateWhatsApp(ctx context.Context, params d.UpdateUserWhatsAppParams) error {
	_, err := dal.ExecProc[dal.DbResult](
		r.dal,
		ctx,
		spUpdateUserWhatsApp,
		params.IdentityNumber,
		params.WhatsApp,
	)

	if err != nil {
		return fmt.Errorf("failed to execute %s: %w", spUpdateUserWhatsApp, err)
	}

	return nil
}
