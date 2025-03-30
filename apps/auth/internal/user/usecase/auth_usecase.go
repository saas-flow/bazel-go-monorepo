package usecase

import (
	"context"

	"github.com/saas-flow/monorepo/apps/auth/internal/user/domain"
	"github.com/saas-flow/monorepo/apps/auth/internal/user/dto"
	"github.com/saas-flow/monorepo/libs/pagination"
	"github.com/saas-flow/monorepo/libs/response"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AuthUsecaseModule = fx.Module("auth.usecase",
	fx.Provide(NewAuthUsecase),
)

type AuthUsecase struct {
	passwordVersionRepository domain.PasswordVersionRepository
	authProviderRepository    domain.AuthProviderRepository
	accountRepository         domain.AccountRepository
}

func NewAuthUsecase(
	passwordVersionRepository domain.PasswordVersionRepository,
	authProviderRepository domain.AuthProviderRepository,
	accountRepository domain.AccountRepository,
) domain.AccountService {
	return &AuthUsecase{
		passwordVersionRepository: passwordVersionRepository,
		authProviderRepository:    authProviderRepository,
		accountRepository:         accountRepository,
	}
}

func (s *AuthUsecase) ListProvider(ctx context.Context) ([]*domain.AuthProvider, error) {
	providers, err := s.authProviderRepository.Find(ctx, &pagination.PaginationRequest{
		SortBy:  "created_at",
		OrderBy: "DESC",
	}, &domain.AuthProvider{})
	if err != nil {
		return nil, response.SendError("InternalServerError", err.Error())
	}

	return providers, nil
}

// GetPasswordVersion
func (s *AuthUsecase) GetPasswordVersion(ctx context.Context) (*domain.PasswordVersion, error) {

	passwordVersion, err := s.passwordVersionRepository.FindOne(ctx, &domain.PasswordVersion{})
	if err != nil {
		return nil, response.SendError("InternalServerError", err.Error())
	}

	if passwordVersion == nil {
		return nil, response.SendError("InternalServerError", "password version not found!")
	}

	return passwordVersion, nil
}

// ValidatePassword
func (s *AuthUsecase) ValidatePassword(ctx context.Context, req *dto.ValidatePasswordRequest) error {
	passwordVersion, err := s.GetPasswordVersion(ctx)
	if err != nil {
		return err
	}

	if !passwordVersion.ValidatePassword(req.Password) {
		return response.SendError("ErrPasswordMissMatch", "password does not meet the security requirements.")
	}

	return nil
}

// SkipUpdatePassword
func (s *AuthUsecase) SkipUpdatePassword(ctx context.Context) error {
	return nil
}

// Lookup
func (s *AuthUsecase) Lookup(ctx context.Context, username string) error {

	exist, err := s.accountRepository.FindOne(ctx, &domain.Account{
		Type:     domain.User,
		Username: username,
	})
	if err != nil {
		zap.L().Error("failed query user", zap.Error(err))
		return response.SendError("InternalServerError", err.Error())
	}

	if exist != nil {
		return response.SendError("UserExist", "user already exist!")
	}

	return nil
}

// SignInWithPassword
func (s *AuthUsecase) SignInWithPassword(ctx context.Context, req *dto.SignInWithPasswordRequest) (*domain.Account, error) {

	exist, err := s.accountRepository.FindOne(ctx, &domain.Account{
		Type:     domain.User,
		Username: req.Username,
	})
	if err != nil {
		zap.L().Error("failed query user", zap.Error(err))
		return nil, response.SendError("InternalServerError", err.Error())
	}

	if exist == nil {
		return nil, response.SendError("InvalidCredential", "invalid email or password")
	}

	if !exist.ComparePassword(req.Password) {
		return nil, response.SendError("InvalidCredential", "invalid email or password")
	}

	return exist, nil
}

// SignUp
func (s *AuthUsecase) SignUp(ctx context.Context, req *dto.SignUpRequest) (*domain.Account, error) {

	lastPasswordVersion, err := s.GetPasswordVersion(ctx)
	if err != nil {
		return nil, err
	}

	exist, err := s.accountRepository.FindOne(ctx, &domain.Account{
		Type:     domain.User,
		Username: req.Username,
	})
	if err != nil {
		zap.L().Error("failed query user", zap.Error(err))
		return nil, err
	}

	if exist != nil {
		return nil, response.SendError("ErrAccountExist", "account already exist.")
	}

	if !lastPasswordVersion.ValidatePassword(req.Password) {
		return nil, response.SendError("ErrPasswordMissMatch", "password does not meet the security requirements.")
	}

	user, err := s.accountRepository.Create(ctx, &domain.Account{
		Type:              domain.User,
		Username:          req.Username,
		Password:          req.Password,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		PasswordVersionID: lastPasswordVersion.ID,
	})
	if err != nil {
		return nil, response.SendError("InternalServerError", err.Error())
	}

	return user, nil
}
