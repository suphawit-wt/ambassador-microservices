package controllers

import (
	"ambassador/users/database"
	"ambassador/users/models"
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetProductsFrontend(c *fiber.Ctx) error {
	products := []models.Product{}
	ctx := context.Background()

	productsCache, err := database.RedisClient.Get(ctx, "products_frontend").Result()
	if err != nil {
		database.DB.Find(&products)

		productsBytes, err := json.Marshal(products)
		if err != nil {
			panic(err)
		}

		database.RedisClient.Set(ctx, "products_frontend", productsBytes, time.Minute*5)
	} else {
		json.Unmarshal([]byte(productsCache), &products)
	}

	return c.Status(200).JSON(products)
}

func GetProductsBackend(c *fiber.Ctx) error {
	products := []models.Product{}
	ctx := context.Background()

	productsCache, err := database.RedisClient.Get(ctx, "products_backend").Result()
	if err != nil {
		database.DB.Find(&products)

		productsBytes, err := json.Marshal(products)
		if err != nil {
			panic(err)
		}

		database.RedisClient.Set(ctx, "products_backend", productsBytes, time.Minute*5)
	} else {
		json.Unmarshal([]byte(productsCache), &products)
	}

	var searchedProducts []models.Product

	if s := c.Query("s"); s != "" {
		srcLower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), srcLower) || strings.Contains(strings.ToLower(product.Description), srcLower) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		searchedProducts = products
	}

	if sortQuery := c.Query("sort"); sortQuery != "" {
		sortLower := strings.ToLower(sortQuery)
		if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	var total = len(searchedProducts)
	var data []models.Product
	perPage := 9

	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchedProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		data = searchedProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}

	return c.Status(200).JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": total/perPage + 1,
	})
}
