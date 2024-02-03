package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	GetHealth(context *gin.Context)
}
