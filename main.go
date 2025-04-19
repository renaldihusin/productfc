package main

import (
	// golang package
	"context"
	"productfc/cmd/product/handler"
	"productfc/cmd/product/repository"
	"productfc/cmd/product/resource"
	"productfc/cmd/product/service"
	"productfc/cmd/product/usecase"
	"productfc/config"
	"productfc/infrastructure/log"
	"productfc/kafka/consumer"
	"productfc/routes"

	// external package
	"github.com/gin-gonic/gin"
)

// main main.
func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)

	log.SetupLogger()

	productRepository := repository.NewProductRepository(db, redis)
	productService := service.NewProductService(*productRepository)
	productUsecase := usecase.NewProductUsecase(*productService)
	productHandler := handler.NewProductHandler(*productUsecase)

	kafkaProductUpdateStockConsumer := consumer.NewProductUpdateStockConsumer(
		[]string{"localhost:9093"},
		"stock.update",
		*productService,
	)

	kafkaProductUpdateStockConsumer.Start(context.Background())

	kafkaProductRollbackStockConsumer := consumer.NewProductRollbackStockConsumer(
		[]string{"localhost:9093"},
		"stock.rollback",
		*productService,
	)

	kafkaProductRollbackStockConsumer.Start(context.Background())

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *productHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
