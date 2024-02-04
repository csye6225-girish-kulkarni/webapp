package service

import (
	"Health-Check/repository"
	"Health-Check/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) CreateUser(ctx *gin.Context, userRequest types.UserRequest) (types.UserResponse, error) {
	var (
		user types.User
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Error hashing the password : %v", err)
		return types.UserResponse{}, err
	}
	user = types.User{
		Username:  userRequest.Username,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Password:  string(hashedPassword),
	}

	updatedUser, err := us.repo.Create(ctx, user)
	if err != nil {
		log.Errorf("Error creating the user : %v", err)
		return types.UserResponse{}, err
	}

	log.Println("User created successfully")
	return types.UserResponse{
		Username:  updatedUser.Username,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		ID:        updatedUser.ID.String(),
	}, nil
}
