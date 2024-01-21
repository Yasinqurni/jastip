package usecase

import (
	"jastip-app/internal/usecase/repo"
	"jastip-app/pkg/app"
)

type UseCase struct {
	User UserUsecase
	App  *app.App
}

func New(app *app.App) *UseCase {
	userUsecase := NewUserUsecase(repo.NewUserRepo(app.DB.Gorm), app.Config)
	return &UseCase{
		User: userUsecase,
		App:  app,
	}
}
