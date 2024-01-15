package usecase

import (
	"jastip-app/pkg/app"
)

type UseCase struct {
}

func New(app *app.App) *UseCase {
	return &UseCase{}
}
