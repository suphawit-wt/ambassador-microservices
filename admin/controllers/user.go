package controllers

import (
	"ambassador/admin/database"
	"ambassador/admin/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllAmbassador(c *fiber.Ctx) error {
	users := []models.User{}

	database.DB.Where("is_ambassador = true").Find(&users)

	return c.Status(200).JSON(users)
}
