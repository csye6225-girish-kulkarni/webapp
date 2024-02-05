package controller

import (
	"Health-Check/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	ctx.Header("cache-control", "no-cache")
	// Request Payload validation
	if ctx.Request.ContentLength > 0 {
		log.Println("Request has a payload")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Request query params validation
	if len(ctx.Request.URL.RawQuery) > 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := hs.healthService.Ping(ctx)
	if err != nil {
		log.Printf("Unable to Ping to DB err : %v", err)
		ctx.Status(http.StatusServiceUnavailable)
		return
	}
	log.Println("Database Successfully pinged")
	ctx.Status(http.StatusOK)
	return
}
