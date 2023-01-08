package repository

import "context"

type UserModel struct {
	Id          *string
	Email       *string
	Password    *string
	DisplayName *string
}

type UserRepo interface {
	Create(ctx context.Context, model UserModel) (*UserModel, error)
	Get(ctx context.Context, id string) (*UserModel, error)
	GetByEmail(ctx context.Context, email string) (*UserModel, error)
	Update(ctx context.Context, id string, model UserModel) (*UserModel, error)
}
