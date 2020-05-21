package logic

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/rahulbharuka/train-route-finder/repository"

	"github.com/gin-gonic/gin"
)

// FindRealtimeRoute finds route(s) from source to destination.
func (h *handlerImpl) Routes(ctx *gin.Context) {
	source := ctx.Query("src")
	destination := ctx.Query("dst")

	if source == destination {
		log.Println("source and destination cannot be same")
		handlerError(ctx, http.StatusBadRequest, errors.New("source and destination cannot be same"))
		return
	}

	var err error
	var journeyTime time.Time
	var computeTimeCost bool
	jTime := ctx.Query("journeyTime")
	if jTime != "" {
		journeyTime, err = time.Parse("2006-01-02T15:04", jTime)
		if err != nil {
			log.Println("invalid journey start time")
			handlerError(ctx, http.StatusBadRequest, errors.New("invalid journey start time"))
			return
		}
		computeTimeCost = true
	}

	resp, err := h.repo.FindRoutes(source, destination, journeyTime, computeTimeCost)
	if err == repository.ErrInvalidRequest {
		handlerError(ctx, http.StatusBadRequest, err)
		return
	}
	if err == repository.ErrRouteNotFound {
		handlerError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
