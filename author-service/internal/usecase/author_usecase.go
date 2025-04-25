package usecase

import (
	"author-service/internal/domain"
	"author-service/internal/grpcservice"

	"author-service/internal/repository"
	sharedDomain "author-service/pkg/shared/domain"
	"context"
)

type authorUsecase struct {
	repo       repository.AuthorRepository
	bookClient *grpcservice.BookGRPCClient
}

func NewAuthorUsecase(
	repo repository.AuthorRepository,
	bookClient *grpcservice.BookGRPCClient,
) AuthorUsecase {
	return &authorUsecase{
		repo:       repo,
		bookClient: bookClient,
	}
}

func (uc *authorUsecase) CreateAuthor(ctx context.Context, req *domain.CreateAuthorRequest) (*sharedDomain.Author, error) {
	var author *sharedDomain.Author
	newAuthor := &sharedDomain.Author{
		Name:        req.Name,
		Bio:         req.Bio,
		CreatedByID: req.CreatedByID,
		UpdatedByID: req.CreatedByID,
	}

	err := uc.repo.SaveAuthor(ctx, newAuthor)
	if err != nil {
		return nil, err
	}
	author = newAuthor

	grpcRequest := author.ToProto()
	_, err = uc.bookClient.SaveAuthor(ctx, grpcRequest)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (uc *authorUsecase) GetAllAuthors(ctx context.Context, req *domain.PaginationRequest) (*domain.PaginatedResponse, error) {
	authors, totalRows, err := uc.repo.GetAllAuthors(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	paginatedResponse := &domain.PaginatedResponse{
		Data:       authors,
		Total:      totalRows,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((totalRows + int64(req.Limit) - 1) / int64(req.Limit)),
	}

	return paginatedResponse, nil
}

func (uc *authorUsecase) GetAuthorByID(ctx context.Context, id uint) (*sharedDomain.Author, error) {
	return uc.repo.GetAuthorByID(ctx, id)
}

func (uc *authorUsecase) UpdateAuthor(ctx context.Context, req *domain.UpdateAuthorRequest) (*sharedDomain.Author, error) {
	var author *sharedDomain.Author
	existingAuthor, err := uc.repo.GetAuthorByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existingAuthor.Name = *req.Name
	}
	if req.Bio != nil {
		existingAuthor.Bio = *req.Bio
	}

	existingAuthor.UpdatedByID = req.UpdatedByID

	err = uc.repo.SaveAuthor(ctx, existingAuthor)
	if err != nil {
		return nil, err
	}
	author = existingAuthor

	grpcRequest := author.ToProto()
	_, err = uc.bookClient.SaveAuthor(ctx, grpcRequest)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (uc *authorUsecase) DeleteAuthor(ctx context.Context, id uint) error {
	err := uc.repo.DeleteAuthor(ctx, id)
	if err != nil {
		return err
	}

	_, err = uc.bookClient.DeleteAuthor(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
