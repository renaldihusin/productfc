package consumer

import (
	// golang package
	"context"
	"encoding/json"
	"log"
	"productfc/cmd/product/service"
	"productfc/models"

	// external package
	"github.com/segmentio/kafka-go"
)

type ProductRollbackStockConsumer struct {
	Reader         *kafka.Reader
	ProductService service.ProductService
}

// NewProductRollbackStockConsumer new product rollback stock consumer by given slice of brokers, topic, and ProductService.
//
// It returns pointer of ProductRollbackStockConsumer when successful.
// Otherwise, nil pointer of ProductRollbackStockConsumer will be returned.
func NewProductRollbackStockConsumer(brokers []string, topic string, productService service.ProductService) *ProductRollbackStockConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "productfc",
	})

	return &ProductRollbackStockConsumer{
		Reader:         reader,
		ProductService: productService,
	}
}

// Start start.
func (c *ProductRollbackStockConsumer) Start(ctx context.Context) {
	log.Println("[Kafka] Listening to topic 'stock.rollback'")

	for {
		message, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			continue
		}

		var event models.ProductStockUpdateEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			continue
		}

		// looping based on product stock update event
		for _, product := range event.Products {
			// rollback stock
			err = c.ProductService.AddProductStockByProductID(ctx, product.ProductID, product.Qty)
			if err != nil {
				continue
			}
		}
	}
}
