package service

import (
	"context"
	"payments/models"
	"payments/repository"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (id int, err error)
	GenerateToken(ctx context.Context, user models.User) (signedToken string, err error)
	ParseToken(accessToken string) (id int, err error)
}

type Service struct {
	Authorization
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthorizationService(r),
	}
}
