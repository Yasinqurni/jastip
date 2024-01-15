package recovery

import (
	"fmt"
	"jastip-app/pkg/logger"
	"jastip-app/pkg/res"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recover(ctx *gin.Context) {
	if r := recover(); r != nil {
		logger.L().Error(fmt.Sprintf("Recovered by error : %v", r))
		ctx.JSON(http.StatusInternalServerError, res.JSON(false, "Something went wrong", nil))
	}
}
