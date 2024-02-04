package controller

import (
	"Health-Check/service"
	"Health-Check/types"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var request types.UserRequest
	err := ctx.ShouldBindBodyWith(&request, binding.JSON)
	if err != nil {
		log.Printf("Bad Request with error : %v", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}

	response, err := uc.userService.CreateUser(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
			return
		}
		log.Printf("Failed to create user with error : %v", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, response)
	return
}
