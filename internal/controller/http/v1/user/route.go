package user

import (
	"jastip-app/internal/usecase"
	"jastip-app/pkg/app"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, userUsecase usecase.UserUsecase, app *app.App) {
	h := NewUserHandler(userUsecase)
	routes := r.Group("/user")
	{
		//register
		routes.POST("/register", h.Register)
		//login
		routes.POST("/login", h.Login)

	}
}
