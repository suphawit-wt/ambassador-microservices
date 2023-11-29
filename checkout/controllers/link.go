package controllers

import (
	"ambassador/checkout/database"
	"ambassador/checkout/models"

	"github.com/gofiber/fiber/v2"
)

func GetLink(c *fiber.Ctx) error {
	code := c.Params("code")

	link := models.Link{
		Code: code,
	}

	database.DB.Preload("User").Preload("Products").First(&link)

	return c.Status(200).JSON(link)
}
