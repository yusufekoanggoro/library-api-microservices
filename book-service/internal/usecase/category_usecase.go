package usecase

import (
	"book-service/internal/repository"
	sharedDomain "book-service/pkg/shared/domain"
	protoBook "book-service/proto/book"
	"context"
)

type categoryUsecase struct {
	repo repository.Repository
}

func NewCategoryUsecase(
	repo repository.Repository,
) CategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (uc *categoryUsecase) SaveCategory(ctx context.Context, req *protoBook.CategoryData) error {
	category := &sharedDomain.Category{}
	category, err := category.FromProto(req)
	if err != nil {
		return err
	}
	return uc.repo.Category().Upsert(ctx, category)
}

func (uc *categoryUsecase) DeleteCategory(ctx context.Context, id uint) error {
	return uc.repo.Category().Delete(ctx, id)
}
