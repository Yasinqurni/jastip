package repo

import (
	"jastip-app/internal/customerror"
	"jastip-app/internal/entity/model"
	"jastip-app/pkg/logger"

	"gorm.io/gorm"
)

type SubscribeClassRepo interface {
	Create(model.SubscribeClassModel) error
	Update(data model.SubscribeClassModel) error
	GetAll() ([]model.SubscribeClassModel, error)
	GetByID(ID string) (model.SubscribeClassModel, error)
	Delete(ID string) error
}

type subscribeClassRepo struct {
	db *gorm.DB
}

func NewSubscribeClassRepo(db *gorm.DB) SubscribeClassRepo {
	return &subscribeClassRepo{
		db: db,
	}
}

func (r *subscribeClassRepo) Create(data model.SubscribeClassModel) error {
	if err := r.db.Debug().Table(data.TableName()).Create(&data).Error; err != nil {
		logger.L().Error(err.Error())
		return err
	}
	return nil
}

func (r *subscribeClassRepo) Update(data model.SubscribeClassModel) error {
	if err := r.db.Debug().Table(data.TableName()).Where("id = ?", data.ID).Updates(&data).Error; err != nil {
		logger.L().Error(err.Error())
		return err
	}
	return nil
}

func (r *subscribeClassRepo) GetAll() ([]model.SubscribeClassModel, error) {
	var (
		subscribeClasses []model.SubscribeClassModel
		subscribeClass   model.SubscribeClassModel
	)
	if err := r.db.
		Debug().
		Table(subscribeClass.TableName()).
		Find(&subscribeClasses).
		Error; err != nil {
		logger.L().Error(err.Error())
		return nil, err
	}
	return subscribeClasses, nil
}

func (r *subscribeClassRepo) GetByID(ID string) (model.SubscribeClassModel, error) {
	var subscribeClass model.SubscribeClassModel
	if err := r.db.
		Debug().
		Table(subscribeClass.TableName()).
		Where("id = ?", ID).
		Find(&subscribeClass).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound || subscribeClass.ID == "" {
			return subscribeClass, &customerror.Err{
				Code:   customerror.CodeErrNotFound,
				Errors: nil,
			}
		}
		if err != nil {
			logger.L().Error(err.Error())
			return subscribeClass, err
		}
	}
	return subscribeClass, nil
}

func (r *subscribeClassRepo) Delete(ID string) error {
	var subscribeClass model.SubscribeClassModel

	condition := r.db.Debug().Where("id = ?", ID).Delete(&subscribeClass)

	if condition.Error != nil {
		logger.L().Error(condition.Error.Error())
		return condition.Error
	}

	if condition.RowsAffected == 0 {
		return &customerror.Err{
			Code:   customerror.CodeErrNotFound,
			Errors: nil,
		}
	}

	return nil
}
