package usecase

import (
	// golang package
	"context"
	"productfc/cmd/product/service"
	"productfc/infrastructure/log"
	"productfc/models"

	// external package
	"github.com/sirupsen/logrus"
)

type ProductUsecase struct {
	ProductService service.ProductService
}

// NewProductUsecase new product usecase by given ProductService.
//
// It returns pointer of ProductUsecase when successful.
// Otherwise, nil pointer of ProductUsecase will be returned.
func NewProductUsecase(orderService service.ProductService) *ProductUsecase {
	return &ProductUsecase{
		ProductService: orderService,
	}
}

// GetProductByID get product by id by given productID.
//
// It returns pointer of models.Product, and nil error when successful.
// Otherwise, nil pointer of models.Product, and error will be returned.
func (uc *ProductUsecase) GetProductByID(ctx context.Context, productID int64) (*models.Product, error) {
	product, err := uc.ProductService.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductCategoryByID get product category by id by given productCategoryID.
//
// It returns pointer of models.ProductCategory, and nil error when successful.
// Otherwise, nil pointer of models.ProductCategory, and error will be returned.
func (uc *ProductUsecase) GetProductCategoryByID(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	productCategory, err := uc.ProductService.GetProductCategoryByID(ctx, productCategoryID)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

// CreateNewProduct create new product by given param pointer of models.Product.
//
// It returns int64, and nil error when successful.
// Otherwise, empty int64, and error will be returned.
func (uc *ProductUsecase) CreateNewProduct(ctx context.Context, param *models.Product) (int64, error) {
	productID, err := uc.ProductService.CreateNewProduct(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name":     param.Name,
			"category": param.CategoryID,
		}).Errorf("uc.ProductService.CreateNewProduct got error %v", err)
		return 0, err
	}

	return productID, nil
}

// CreateNewProductCategory create new product category by given param pointer of models.ProductCategory.
//
// It returns int, and nil error when successful.
// Otherwise, empty int, and error will be returned.
func (uc *ProductUsecase) CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int, error) {
	productCategoryID, err := uc.ProductService.CreateNewProductCategory(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name": param.Name,
		}).Errorf("uc.ProductService.CreateNewProductCategory got error %v", err)
		return 0, err
	}

	return productCategoryID, nil
}

// EditProduct edit product by given param pointer of models.Product.
//
// It returns pointer of models.Product, and nil error when successful.
// Otherwise, nil pointer of models.Product, and error will be returned.
func (uc *ProductUsecase) EditProduct(ctx context.Context, param *models.Product) (*models.Product, error) {
	product, err := uc.ProductService.EditProdut(ctx, param)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// EditProductCategory edit product category by given param pointer of models.ProductCategory.
//
// It returns pointer of models.ProductCategory, and nil error when successful.
// Otherwise, nil pointer of models.ProductCategory, and error will be returned.
func (uc *ProductUsecase) EditProductCategory(ctx context.Context, param *models.ProductCategory) (*models.ProductCategory, error) {
	productCategory, err := uc.ProductService.EditProductCategory(ctx, param)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

// DeleteProduct delete product by given productID.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (uc *ProductUsecase) DeleteProduct(ctx context.Context, productID int64) error {
	err := uc.ProductService.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProductCategory delete product category by given productCategoryID.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (uc *ProductUsecase) DeleteProductCategory(ctx context.Context, productCategoryID int) error {
	err := uc.ProductService.DeleteProductCategory(ctx, productCategoryID)
	if err != nil {
		return err
	}

	return nil
}

// SearchProduct search product by given SearchProductParameter.
//
// It returns slice of models.Product, int, and nil error when successful.
// Otherwise, nil value of models.Product slice, empty int, and error will be returned.
func (uc *ProductUsecase) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	products, totalCount, err := uc.ProductService.SearchProduct(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}
