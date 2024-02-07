package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
	"net/http"
	"webapp/service"
	"webapp/types"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (uc *userController) CreateUser(ctx *gin.Context) {
	var (
		request types.UserRequest
	)
	// Validating If Request Payload has any unknown fields
	j := json.NewDecoder(ctx.Request.Body)
	j.DisallowUnknownFields()
	err := j.Decode(&request)
	if err != nil {
		log.Printf("Bad Request with error : %v", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}
	// Validating the Content of the Request
	err = request.Validate()
	if err != nil {
		log.Printf("Bad Request with error : %v", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = request.Validate()
	if err != nil {
		log.Printf("Bad Request with apperror : %v", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}

	response, err := uc.userService.CreateUser(ctx, request)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				log.Error("Username Already Exists")
				ctx.Status(http.StatusBadRequest)
				return
			}
		}
		log.Printf("Failed to create user with apperror : %v", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, response)
	return
}

func (uc *userController) GetUser(ctx *gin.Context) {
	var (
		response types.UserResponse
	)
	if ctx.Request.Body != http.NoBody {
		log.Println("Request has a payload")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Request query params validation
	if len(ctx.Request.URL.RawQuery) > 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userDetails := user.(types.User)
	response = types.UserResponse{
		Username:  userDetails.Username,
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		CreatedAt: userDetails.CreatedAt,
		UpdatedAt: userDetails.UpdatedAt,
		ID:        userDetails.ID.String(),
	}

	ctx.JSON(http.StatusOK, response)
	return
}

func (uc *userController) UpdateUser(ctx *gin.Context) {
	var (
		request types.UpdateUserRequest
	)
	// Validating If Request Payload has any unknown fields
	j := json.NewDecoder(ctx.Request.Body)
	j.DisallowUnknownFields()
	err := j.Decode(&request)
	if err != nil {
		log.Printf("Bad Request with error : %v", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}
	// Validating the Content of the Request
	err = request.Validate()
	if err != nil {
		log.Printf("Bad Request with error : %v", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}

	_, err = uc.userService.UpdateUser(ctx, request)
	if err != nil {
		log.Printf("Failed to update user with error : %v", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
	return
}
