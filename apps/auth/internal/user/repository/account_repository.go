package repository

import (
	"context"
	"errors"

	"github.com/saas-flow/monorepo/apps/auth/internal/user/domain"
	"github.com/saas-flow/monorepo/libs/pagination"
	"gorm.io/gorm"
)

type AccountRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepositoryImpl(db *gorm.DB) domain.AccountRepository {
	return &AccountRepositoryImpl{db}
}

func (r *AccountRepositoryImpl) Find(ctx context.Context, p *pagination.PaginationRequest, f *domain.Account) ([]*domain.Account, error) {
	var result []*domain.Account
	if err := r.db.WithContext(ctx).Scopes(pagination.ApplyPaginationAndFilter(p)).Where(f).Preload("PasswordVersion").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *AccountRepositoryImpl) FindOne(ctx context.Context, f *domain.Account) (*domain.Account, error) {
	var result *domain.Account
	if err := r.db.WithContext(ctx).Where(f).Preload("PasswordVersion").First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (r *AccountRepositoryImpl) Create(ctx context.Context, param *domain.Account) (*domain.Account, error) {
	if err := r.db.WithContext(ctx).Create(param).Error; err != nil {
		return nil, err
	}
	return r.FindOne(ctx, &domain.Account{ID: param.ID})
}

func (r *AccountRepositoryImpl) Update(ctx context.Context, param *domain.Account) (*domain.Account, error) {
	if err := r.db.WithContext(ctx).Save(param).Error; err != nil {
		return nil, err
	}
	return r.FindOne(ctx, &domain.Account{ID: param.ID})
}

func (r *AccountRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Account{}).Error
}
