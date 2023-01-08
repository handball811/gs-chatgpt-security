//go:build password_oneway_3
// +build password_oneway_3

package app

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/handball811/gs-chatgpt-security/pkg/interface/repository"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, email, password, displayName string) (*repository.UserModel, error) {
	// メールアドレスのバリデーション処理
	if !validateEmail(email) {
		return nil, errors.New("invalid email format")
	}

	// パスワードのバリデーション処理
	if !validatePassword(password) {
		return nil, errors.New("invalid password format")
	}

	// パスワードをハッシュ化する
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// バイト列を文字列に変換する
	passwordStr := string(hashedPassword)

	model := repository.UserModel{
		Email:       &email,
		Password:    &passwordStr,
		DisplayName: &displayName,
	}
	return s.repo.Create(ctx, model)
}

func (s *UserService) Validate(ctx context.Context, email, password string) (*repository.UserModel, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, email, password, newPassword string) (*repository.UserModel, error) {
	user, err := s.Validate(ctx, email, password)
	if err != nil {
		return nil, err
	}

	// パスワードのバリデーション処理
	if !validatePassword(newPassword) {
		return nil, errors.New("invalid password format")
	}

	// パスワードをハッシュ化する
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// バイト列を文字列に変換する
	newPasswordStr := string(hashedPassword)
	user.Password = &newPasswordStr
	return s.repo.Update(ctx, *user.Id, *user)
}

func validateEmail(email string) bool {
	// emailのバリデーション処理
	return true
}

func validatePassword(password string) bool {
	// passwordのバリデーション処理
	return true
}
