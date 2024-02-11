package db

import (
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
	"webapp/types"
)

type PostgreSQL struct {
	DB *gorm.DB
}

// NewPostgreSQL constructor for PostgreSQL struct.
func NewPostgreSQL(connStr string) *PostgreSQL {
	pg := &PostgreSQL{}
	connected := make(chan bool)
	go pg.connectWithRetry(connStr, connected)

	select {
	case <-connected:
		return pg
	case <-time.After(30 * time.Second):
		log.Println("Timeout while connecting to the database")
		return nil
	}
}

func (p *PostgreSQL) connectWithRetry(connStr string, connected chan bool) {
	var db *gorm.DB
	var err error

	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			p.DB = db
			log.Println("Database connection established")
			// Adding the UUID extension to the database
			db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

			err = db.AutoMigrate(&types.User{})
			if err != nil {
				log.Printf("Error while migrating the table: %v", err)
				return
			}
			log.Println("User table migrated successfully")
			connected <- true
			return
		}
		log.Printf("Attempt %d: could not connect to database: %v", i, err)
		time.Sleep(time.Duration(5) * time.Second)
	}
	close(connected)
	log.Printf("Failed to connect to the database after %d attempts: %v", maxAttempts, err)
}

func (p *PostgreSQL) Ping(ctx *gin.Context) error {
	if p == nil || p.DB == nil {
		return errors.New("DB object is not initialized")
	}

	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) CreateUser(ctx *gin.Context, user types.User) (types.User, error) {
	if p == nil || p.DB == nil {
		return types.User{}, errors.New("DB object is not initialized")
	}

	if err := p.DB.Create(&user).Error; err != nil {
		return types.User{}, err
	}
	return user, nil
}

func (p *PostgreSQL) GetByUsername(ctx *gin.Context, username string) (types.User, error) {
	if p == nil || p.DB == nil {
		return types.User{}, errors.New("DB object is not initialized")
	}

	var user types.User
	if err := p.DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Error getting the user by username : %v", err)
		return types.User{}, err
	}
	return user, nil
}

func (p *PostgreSQL) UpdateUser(ctx *gin.Context, user types.User) (types.User, error) {
	if p == nil || p.DB == nil {
		return types.User{}, errors.New("DB object is not initialized")
	}
	// The user details are stored in the context in the middleware
	u, ok := ctx.Get("user")
	if !ok {
		return types.User{}, errors.New("user not found in context")
	}
	existingUser := u.(types.User)

	if err := p.DB.Model(&existingUser).Updates(user).Error; err != nil {
		return types.User{}, err
	}
	return user, nil
}

func (p *PostgreSQL) Close() error {
	if p == nil || p.DB == nil {
		return errors.New("DB object is not initialized")
	}
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

//func (p *PostgreSQL) Exec(query string, args ...interface{}) error {
//	if p == nil || p.DB == nil {
//		return errors.New("DB object is not initialized")
//	}
//	return p.DB.Exec(query, args...).Error
//}
