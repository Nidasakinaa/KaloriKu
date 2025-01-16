package controller

import (
	"KaloriKu/repository"

	"KaloriKu/model"

	"github.com/gofiber/fiber/v2"
	// "KSI-BE/config"
)

func CreateUser(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	userID, err := repository.CreateUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving user")
	}

	return c.JSON(fiber.Map{
		"message": "User created successfully",
		"user_id": userID,
	})
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := repository.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}
	return c.JSON(user)
}

func GetAllUser(c *fiber.Ctx) error {
	users, err := repository.GetAllUser()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching users")
	}
	return c.JSON(users)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	userID, err := repository.UpdateUser(id, &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating user")
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
		"user_id": userID,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := repository.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting user")
	}
	return c.SendString("User deleted successfully")
}
