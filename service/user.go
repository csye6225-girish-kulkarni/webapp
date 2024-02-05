package service

import (
	"Health-Check/repository"
	"Health-Check/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepo
}

type UserService interface {
	CreateUser(ctx *gin.Context, userRequest types.UserRequest) (types.UserResponse, error)
	ValidateUser(ctx *gin.Context, username, password string) (bool, types.User, error)
	GetUserByUsername(ctx *gin.Context, username string) (types.User, error)
	UpdateUser(ctx *gin.Context, userRequest types.UserRequest) (types.UserResponse, error)
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) CreateUser(ctx *gin.Context, userRequest types.UserRequest) (types.UserResponse, error) {
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

func (us *userService) ValidateUser(ctx *gin.Context, username, password string) (bool, types.User, error) {
	user, err := us.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Errorf("Error getting the user by username : %v", err)
		return false, types.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Errorf("Error comparing the password : %v", err)
		return false, types.User{}, err
	}

	return true, user, nil
}

func (us *userService) GetUserByUsername(ctx *gin.Context, username string) (types.User, error) {
	user, err := us.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Errorf("Error getting the user by username : %v", err)
		return types.User{}, err
	}
	return user, nil
}

func (us *userService) UpdateUser(ctx *gin.Context, userRequest types.UserRequest) (types.UserResponse, error) {
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
