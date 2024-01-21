package usecase

import (
	"jastip-app/internal/entity/request"
)

type UserUsecase interface {
	Register(request.RegisterRequest) error
	Login(request.LoginRequest) (token string, err error)
}
