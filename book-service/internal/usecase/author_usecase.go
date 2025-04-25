package usecase

import (
	"book-service/internal/repository"
	sharedDomain "book-service/pkg/shared/domain"
	protoBook "book-service/proto/book"
	"context"
)

type authorUsecase struct {
	repo repository.Repository
}

func NewAuthorUsecase(
	repo repository.Repository,
) AuthorUsecase {
	return &authorUsecase{repo: repo}
}

func (uc *authorUsecase) SaveAuthor(ctx context.Context, req *protoBook.AuthorData) error {
	author := &sharedDomain.Author{}
	author, err := author.FromProto(req)
	if err != nil {
		return err
	}
	return uc.repo.Author().Upsert(ctx, author)
}

func (uc *authorUsecase) DeleteAuthor(ctx context.Context, id uint) error {
	err := uc.repo.Author().Delete(ctx, id)
	return err
}
