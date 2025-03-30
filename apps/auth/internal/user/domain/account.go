package domain

import (
	"context"

	"github.com/saas-flow/monorepo/apps/auth/internal/user/dto"
	"github.com/saas-flow/monorepo/libs/pagination"
)

type AccountRepository interface {
	Find(context.Context, *pagination.PaginationRequest, *Account) ([]*Account, error)
	FindOne(context.Context, *Account) (*Account, error)
	Create(context.Context, *Account) (*Account, error)
	Update(context.Context, *Account) (*Account, error)
	Delete(context.Context, string) error
}

type PasswordVersionRepository interface {
	Find(context.Context, *pagination.PaginationRequest, *PasswordVersion) ([]*PasswordVersion, error)
	FindOne(context.Context, *PasswordVersion) (*PasswordVersion, error)
	Create(context.Context, *PasswordVersion) (*PasswordVersion, error)
	Update(context.Context, *PasswordVersion) (*PasswordVersion, error)
	Delete(context.Context, string) error
}

type AuthProviderRepository interface {
	Find(context.Context, *pagination.PaginationRequest, *AuthProvider) ([]*AuthProvider, error)
	FindOne(context.Context, *AuthProvider) (*AuthProvider, error)
	Create(context.Context, *AuthProvider) (*AuthProvider, error)
	Update(context.Context, *AuthProvider) (*AuthProvider, error)
	Delete(context.Context, string) error
}

type AccountService interface {
	ListProvider(context.Context) ([]*AuthProvider, error)
	GetPasswordVersion(context.Context) (*PasswordVersion, error)
	ValidatePassword(context.Context, *dto.ValidatePasswordRequest) error
	SkipUpdatePassword(context.Context) error
	Lookup(context.Context, string) error
	SignInWithPassword(context.Context, *dto.SignInWithPasswordRequest) (*Account, error)
	SignUp(context.Context, *dto.SignUpRequest) (*Account, error)
	// ForgotPassword(context.Context) error
	// ChangePassword(context.Context) error
}
