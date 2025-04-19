package repository

import (
	// golang package
	"context"
	"encoding/json"
	"fmt"
	"productfc/models"
	"time"

	// external package
	"github.com/redis/go-redis/v9"
)

var (
	cacheKeyProductInfo         = "product:%d" // format: product:{productID} product:1
	cacheKeyProductCategoryInfo = "product_category:%d"
)

// GetProductByIDFromRedis get product by id from redis by given productID.
//
// It returns pointer of models.Product, and nil error when successful.
// Otherwise, nil pointer of models.Product, and error will be returned.
func (r *ProductRepository) GetProductByIDFromRedis(ctx context.Context, productID int64) (*models.Product, error) {
	cacheKey := fmt.Sprintf(cacheKeyProductInfo, productID)

	var product models.Product
	productStr, err := r.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return &models.Product{}, nil
		}

		return nil, err
	}

	// unmarshal
	err = json.Unmarshal([]byte(productStr), &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetProductCategoryByIDFromRedis get product category by id from redis by given productCategoryID.
//
// It returns pointer of models.ProductCategory, and nil error when successful.
// Otherwise, nil pointer of models.ProductCategory, and error will be returned.
func (r *ProductRepository) GetProductCategoryByIDFromRedis(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	cacheKey := fmt.Sprintf(cacheKeyProductCategoryInfo, productCategoryID)

	var productCategory models.ProductCategory
	productCategoryStr, err := r.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return &models.ProductCategory{}, nil
		}

		return nil, err
	}

	err = json.Unmarshal([]byte(productCategoryStr), &productCategory)
	if err != nil {
		return nil, err
	}

	return &productCategory, nil
}

// SetProductByID set product by id by given product pointer of models.Product, and productID.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (r *ProductRepository) SetProductByID(ctx context.Context, product *models.Product, productID int64) error {
	cacheKey := fmt.Sprintf(cacheKeyProductInfo, productID)
	fmt.Println(cacheKey)
	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	err = r.Redis.SetEx(ctx, cacheKey, productJSON, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

// SetProductCategoryByID set product category by id by given productCategory pointer of models.ProductCategory, and productCategoryID.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (r *ProductRepository) SetProductCategoryByID(ctx context.Context, productCategory *models.ProductCategory, productCategoryID int) error {
	cacheKey := fmt.Sprintf(cacheKeyProductCategoryInfo, productCategoryID)

	productCategoryJSON, err := json.Marshal(productCategory)
	if err != nil {
		return err
	}

	err = r.Redis.SetEx(ctx, cacheKey, productCategoryJSON, 1*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}
