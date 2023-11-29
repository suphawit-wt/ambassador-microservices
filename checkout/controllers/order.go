package controllers

import (
	"ambassador/checkout/database"
	"ambassador/checkout/models"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

func CreateOrder(c *fiber.Ctx) error {
	req := models.CreateOrderRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	link := models.Link{
		Code: req.Code,
	}

	if result := database.DB.Preload("User").First(&link); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found.",
		})
	}

	order := models.Order{
		Code:            link.Code,
		UserId:          link.UserId,
		AmbassadorEmail: link.User.Email,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		Address:         req.Address,
		Country:         req.Country,
		City:            req.City,
		Zip:             req.Zip,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, reqProduct := range req.Products {
		product := models.Product{}
		product.Id = uint(reqProduct["product_id"])
		database.DB.First(&product)

		total := product.Price * float64(reqProduct["quantity"])

		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(reqProduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(product.Title),
					Description: stripe.String(product.Description),
					Images:      []*string{stripe.String(product.Image)},
				},
				UnitAmount: stripe.Int64(100 * int64(product.Price)),
				Currency:   stripe.String("usd"),
			},
			Quantity: stripe.Int64(int64(reqProduct["quantity"])),
		})
	}

	stripe.Key = os.Getenv("STRIPE_KEY")

	params := stripe.CheckoutSessionParams{
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String("http://localhost:5000/success?source={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("http://localhost:5000/error"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
	}

	source, err := session.New(&params)
	if err != nil {
		tx.Rollback()
		log.Printf("session.New: %v", err)
	}

	order.TransactionId = source.ID

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	tx.Commit()

	return c.Status(200).JSON(source)
}

func CompleteOrder(c *fiber.Ctx) error {
	req := models.ConfirmOrderRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	order := models.Order{}

	if result := database.DB.Preload("OrderItems").First(&order, models.Order{
		TransactionId: req.Source,
	}); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found.",
		})
	}

	order.Complete = true
	database.DB.Save(&order)

	go func(order models.Order) {
		ambassadorRevenue := 0.0
		adminRevenue := 0.0

		for _, item := range order.OrderItems {
			ambassadorRevenue += item.AmbassadorRevenue
			adminRevenue += item.AdminRevenue
		}

		user := models.User{}
		user.Id = order.Id

		database.DB.First(&user)

		database.RedisClient.ZIncrBy(context.Background(), "rankings", ambassadorRevenue, user.Name())

		producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_SERVERS")})
		if err != nil {
			panic(err)
		}

		defer producer.Close()

		topic := os.Getenv("KAFKA_TOPIC")

		message := models.EmailKafkaMessage{
			Id:                order.Id,
			AmbassadorRevenue: ambassadorRevenue,
			AdminRevenue:      adminRevenue,
			Code:              order.Code,
			AmbassadorEmail:   order.AmbassadorEmail,
		}

		messageMarshal, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          messageMarshal,
		}, nil)

		producer.Flush(15 * 1000)
	}(order)

	return c.Status(200).JSON(fiber.Map{
		"message": "Order Confirm Successfully!",
	})
}
