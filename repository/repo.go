package repository

import "github.com/gin-gonic/gin"

type UserRepo interface {
	Ping(ctx *gin.Context) error
}
