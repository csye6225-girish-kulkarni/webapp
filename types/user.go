package types

import (
	"github.com/gofrs/uuid"
	"time"
)

type UserRequest struct {
	Username  string `json:"username" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=20"`
	FirstName string `json:"firstName" validate:"required,alpha"`
	LastName  string `json:"lastName" validate:"required,alpha"`
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"type:varchar(255);unique"`
	Password  string `gorm:"type:varchar(255)"`
	FirstName string `gorm:"type:varchar(255)"`
	LastName  string `gorm:"type:varchar(255)"`
}

type UserResponse struct {
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
