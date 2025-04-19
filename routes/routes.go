package routes

import (
	// golang package
	"productfc/cmd/product/handler"
	"productfc/middleware"

	// external package
	"github.com/gin-gonic/gin"
)

// SetupRoutes setup routes by given router pointer of gin.Engine, and ProductHandler.
func SetupRoutes(router *gin.Engine, orderHandler handler.ProductHandler) {
	router.Use(middleware.RequestLogger())
	router.POST("/v1/product", orderHandler.ProductManagement)
	router.POST("/v1/product_category", orderHandler.ProductCategoryManagement)

	router.GET("/v1/product/:id", orderHandler.GetProductInfo)
	router.GET("/v1/product_category/:id", orderHandler.GetProductCategoryInfo)

	router.GET("/v1/product/search", orderHandler.SearchProduct)
}
