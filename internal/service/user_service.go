package service

import (
	"context"
	"errors"

	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrWrongPassword      = errors.New("wrong current password")
)

var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)

type UserService struct {
	repo         *repository.UserRepository
	tokenService *TokenService
}

func NewUserService(repo *repository.UserRepository, tokenService *TokenService) *UserService {
	return &UserService{repo: repo, tokenService: tokenService}
}

func (s *UserService) Register(ctx context.Context, body dto.UserRegisterDTO) (*dto.AuthResponseDTO, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: body.Username,
		Password: string(hash),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.authResponse(ctx, user)
}

func (s *UserService) Login(ctx context.Context, body dto.UserLoginDTO) (*dto.AuthResponseDTO, error) {
	user, err := s.repo.FindByUsername(ctx, body.Username)
	if errors.Is(err, repository.ErrUserNotFound) {
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(body.Password))
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.authResponse(ctx, user)
}

func (s *UserService) Refresh(ctx context.Context, body dto.RefreshTokenDTO) (*dto.AuthResponseDTO, error) {
	userID, err := s.tokenService.Consume(ctx, body.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(ctx, userID)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrInvalidRefreshToken
	}
	if err != nil {
		return nil, err
	}

	return s.authResponse(ctx, user)
}

func (s *UserService) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, current *model.User, body dto.UserUpdateDTO) (*model.User, error) {
	user := *current

	if body.KcalGoal != nil {
		user.KcalGoal = uint16(*body.KcalGoal)
	}
	if body.ProteinGoal != nil {
		user.ProteinGoal = uint16(*body.ProteinGoal)
	}
	if body.CarbohydrateGoal != nil {
		user.CarbohydrateGoal = uint16(*body.CarbohydrateGoal)
	}
	if body.FatGoal != nil {
		user.FatGoal = uint16(*body.FatGoal)
	}

	if err := s.repo.Update(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, user *model.User, body dto.UserPasswordUpdateDTO) (*dto.AuthResponseDTO, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.CurrentPassword)); err != nil {
		return nil, ErrWrongPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePassword(ctx, user.ID, string(hash)); err != nil {
		return nil, err
	}

	if err := s.tokenService.RevokeAll(ctx, user.ID); err != nil {
		return nil, err
	}

	return s.authResponse(ctx, user)
}

func (s *UserService) authResponse(ctx context.Context, user *model.User) (*dto.AuthResponseDTO, error) {
	tokens, err := s.tokenService.Issue(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponseDTO{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         user,
	}, nil
}
