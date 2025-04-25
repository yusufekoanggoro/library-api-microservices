package usecase

import (
	"book-service/internal/domain"
	"book-service/internal/repository"
	"book-service/pkg/redis"
	sharedDomain "book-service/pkg/shared/domain"
	"context"
	"errors"
	"fmt"
	"time"
)

type borrowingUseCase struct {
	repo   repository.Repository
	locker redis.Locker
}

func NewBorrowingUseCase(
	repo repository.Repository,
	locker redis.Locker,
) BorrowingUseCase {
	return &borrowingUseCase{
		repo:   repo,
		locker: locker,
	}
}

func (u *borrowingUseCase) BorrowBook(ctx context.Context, req *domain.BorrowBookRequest) (*sharedDomain.Borrowing, error) {
	lockKey := fmt.Sprintf("book:%d", req.BookID)

	success, err := u.locker.AcquireLock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return nil, err
	}

	if !success {
		return nil, err
	}

	defer u.locker.ReleaseLock(ctx, lockKey)

	_, err = u.repo.Borrowing().FindByUserAndStatus(ctx, req.UserID, "borrowed")
	if err == nil {
		return nil, errors.New("already borrowing")
	}

	book, err := u.repo.Book().FindByID(ctx, req.BookID)
	if err != nil {
		return nil, err
	}
	if book.Stock <= 0 {
		return nil, err
	}

	user, err := u.repo.User().FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	borrowing := &sharedDomain.Borrowing{
		UserID:      user.ID,
		BookID:      book.ID,
		BorrowDate:  time.Now(),
		Status:      "borrowed",
		CreatedByID: req.CreatedByID,
		UpdatedByID: req.CreatedByID,
	}

	err = u.repo.WithTransaction(func(txRepo repository.Repository) error {
		if err := txRepo.Borrowing().Create(ctx, borrowing); err != nil {
			return err
		}

		if err = txRepo.Stock().DecreaseStock(ctx, req.BookID, 1); err != nil {
			return err
		}

		return nil
	})

	borrowing, err = u.repo.Borrowing().FindByID(ctx, borrowing.ID)
	if err != nil {
		return nil, err
	}

	return borrowing, nil
}

func (u *borrowingUseCase) ReturnBook(ctx context.Context, req *domain.ReturnBookRequest) (*sharedDomain.Borrowing, error) {
	lockKey := fmt.Sprintf("book:%d", req.BookID)

	success, err := u.locker.AcquireLock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, err
	}

	defer u.locker.ReleaseLock(ctx, lockKey)

	_, err = u.repo.Borrowing().FindByUserAndBook(ctx, req.UserID, req.BookID)
	if err != nil {
		return nil, err
	}

	borrowing, err := u.repo.Borrowing().FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if borrowing.Status != "borrowed" {
		return nil, errors.New("book returned")
	}

	borrowing.Status = "returned"
	now := time.Now()
	borrowing.ReturnDate = &now
	borrowing.UpdatedByID = req.UpdatedByID

	err = u.repo.WithTransaction(func(txRepo repository.Repository) error {
		if err := txRepo.Borrowing().Update(ctx, borrowing); err != nil {
			return err
		}

		if err = txRepo.Stock().IncreaseStock(ctx, req.BookID, 1); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return borrowing, nil
}

func (uc *borrowingUseCase) ListBorrowings(ctx context.Context) ([]*sharedDomain.Borrowing, error) {
	borrowings, err := uc.repo.Borrowing().FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return borrowings, nil
}
