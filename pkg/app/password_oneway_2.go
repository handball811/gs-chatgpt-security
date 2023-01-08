//go:build password_oneway_2
// +build password_oneway_2

package app

import (
	"context"
	"errors"

	"github.com/handball811/gs-chatgpt-security/pkg/interface/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, model repository.UserModel) (*repository.UserModel, error) {
	if model.Email == nil || model.Password == nil {
		return nil, errors.New("email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*model.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	model.Password = &string(hashedPassword)

	return s.repo.Create(ctx, model)
}

func (s *UserService) Validate(ctx context.Context, email string, password string) (*repository.UserModel, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, email string, oldPassword string, newPassword string) (*repository.UserModel, error) {
	user, err := s.Validate(ctx, email, oldPassword)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = &string(hashedPassword)

	return s.repo.Update(ctx, *user.Id, *user)
}
