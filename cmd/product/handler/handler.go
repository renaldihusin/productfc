package handler

import (
	// golang package
	"fmt"
	"net/http"
	"productfc/cmd/product/usecase"
	"productfc/infrastructure/log"
	"productfc/models"
	"strconv"

	// external package
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	ProductUsecase usecase.ProductUsecase
}

// NewProductHandler new product handler by given ProductUsecase.
//
// It returns pointer of ProductHandler when successful.
// Otherwise, nil pointer of ProductHandler will be returned.
func NewProductHandler(orderUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		ProductUsecase: orderUsecase,
	}
}

// GetProductInfo get product info by given c pointer of gin.Context.
func (h *ProductHandler) GetProductInfo(c *gin.Context) {
	productIDstr := c.Param("id")

	productID, err := strconv.ParseInt(productIDstr, 10, 64)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productIDstr,
		}).Errorf("strconv.ParseInt got error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Product ID",
		})

		return
	}

	product, err := h.ProductUsecase.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productID,
		}).Errorf("h.ProductUsecase.GetProductByID() got error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err,
		})

		return
	}

	if product.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"productID": productID,
		}).Info("Product ID not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Product Not Exists",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// GetProductCategoryInfo get product category info by given c pointer of gin.Context.
func (h *ProductHandler) GetProductCategoryInfo(c *gin.Context) {
	productCategoryIDstr := c.Param("id")

	productCategoryID, err := strconv.Atoi(productCategoryIDstr)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productCategoryIDstr,
		}).Errorf("strconv.Atoi got error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Product ID",
		})

		return
	}

	productCategory, err := h.ProductUsecase.GetProductCategoryByID(c.Request.Context(), productCategoryID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productCategoryID": productCategoryID,
		}).Errorf("h.ProductUsecase.GetProductCategoryByID() got error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err,
		})

		return
	}

	if productCategory.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"productCategoryID": productCategoryID,
		}).Info("Product Category Not Found")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Product Category Not Exists",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"productCategory": productCategory,
	})
}

// ProductManagement product management by given c pointer of gin.Context.
func (h *ProductHandler) ProductManagement(c *gin.Context) {
	var param models.ProductManagementParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Input",
		})

		return
	}

	// validateProductManagementParameter
	if param.Action == "" {
		log.Logger.Error("missing parameter action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing required parameter",
		})

		return
	}

	switch param.Action {
	case "add":
		if param.ID != 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		productID, err := h.ProductUsecase.CreateNewProduct(c.Request.Context(), &param.Product)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.CreateNewProduct() got error %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Sucessfully create new product: %d", productID),
		})

		return
	case "edit":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		product, err := h.ProductUsecase.EditProduct(c.Request.Context(), &param.Product)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.EditProduct() got error %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Success edit product!",
			"product": product,
		})

		return
	case "delete":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		err := h.ProductUsecase.DeleteProduct(c.Request.Context(), param.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.DeleteProduct() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Product %d successfully deleted!", param.ID),
		})
	default:
		log.Logger.Errorf("Invalid action: %s", param.Action)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Action",
		})

		return
	}

}

// ProductCategoryManagement product category management by given c pointer of gin.Context.
func (h *ProductHandler) ProductCategoryManagement(c *gin.Context) {
	var param models.ProductCategoryManagementParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Error(err.Error()) // utk debugging
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Input",
		})
		return
	}

	if param.Action == "" {
		log.Logger.Error("missing parameter action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing required parameter",
		})

		return
	}

	switch param.Action {
	case "add":
		if param.ID != 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product category id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		productCategoryID, err := h.ProductUsecase.CreateNewProductCategory(c.Request.Context(), &param.ProductCategory)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.CreateNewProductCategory got error %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully create new product category: %d", productCategoryID),
		})

		return
	case "edit":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		productCategory, err := h.ProductUsecase.EditProductCategory(c.Request.Context(), &param.ProductCategory)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.EditProductCategory got error %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":         "Success Edit Product",
			"productCategory": productCategory,
		})

		return
	case "delete":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})
			return
		}

		err := h.ProductUsecase.DeleteProductCategory(c.Request.Context(), param.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param, // notes: kalau ada PII --> prevent print log PII data
			}).Errorf("h.ProductUsecase.DeleteProductCategory() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Product Category ID %d successfully deleted!", param.ID),
		})

		return
	default:
		log.Logger.Errorf("Invalid action: %s", param.Action)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Action",
		})

		return
	}
}

// /v1/search/product?name=iphone...
func (h *ProductHandler) SearchProduct(c *gin.Context) {
	name := c.Query("name")
	category := c.Query("category")

	minPrice, _ := strconv.ParseFloat(c.Query("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice"), 64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "2"))

	orderBy := c.Query("orderBy")
	sort := c.Query("sort")

	param := models.SearchProductParameter{
		Name:     name,
		Category: category,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
		Sort:     sort,
	}
	products, totalCount, err := h.ProductUsecase.SearchProduct(c.Request.Context(), param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("h.ProductUsecase.SearchProduct() got error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err,
		})

		return
	}

	// next page url

	totalPages := (totalCount + pageSize - 1) / pageSize

	var nextPageUrl *string
	if page < totalPages {
		url := fmt.Sprintf("%s/v1/product/search?name=%s&category=%s&minPrice=%0.f&maxPrice=%0.f&page=%d&pageSize=%d",
			c.Request.Host, name, category, minPrice, maxPrice, page+1, pageSize)
		nextPageUrl = &url
	}

	c.JSON(http.StatusOK, gin.H{
		// auto sorting a-z
		"data": models.SearchProductResponse{
			Products:    products,
			Page:        page,
			PageSize:    pageSize,
			TotalCount:  totalCount,
			TotalPages:  totalPages,
			NextPageUrl: nextPageUrl,
		},
	})
}
