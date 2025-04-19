package repository

import (
	// external package
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductRepository struct {
	Database *gorm.DB
	Redis    *redis.Client
}

// NewProductRepository new order repository by given db pointer of gorm.DB, and redis pointer of redis.Client.
//
// It returns pointer of ProductRepository when successful.
// Otherwise, nil pointer of ProductRepository will be returned.
func NewProductRepository(db *gorm.DB, redis *redis.Client) *ProductRepository {
	return &ProductRepository{
		Database: db,
		Redis:    redis,
	}
}
