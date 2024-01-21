package usecase

import (
	"fmt"
	"jastip-app/config"
	"jastip-app/internal/customerror"
	"jastip-app/internal/entity/model"
	"jastip-app/internal/entity/request"
	"jastip-app/internal/usecase/repo"
	"jastip-app/pkg/bcrypt"
	"jastip-app/pkg/jwt"
	"jastip-app/pkg/logger"

	"github.com/google/uuid"
)

type UserUsecase interface {
	Register(request.RegisterRequest) error
	Login(request.LoginRequest) (token string, err error)
}

type userUsecase struct {
	userRepo repo.UserRepo
	config   *config.Config
}

func NewUserUsecase(userRepo repo.UserRepo, config *config.Config) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		config:   config,
	}
}

func (uc *userUsecase) Register(req request.RegisterRequest) error {

	userd, err := uc.userRepo.GetByEmail(req.Email)
	fmt.Println(userd)
	if err == nil {
		return &customerror.Err{
			Code:   customerror.UserAlreadyRegistered,
			Errors: nil,
		}
	}
	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return err
	}
	id, err := uuid.NewRandom()
	if err != nil {
		logger.L().Error(err.Error())
		return err
	}

	user := model.UserModel{
		ID:          id.String(),
		Name:        req.Name,
		Address:     req.Address,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
		Role:        model.RoleUser(model.ROLE_USER),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return err
	}
	return nil
}

func (uc *userUsecase) Login(req request.LoginRequest) (token string, err error) {
	user, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", err
	}
	if !bcrypt.CheckPasswordHash(req.Password, user.Password) {
		return "", &customerror.Err{
			Code:   customerror.CodeErrInvalidRequest,
			Errors: "wrong email or password",
		}
	}
	token, err = jwt.GenerateToken(jwt.JWTPayload{
		SecretKey: uc.config.JWT.SECRET,
		Expired:   float64(uc.config.JWT.TTL),
		ID:        user.ID,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
