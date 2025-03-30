package repository

import (
	"context"
	"errors"

	"github.com/saas-flow/monorepo/apps/auth/internal/user/domain"
	"github.com/saas-flow/monorepo/libs/pagination"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var AuthProviderModule = fx.Module("auth_provider.repository",
	fx.Provide(NewAuthProviderImpl),
)

type AuthProviderImpl struct {
	db *gorm.DB
}

func NewAuthProviderImpl(db *gorm.DB) domain.AuthProviderRepository {
	return &AuthProviderImpl{db}
}

func (r *AuthProviderImpl) Find(ctx context.Context, p *pagination.PaginationRequest, f *domain.AuthProvider) ([]*domain.AuthProvider, error) {
	var result []*domain.AuthProvider
	if err := r.db.WithContext(ctx).Scopes(pagination.ApplyPaginationAndFilter(p)).Where(f).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *AuthProviderImpl) FindOne(ctx context.Context, f *domain.AuthProvider) (*domain.AuthProvider, error) {
	var result *domain.AuthProvider
	if err := r.db.WithContext(ctx).Where(f).First(&result).Order("created_at DESC").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (r *AuthProviderImpl) Create(ctx context.Context, param *domain.AuthProvider) (*domain.AuthProvider, error) {
	if err := r.db.WithContext(ctx).Create(param).Error; err != nil {
		return nil, err
	}
	return r.FindOne(ctx, &domain.AuthProvider{ID: param.ID})
}

func (r *AuthProviderImpl) Update(ctx context.Context, param *domain.AuthProvider) (*domain.AuthProvider, error) {
	if err := r.db.WithContext(ctx).Save(param).Error; err != nil {
		return nil, err
	}
	return r.FindOne(ctx, &domain.AuthProvider{ID: param.ID})
}

func (r *AuthProviderImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.AuthProvider{}).Error
}
