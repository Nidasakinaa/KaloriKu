package controller

import (
	"KaloriKu/model"
	"KaloriKu/repository"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMenuItem(c *fiber.Ctx) error {
	// Cek apakah folder uploads ada, jika tidak, buat folder tersebut
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadsDir, os.ModePerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create uploads directory: " + err.Error())
		}

		return nil
	}
	// Remove the misplaced closing brace

	// Parse multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid form data: " + err.Error())
	}

	// Ambil data dari form
	userID := form.Value["id_user"][0]
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("User ID is required")
	}

	// Ambil data user berdasarkan user_id
	_, err = repository.GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch user: " + err.Error())
	}

	// Ambil data image dari form
	Image := form.File["menu_image"][0]

	// Cek apakah file desain image ada
	if Image == nil {
		return c.Status(fiber.StatusBadRequest).SendString("All design fields are required")
	}

	// Simpan file desain image ke folder uploads
	err = c.SaveFile(Image, filepath.Join(uploadsDir, Image.Filename))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save menu image: " + err.Error())
	}

	// Isi data menu dengan informasi dari user dan data foto
	var input model.MenuItem
	input.ID = primitive.NewObjectID()
	input.Image = Image.Filename // Menyimpan nama file gambar

	// Simpan portofolio ke database
	menuItemID, err := repository.CreateMenuItem(&input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create menu: " + err.Error())
	}

	// Return response berhasil
	return c.JSON(fiber.Map{
		"message":     "Menu created successfully",
		"menuItem_id": menuItemID,
	})
}

func GetAllMenuItem(c *fiber.Ctx) error {
	// Ambil data menuItem dari database
	menuItem, err := repository.GetAllMenuItem()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching menu")
	}

	// Return success response with all portfolios
	return c.JSON(fiber.Map{
		"menuItem": menuItem,
	})
}

func UpdateMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var menuItem model.MenuItem

	if err := c.BodyParser(&menuItem); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	menuItemID, err := repository.UpdateMenuItem(id, &menuItem)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating menu")
	}

	return c.JSON(fiber.Map{
		"message":     "Menu updated successfully",
		"menuItem_id": menuItemID,
	})
}

func DeleteMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	err := repository.DeleteMenuItem(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting menu")
	}
	return c.SendString("Menu deleted successfully")
}
