package service

import (
	// golang package
	"context"
	"productfc/cmd/product/repository"
	"productfc/infrastructure/log"
	"productfc/models"

	// external package
	"github.com/sirupsen/logrus"
)

type ProductService struct {
	ProductRepository repository.ProductRepository
}

// NewProductService new product service by given ProductRepository.
//
// It returns pointer of ProductService when successful.
// Otherwise, nil pointer of ProductService will be returned.
func NewProductService(productRepository repository.ProductRepository) *ProductService {
	return &ProductService{
		ProductRepository: productRepository,
	}
}

// DeductProductStockByProductID deduct product stock by product id by given productID, and qty.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (s *ProductService) DeductProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	err := s.ProductRepository.DeductProductStockByProductID(ctx, productID, qty)
	if err != nil {
		return err
	}

	return nil
}

// AddProductStockByProductID add product stock by product id by given productID, and qty.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (s *ProductService) AddProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	err := s.ProductRepository.AddProductStockByProductID(ctx, productID, qty)
	if err != nil {
		return err
	}

	return nil
}

// di layer service
// kita akan tentukan mau menggunakan resource yg mana
// db or redis

func (s *ProductService) GetProductByID(ctx context.Context, productID int64) (*models.Product, error) {
	// get from Redis
	product, err := s.ProductRepository.GetProductByIDFromRedis(ctx, productID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productID,
		}).Errorf("s.ProductRepository.GetProductByIDFromRedis() got error %v", err)
	}

	if product.ID != 0 {
		return product, nil
	}

	// get from DB
	product, err = s.ProductRepository.FindProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	ctxConcurrent := context.WithValue(ctx, context.Background(), ctx.Value("request_id"))
	go func(ctx context.Context, product *models.Product, productID int64) {
		errConcurrent := s.ProductRepository.SetProductByID(ctx, product, productID)
		if errConcurrent != nil {
			log.Logger.WithFields(logrus.Fields{
				"product": product,
			}).Errorf("s.ProductRepository.SetProductByID() got error %v", errConcurrent)
		}
	}(ctxConcurrent, product, productID)

	return product, nil
}

// GetProductCategoryByID get product category by id by given productCategoryID.
//
// It returns pointer of models.ProductCategory, and nil error when successful.
// Otherwise, nil pointer of models.ProductCategory, and error will be returned.
func (s *ProductService) GetProductCategoryByID(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	productCategory, err := s.ProductRepository.FindProductCategoryByID(ctx, productCategoryID)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

// CreateNewProduct create new product by given param pointer of models.Product.
//
// It returns int64, and nil error when successful.
// Otherwise, empty int64, and error will be returned.
func (s *ProductService) CreateNewProduct(ctx context.Context, param *models.Product) (int64, error) {
	productID, err := s.ProductRepository.InsertNewProduct(ctx, param)
	if err != nil {
		return 0, err
	}

	return productID, nil
}

// CreateNewProductCategory create new product category by given param pointer of models.ProductCategory.
//
// It returns int, and nil error when successful.
// Otherwise, empty int, and error will be returned.
func (s *ProductService) CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int, error) {
	productCategoryID, err := s.ProductRepository.InsertNewProductCategory(ctx, param)
	if err != nil {
		return 0, err
	}

	return productCategoryID, nil
}

// EditProdut edit produt by given product pointer of models.Product.
//
// It returns pointer of models.Product, and nil error when successful.
// Otherwise, nil pointer of models.Product, and error will be returned.
func (s *ProductService) EditProdut(ctx context.Context, product *models.Product) (*models.Product, error) {
	product, err := s.ProductRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// EditProductCategory edit product category by given productCategory pointer of models.ProductCategory.
//
// It returns pointer of models.ProductCategory, and nil error when successful.
// Otherwise, nil pointer of models.ProductCategory, and error will be returned.
func (s *ProductService) EditProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	productCategory, err := s.ProductRepository.UpdateProductCategory(ctx, productCategory)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

// DeleteProduct delete product by given productID.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (s *ProductService) DeleteProduct(ctx context.Context, productID int64) error {
	err := s.ProductRepository.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProductCategory delete product category by given productCategoryID.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func (s *ProductService) DeleteProductCategory(ctx context.Context, productCategoryID int) error {
	err := s.ProductRepository.DeleteProductCategory(ctx, productCategoryID)
	if err != nil {
		return err
	}

	return nil
}

// SearchProduct search product by given SearchProductParameter.
//
// It returns slice of models.Product, int, and nil error when successful.
// Otherwise, nil value of models.Product slice, empty int, and error will be returned.
func (s *ProductService) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	products, totalCount, err := s.ProductRepository.SearchProduct(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}
