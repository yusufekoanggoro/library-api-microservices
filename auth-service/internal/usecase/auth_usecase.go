package usecase

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"auth-service/pkg/logger"
	"auth-service/pkg/token"
	"context"
	"errors"

	"time"

	"auth-service/internal/grpcservice"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo     repository.UserRepository
	tokenService token.Token
	bookClient   *grpcservice.BookGRPCClient
	logger       logger.Logger
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	tokenService token.Token,
	bookClient *grpcservice.BookGRPCClient,
	logger logger.Logger,
) AuthUsecase {
	return &authUsecase{
		userRepo:     userRepo,
		tokenService: tokenService,
		bookClient:   bookClient,
		logger:       logger,
	}
}

func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessTokenExpiration := time.Hour * 1
	token, err := u.tokenService.GenerateToken(user.ID, user.Role, accessTokenExpiration)
	if err != nil {
		return nil, err
	}

	refreshTokenExpiration := 7 * 24 * time.Hour
	refreshToken, err := u.tokenService.GenerateToken(user.ID, user.Role, refreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	grpcRequest := user.ToProto()
	_, err = u.bookClient.SaveUser(ctx, grpcRequest)
	if err != nil {
		u.logger.Error(err.Error(), "bookClient.SaveUser", "grpc request")
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:           token,
		AccessTokenExpiresIn:  int(accessTokenExpiration.Seconds()),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: int(refreshTokenExpiration.Seconds()),
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponse, error) {
	tokenClaims, err := u.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := u.userRepo.GetUserByID(ctx, tokenClaims.UserID)
	if err != nil {
		return nil, err
	}

	accessTokenExpiration := time.Hour * 1
	newAccessToken, err := u.tokenService.GenerateToken(tokenClaims.UserID, user.Role, accessTokenExpiration)
	if err != nil {
		return nil, err
	}

	refreshTokenExpiration := 7 * 24 * time.Hour
	refreshToken, err = u.tokenService.GenerateToken(tokenClaims.UserID, user.Role, refreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshTokenResponse{
		AccessToken:           newAccessToken,
		AccessTokenExpiresIn:  int(accessTokenExpiration.Seconds()),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: int(refreshTokenExpiration.Seconds()),
	}, nil
}
