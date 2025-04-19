package consumer

import (
	// golang package
	"context"
	"encoding/json"
	"productfc/cmd/product/service"
	"productfc/infrastructure/log"
	"productfc/models"

	// external package
	"github.com/segmentio/kafka-go"
)

type ProductUpdateStockConsumer struct {
	Reader         *kafka.Reader
	productService service.ProductService
}

// NewProductUpdateStockConsumer new product update stock consumer by given slice of brokers, topic, and ProductService.
//
// It returns pointer of ProductUpdateStockConsumer when successful.
// Otherwise, nil pointer of ProductUpdateStockConsumer will be returned.
func NewProductUpdateStockConsumer(brokers []string, topic string, productService service.ProductService) *ProductUpdateStockConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "productfc",
	})

	return &ProductUpdateStockConsumer{
		productService: productService,
		Reader:         reader,
	}
}

// Start start.
func (c *ProductUpdateStockConsumer) Start(ctx context.Context) {
	log.Logger.Println("[KAFKA] Listening to topic stock.update")

	for {
		message, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Logger.Println("[KAFKA] Error ReadMessage: ", err)
			continue
		}

		// unmarshal event to product update stock struct
		var event models.ProductStockUpdateEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Logger.Println("[KAFKA] Error Unmarshal Event Message: ", err)
			continue
		}

		// update stock
		for _, product := range event.Products {
			err = c.productService.DeductProductStockByProductID(ctx, product.ProductID, product.Qty)
			if err != nil {
				log.Logger.Printf("[KAFKA] Error Update Product Stock Product ID #%d", product.ProductID)
				continue
			}
		}
	}
}
