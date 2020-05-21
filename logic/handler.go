package logic

import (
	"errors"

	"github.com/rahulbharuka/train-route-finder/repository"

	"github.com/gin-gonic/gin"
)

var errRouteNotFound = errors.New("no route exist")

// Handler is the logic handler interface
type Handler interface {
	Routes(ctx *gin.Context)
}

// handlerImpl is a implementation of Handler interface
type handlerImpl struct {
	repo repository.Handler
}

// GetHandler initializes and returns the logic layer handler.
func GetHandler() Handler {
	return &handlerImpl{
		repo: repository.GetHandler(),
	}
}

// handlerError is a helper function to return JSON error.
func handlerError(ctx *gin.Context, errCode int, err error) {
	ctx.JSON(errCode, gin.H{
		"message": err.Error(),
	})
}
