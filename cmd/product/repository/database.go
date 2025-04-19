package repository

import (
	// golang package
	"context"
	"errors"
	"fmt"
	"productfc/models"

	// external package
	"gorm.io/gorm"
)

// FindProductByID find product by id by given productID.
//
// It returns pointer of models.Product, and nil error when successful.
// Otherwise, nil pointer of models.Product, and error will be returned.
func (r *ProductRepository) FindProductByID(ctx context.Context, productID int64) (*models.Product, error) {
	var product models.Product
	err := r.Database.WithContext(ctx).Table("product").Where("id = ?", productID).Last(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Product{}, nil
		}

		return nil, err
	}

	return &product, nil
}

// FindProductCategoryByID find product category by id by given productCategoryID.
//
// It returns pointer of models.ProductCategory, and nil error when successful.
// Otherwise, nil pointer of models.ProductCategory, and error will be returned.
func (r *ProductRepository) FindProductCategoryByID(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	var productCategory models.ProductCategory
	err := r.Database.WithContext(ctx).Table("product_category").Where("id = ?", productCategoryID).Last(&productCategory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.ProductCategory{}, nil
		}

		return nil, err
	}

	return &productCategory, nil
}

// InsertNewProduct insert new product by given product pointer of models.Product.
//
// It returns int64, and nil error when successful.
// Otherwise, empty int64, and error will be returned.
func (r *ProductRepository) InsertNewProduct(ctx context.Context, product *models.Product) (int64, error) {
	err := r.Database.WithContext(ctx).Table("product").Create(product).Error
	if err != nil {
		return 0, err
	}

	return product.ID, nil
}

// InsertNewProductCategory insert new product category by given productCategory pointer of models.ProductCategory.
//
// It returns int, and nil error when successful.
// Otherwise, empty int, and error will be returned.
func (r *ProductRepository) InsertNewProductCategory(ctx context.Context, productCategory *models.ProductCategory) (int, error) {
	err := r.Database.WithContext(ctx).Table("product_category").Create(productCategory).Error
	if err != nil {
		return 0, err
	}

	return productCategory.ID, nil
}

// UpdateProductStockByProductID update product stock by product id by given productID, and qty.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (r *ProductRepository) DeductProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	err := r.Database.Table("product").WithContext(ctx).Model(&models.Product{}).
		Updates(map[string]interface{}{
			"stock": gorm.Expr("stock - %d", qty),
		}).Where("id = ?", productID).Error
	if err != nil {
		return err
	}

	return nil
}

// AddProductStockByProductID add product stock by product id by given productID, and qty.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (r *ProductRepository) AddProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	err := r.Database.Table("product").WithContext(ctx).Model(&models.Product{}).
		Updates(map[string]interface{}{
			"stock": gorm.Expr("stock + %d", qty),
		}).Where("id = ?", productID).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateProduct update product by given product pointer of models.Product.
//
// It returns pointer of models.Product, and nil error when successful.
// Otherwise, nil pointer of models.Product, and error will be returned.
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	err := r.Database.WithContext(ctx).Table("product").Save(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil // updated data
}

// UpdateProductCategory update product category by given productCategory pointer of models.ProductCategory.
//
// It returns pointer of models.ProductCategory, and nil error when successful.
// Otherwise, nil pointer of models.ProductCategory, and error will be returned.
func (r *ProductRepository) UpdateProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	err := r.Database.WithContext(ctx).Table("product_category").Save(productCategory).Error
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

// DeleteProduct delete product by given productID.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (r *ProductRepository) DeleteProduct(ctx context.Context, productID int64) error {
	err := r.Database.WithContext(ctx).Table("product").Delete(&models.Product{}, productID).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteProductCategory delete product category by given productCategoryID.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (r *ProductRepository) DeleteProductCategory(ctx context.Context, productCategoryID int) error {
	err := r.Database.WithContext(ctx).Table("product_category").Delete(&models.ProductCategory{}, productCategoryID).Error
	if err != nil {
		return err
	}

	return nil
}

// SearchProduct search product by given SearchProductParameter.
//
// It returns slice of models.Product, int, and nil error when successful.
// Otherwise, nil value of models.Product slice, empty int, and error will be returned.
func (r *ProductRepository) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	var products []models.Product
	var totalCount int64

	query := r.Database.WithContext(ctx).Table("product").
		Select("product.id, product.name, product.description, product.price, product.stock, product.category_id, product_category.name AS category").
		Joins("JOIN product_category ON product.category_id = product_category.id")

	// filtering
	if param.Name != "" { // iphone --> iphone X, etc.
		query = query.Where("product.name ILIKE ?", "%"+param.Name+"%")
	}

	if param.Category != "" {
		query = query.Where("product_category.name = ?", param.Category)
	}

	if param.MinPrice > 0 {
		query = query.Where("product.price >= ?", param.MinPrice)
	}

	if param.MaxPrice > 0 {
		query = query.Where("product.price <= ?", param.MaxPrice)
	}

	// pagination

	// dapetin total counts dari hasil query
	query.Model(&models.Product{}).Count(&totalCount)

	// default order by
	if param.OrderBy == "" {
		param.OrderBy = "product.name"
	}

	if param.Sort == "" || (param.Sort != "ASC" && param.Sort != "DESC") {
		param.Sort = "ASC"
	}

	orderBy := fmt.Sprintf("%s %s", param.OrderBy, param.Sort)
	query = query.Order(orderBy)

	offset := (param.Page - 1) * param.PageSize
	query = query.Limit(param.PageSize).Offset(offset)

	err := query.Scan(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, int(totalCount), nil
}
