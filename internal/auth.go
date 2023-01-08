package internal

import (
	"context"
)

const (
	BaseURL = "https://example.com"
)

type AuthParam struct {
	Email       string
	DisplayName string
}

type AuthCreateParam struct {
	Email       string
	Password    string
	DisplayName string
}

type AuthOp interface {
	Signup(context.Context, *AuthCreateParam) (*AuthParam, error)
}

func NewAuth(
	apikey string,
	projectid string,
) (*Auth, error) {
	return &Auth{
		apikey:    apikey,
		projectid: projectid,
	}, nil
}

type Auth struct {
	apikey    string
	projectid string
}

func (r *Auth) Signup(
	ctx context.Context,
	param *AuthCreateParam,
) (*AuthParam, error) {
	// code connecting to server
	return nil, nil
}
