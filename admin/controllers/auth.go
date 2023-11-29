package controllers

import (
	"ambassador/admin/database"
	"ambassador/admin/models"
	"ambassador/admin/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	req := models.RegisterRequest{}

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

	user := models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		IsAmbassador: false,
	}

	if err := user.SetPassword(req.Password); err != nil {
		panic(err)
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(409).JSON(fiber.Map{
			"message": "This Email already in used.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Register Successfully!",
	})
}

func Login(c *fiber.Ctx) error {
	req := models.LoginRequest{}

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

	user := models.User{}

	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Email or Password is invalid.",
		})
	}

	if err := user.VerifyPassword(req.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Email or Password is invalid.",
		})
	}

	accessToken, err := utils.GenerateAccessToken(user.Id)
	if err != nil {
		panic(err)
	}

	utils.SetCookie(c, "access_token", accessToken, time.Now().Add(time.Hour*24))

	return c.Status(200).JSON(fiber.Map{
		"message": "Login Successfully!",
	})
}

func User(c *fiber.Ctx) error {
	user := models.User{}

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		panic(err)
	}

	database.DB.Where("id = ?", userId).First(&user)

	return c.Status(200).JSON(user)
}

func Logout(c *fiber.Ctx) error {
	utils.ClearCookie(c, "access_token")

	return c.Status(200).JSON(fiber.Map{
		"message": "Logout Successfully!",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	req := models.UpdateInfoRequest{}

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

	user := models.User{
		Id:        userId,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	if result := database.DB.Model(&user).Updates(&user); result.Error != nil {
		return c.Status(409).JSON(fiber.Map{
			"message": "This Email is already in used.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Update User Info Successfully!",
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	req := models.UpdatePasswordRequest{}

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

	user := models.User{
		Id: userId,
	}

	user.SetPassword(req.Password)

	if result := database.DB.Model(&user).Updates(&user); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Update User Info Successfully!",
	})
}
