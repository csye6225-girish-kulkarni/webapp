package service

import (
	"github.com/gin-gonic/gin"
	"webapp/repository"
)

type Service interface {
	Ping(ctx *gin.Context) error
}

type HealthService struct {
	repo repository.UserRepo
}

func NewHealthService(repo repository.UserRepo) *HealthService {
	return &HealthService{
		repo: repo,
	}
}

func (hs *HealthService) Ping(ctx *gin.Context) error {
	err := hs.repo.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
