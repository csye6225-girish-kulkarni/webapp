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
)

type PostgreSQL struct {
	DB *gorm.DB
}

// NewPostgreSQL constructor for PostgreSQL struct.
func NewPostgreSQL(connStr string) *PostgreSQL {
	pg := &PostgreSQL{}
	go pg.connectWithRetry(connStr)
	return pg
}

func (pg *PostgreSQL) connectWithRetry(connStr string) {
	var db *gorm.DB
	var err error

	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			pg.DB = db
			log.Println("Database connection established")
			return
		}
		log.Printf("Attempt %d: could not connect to database: %v", i, err)
		time.Sleep(time.Duration(5) * time.Second)
	}

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
