package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"webapp/service"
)

type HealthController interface {
	GetHealth(ctx *gin.Context)
}

type healthController struct {
	healthService service.Service
}

func NewHealthController(hs service.Service) HealthController {
	return &healthController{
		healthService: hs,
	}
}

func (hs *healthController) GetHealth(ctx *gin.Context) {
	// Request Payload validation
	if ctx.Request.ContentLength > 0 {
		log.Info().Msg("Request has a payload")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Request query params validation
	if len(ctx.Request.URL.RawQuery) > 0 {
		log.Info().Msg("Request has query parameters")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := hs.healthService.Ping(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Ping to DB")
		ctx.Status(http.StatusServiceUnavailable)
		return
	}
	log.Info().Msg("Database Successfully pinged")
	ctx.Status(http.StatusOK)
	return
}
