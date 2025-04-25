package usecase

import (
	"auth-service/internal/domain"
	"auth-service/internal/grpcservice"
	"auth-service/internal/repository"
	"auth-service/pkg/logger"
	sharedDomain "auth-service/pkg/shared/domain"
	"context"
)

type userUsecase struct {
	repo       repository.UserRepository
	bookClient *grpcservice.BookGRPCClient
	logger     logger.Logger
}

func NewAuthorUsecase(
	repo repository.UserRepository,
	bookClient *grpcservice.BookGRPCClient,
	logger logger.Logger,
) UserUsecase {
	return &userUsecase{repo: repo, bookClient: bookClient, logger: logger}
}

func (uc *userUsecase) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*sharedDomain.User, error) {
	var user *sharedDomain.User

	newUser := &sharedDomain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := newUser.HashPassword(); err != nil {
		return nil, err
	}

	err := uc.repo.SaveUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	user = newUser

	grpcRequest := user.ToProto()
	_, err = uc.bookClient.SaveUser(ctx, grpcRequest)
	if err != nil {
		uc.logger.Error(err.Error(), "bookClient.SaveUser", "grpc request")
		return nil, err
	}

	return user, nil
}

func (uc *userUsecase) GetAllUsers(ctx context.Context, req *domain.PaginationRequest) (*domain.PaginatedResponse, error) {
	users, totalRows, err := uc.repo.GetAllUsers(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	paginatedResponse := &domain.PaginatedResponse{
		Data:       users,
		Total:      totalRows,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((totalRows + int64(req.Limit) - 1) / int64(req.Limit)),
	}

	return paginatedResponse, nil
}

func (uc *userUsecase) GetUserByID(ctx context.Context, id uint) (*sharedDomain.User, error) {
	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *userUsecase) UpdateUser(ctx context.Context, req *domain.UpdateUserRequest) (*sharedDomain.User, error) {
	var user *sharedDomain.User
	existingUser, err := uc.repo.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Username != nil {
		existingUser.Username = *req.Username
	}

	if req.Email != nil {
		existingUser.Email = *req.Email
	}

	if req.Password != nil {
		existingUser.Password = *req.Password
		if err := existingUser.HashPassword(); err != nil {
			return nil, err
		}
	}

	if req.Role != nil {
		existingUser.Role = *req.Role
	}

	err = uc.repo.SaveUser(ctx, existingUser)
	if err != nil {
		return nil, err
	}
	user = existingUser

	grpcRequest := user.ToProto()
	_, err = uc.bookClient.SaveUser(ctx, grpcRequest)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	err := uc.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	_, err = uc.bookClient.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
