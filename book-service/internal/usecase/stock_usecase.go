package usecase

import (
	"book-service/internal/domain"
	"book-service/internal/repository"
	sharedDomain "book-service/pkg/shared/domain"
	"context"
	"log"
)

type stockUsecase struct {
	repo repository.Repository
}

func NewStockUsecase(
	repo repository.Repository,

) StockUsecase {
	return &stockUsecase{
		repo: repo,
	}
}

func (s *stockUsecase) IncreaseStock(ctx context.Context, req *domain.IncreaseStockRequest) (*sharedDomain.Book, error) {
	log.Println("IncreaseStock request:", req)

	err := s.repo.WithTransaction(func(txRepo repository.Repository) error {
		err := txRepo.Stock().IncreaseStock(ctx, req.ID, req.Amount)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.Book().FindByID(ctx, req.ID)
}

func (s *stockUsecase) DecreaseStock(ctx context.Context, req *domain.DecreaseStockRequest) (*sharedDomain.Book, error) {
	log.Println("DecreaseStock request:", req)

	err := s.repo.WithTransaction(func(txRepo repository.Repository) error {
		err := txRepo.Stock().DecreaseStock(ctx, req.ID, req.Amount)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.Book().FindByID(ctx, req.ID)
}
