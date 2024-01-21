package user

import (
	"jastip-app/internal/customerror"
	"jastip-app/internal/entity/request"
	"jastip-app/internal/usecase"
	"jastip-app/pkg/logger"
	"jastip-app/pkg/recovery"
	"jastip-app/pkg/res"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}
type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) Register(ctx *gin.Context) {
	defer recovery.Recover(ctx)

	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationErr, ok := err.(validator.ValidationErrors)
		if !ok {
			logger.L().Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, res.JSON(false, "Something went wrong", nil))
			return
		}

		var errorMessages = make(map[string]string)
		for _, e := range validationErr {
			fieldJSONName := req.GetJsonFieldName(e.Field())
			errorMessages[fieldJSONName] = req.ErrMessages()[fieldJSONName][e.ActualTag()]
		}
		ctx.JSON(
			http.StatusUnprocessableEntity,
			res.JSON(false, "Failed to register", &customerror.Err{
				Code:   customerror.CodeErrInvalidRequest,
				Errors: errorMessages,
			}),
		)
		return
	}

	err := h.userUsecase.Register(req)
	if err != nil {
		if customErr, ok := err.(*customerror.Err); ok {
			ctx.JSON(http.StatusBadRequest, res.JSON(false, "Failed to register", customErr))
			return
		}
		ctx.JSON(http.StatusInternalServerError, res.JSON(false, "Something went wrong", nil))
		return
	}

	ctx.JSON(http.StatusOK, res.JSON(true, "Success to register", nil))
}

func (h *userHandler) Login(ctx *gin.Context) {
	defer recovery.Recover(ctx)

	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationErr, ok := err.(validator.ValidationErrors)
		if !ok {
			logger.L().Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, res.JSON(false, "Something went wrong", nil))
			return
		}

		var errorMessages = make(map[string]string)
		for _, e := range validationErr {
			fieldJSONName := req.GetJsonFieldName(e.Field())
			errorMessages[fieldJSONName] = req.ErrMessages()[fieldJSONName][e.ActualTag()]
		}
		ctx.JSON(
			http.StatusUnprocessableEntity,
			res.JSON(false, "Failed to login", &customerror.Err{
				Code:   customerror.CodeErrInvalidRequest,
				Errors: errorMessages,
			}),
		)
		return
	}

	token, err := h.userUsecase.Login(req)
	if err != nil {
		if customErr, ok := err.(*customerror.Err); ok {
			ctx.JSON(http.StatusBadRequest, res.JSON(false, "Failed to login", customErr))
			return
		}
		ctx.JSON(http.StatusInternalServerError, res.JSON(false, "Something went wrong", nil))
		return
	}

	ctx.JSON(http.StatusOK, res.JSON(true, "Success to login", token))
}
