package repository

import (
	"github.com/gin-gonic/gin"
	"webapp/types"
)

type UserRepo interface {
	Ping(ctx *gin.Context) error
	CreateUser(ctx *gin.Context, user types.User) (types.User, error)
	GetByUsername(ctx *gin.Context, username string) (types.User, error)
	UpdateUser(ctx *gin.Context, user types.User) (types.User, error)
	Close() error
	MarkEmailAsVerified(ctx *gin.Context, uuid string) error
	//GetByEmailVerificationID(ctx *gin.Context, id string) (types.User, error)
	//Exec(query string, args ...interface{}) error
}
