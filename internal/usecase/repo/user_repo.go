package repo

import (
	"jastip-app/internal/customerror"
	"jastip-app/internal/entity/model"
	"jastip-app/pkg/logger"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(user model.UserModel) error
	GetByEmail(email string) (model.UserModel, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(user model.UserModel) error {
	if err := r.db.Debug().Table(user.TableName()).Create(&user).Error; err != nil {
		logger.L().Error(err.Error())
		return err
	}
	return nil

}

func (r *userRepo) GetByEmail(email string) (model.UserModel, error) {

	var user model.UserModel

	err := r.db.Debug().Table(user.TableName()).Where("email= ?", email).Find(&user).Error
	if err == gorm.ErrRecordNotFound || user.ID == "" {
		return user, &customerror.Err{
			Code:   customerror.CodeErrNotFound,
			Errors: nil,
		}
	}
	if err != nil {
		logger.L().Error(err.Error())
		return user, err
	}

	return user, nil
}
