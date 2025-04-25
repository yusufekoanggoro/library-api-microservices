package usecase

import (
	"book-service/internal/repository"
	sharedDomain "book-service/pkg/shared/domain"
	protoBook "book-service/proto/book"
	"context"
)

type userUsecase struct {
	repo repository.Repository
}

func NewUserUsecase(
	repo repository.Repository,
) UserUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) SaveUser(ctx context.Context, req *protoBook.UserData) error {
	user := &sharedDomain.User{}
	user, err := user.FromProto(req)
	if err != nil {
		return err
	}
	return uc.repo.User().Upsert(ctx, user)
}

func (uc *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	return uc.repo.User().Delete(ctx, id)
}
