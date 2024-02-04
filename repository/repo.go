package repository

import (
	"Health-Check/types"
	"github.com/gin-gonic/gin"
)

type UserRepo interface {
	Ping(ctx *gin.Context) error
	Create(ctx *gin.Context, user types.User) (types.User, error)
}
