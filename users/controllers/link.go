package controllers

import (
	"ambassador/users/database"
	"ambassador/users/models"
	"ambassador/users/utils"

	"github.com/bxcodec/faker/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateLink(c *fiber.Ctx) error {
	req := models.CreateLinkRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		panic(err)
	}

	link := models.Link{
		UserId: userId,
		Code:   faker.Username(),
	}

	for _, productId := range req.Products {
		product := models.Product{}
		product.Id = uint(productId)
		link.Products = append(link.Products, product)
	}

	database.DB.Create(&link)

	return c.Status(201).JSON(link)
}

func GetStats(c *fiber.Ctx) error {
	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		panic(err)
	}

	links := []models.Link{}

	database.DB.Find(&links, models.Link{
		UserId: userId,
	})

	var result []interface{}

	orders := []models.Order{}

	for _, link := range links {
		database.DB.Preload("OrderItems").Find(&orders, &models.Order{
			Code:     link.Code,
			Complete: true,
		})

		revenue := 0.0
		for _, order := range orders {
			revenue += order.GetTotal()
		}

		result = append(result, fiber.Map{
			"code":    link.Code,
			"count":   len(orders),
			"revenue": revenue,
		})
	}

	return c.Status(200).JSON(result)
}
