package repository

import (
	"context"
	"errors"

	"github.com/saas-flow/monorepo/apps/auth/internal/user/domain"
	"github.com/saas-flow/monorepo/libs/pagination"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var PasswordVersionModule = fx.Module("password_version.repository",
	fx.Provide(NewPasswordVersionImpl),
)

type PasswordVersionImpl struct {
	db *gorm.DB
}

func NewPasswordVersionImpl(db *gorm.DB) domain.PasswordVersionRepository {
	return &PasswordVersionImpl{db}
}

func (r *PasswordVersionImpl) Find(ctx context.Context, p *pagination.PaginationRequest, f *domain.PasswordVersion) ([]*domain.PasswordVersion, error) {
	var result []*domain.PasswordVersion
	if err := r.db.WithContext(ctx).Scopes(pagination.ApplyPaginationAndFilter(p)).Where(f).Preload("Rules").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PasswordVersionImpl) FindOne(ctx context.Context, f *domain.PasswordVersion) (*domain.PasswordVersion, error) {
	var result *domain.PasswordVersion
	if err := r.db.WithContext(ctx).Where(f).Preload("Rules").First(&result).Order("created_at DESC").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (r *PasswordVersionImpl) Create(ctx context.Context, param *domain.PasswordVersion) (*domain.PasswordVersion, error) {
	if err := r.db.WithContext(ctx).Create(param).Error; err != nil {
		return nil, err
	}
	return r.FindOne(ctx, &domain.PasswordVersion{ID: param.ID})
}

func (r *PasswordVersionImpl) Update(ctx context.Context, param *domain.PasswordVersion) (*domain.PasswordVersion, error) {
	if err := r.db.WithContext(ctx).Save(param).Error; err != nil {
		return nil, err
	}
	return r.FindOne(ctx, &domain.PasswordVersion{ID: param.ID})
}

func (r *PasswordVersionImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.PasswordVersion{}).Error
}
