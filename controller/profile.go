package controller

import (
	"KaloriKu/repository"

	"github.com/gofiber/fiber/v2"
	// "KaloriKu/model"
)

func GetAllProfile(c *fiber.Ctx) error {
	// Ambil token dari header
	token := c.Get("Auth")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Ambil data profile dari database
	profiles, err := repository.GetAllProfile()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching profiles")
	}

	// Return response sukses dengan data profile
	return c.JSON(fiber.Map{
		"profiles": profiles,
	})
}

func GetProfileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	profile, err := repository.GetProfileByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Profile not found")
	}
	return c.JSON(profile)
}
