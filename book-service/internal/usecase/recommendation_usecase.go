package usecase

import (
	"book-service/internal/domain"
	"book-service/internal/repository"
	sharedDomain "book-service/pkg/shared/domain"
	"context"
	"errors"
	"time"
)

type recommendationUsecase struct {
	repo repository.Repository
}

func NewRecommendationUsecase(repo repository.Repository) RecommendationUsecase {
	return &recommendationUsecase{repo: repo}
}

func (uc *recommendationUsecase) CreateRecommendation(ctx context.Context, req *domain.CreateRecommendationRequest) (*sharedDomain.Recommendation, error) {

	_, err := uc.repo.Recommendation().FindByUserAndBook(ctx, req.UserID, req.BookID)
	if err == nil {
		return nil, errors.New("already")
	}

	book, err := uc.repo.Book().FindByID(ctx, req.BookID)
	if err != nil {
		return nil, err
	}

	user, err := uc.repo.User().FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	recommendation := &sharedDomain.Recommendation{
		UserID:    user.ID,
		BookID:    book.ID,
		Reason:    req.Reason,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.Recommendation().Create(ctx, recommendation); err != nil {
		return nil, err
	}

	return recommendation, nil
}

func (uc *recommendationUsecase) GetRecommendationByID(ctx context.Context, id uint) (*sharedDomain.Recommendation, error) {
	return uc.repo.Recommendation().GetByID(ctx, id)
}

func (uc *recommendationUsecase) GetAllRecommendations(ctx context.Context) ([]*sharedDomain.Recommendation, error) {
	return uc.repo.Recommendation().GetAll(ctx)
}

func (uc *recommendationUsecase) DeleteRecommendation(ctx context.Context, id uint) error {
	return uc.repo.Recommendation().Delete(ctx, id)
}
