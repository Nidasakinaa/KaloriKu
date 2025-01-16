package controller

import (
	"KaloriKu/model"
	"KaloriKu/repository"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var input model.User

	// Parse body request untuk mendapatkan data
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Validasi email dan nomor telepon
	if input.Username == "" || input.Email == "" || input.PhoneNumber == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("All fields are required")
	}

    // Jika role tidak disertakan, set role ke 1 (admin)
    if input.Role == 0 {
        input.Role = 2
    } else if input.Role != 1 && input.Role != 2 {
        // Validasi role hanya boleh 1 atau 2
        return c.Status(fiber.StatusBadRequest).SendString("Role must be 1 (admin) or 2 (customer)")
    }

	// Hash password sebelum disimpan
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}
	input.Password = string(hashPassword)

	// Simpan data user ke database
	userID, err := repository.CreateUser(&input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating user")
	}

	// Return response sukses dengan user ID
	return c.JSON(fiber.Map{
		"message": "User registered successfully",
		"user_id": userID,
	})
}

func Login(c *fiber.Ctx) error {
	var input model.User

	// Parse body request untuk mendapatkan data
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Validasi username dan password
	if input.Username == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Username and password are required")
	}

	// Ambil data user berdasarkan username
	user, err := repository.GetUserByUsername(input.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching user")
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	// Validasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid password")
	}

	// Set id user in cookies
	c.Cookie(&fiber.Cookie{
		Name:     "ID",                           // Nama cookie
		Expires:  time.Now().Add(24 * time.Hour), // Masa berlaku cookie
		HTTPOnly: false,                          // Akses hanya melalui HTTP (tidak dapat diakses oleh JS)
		Secure:   false,                          // Gunakan HTTPS jika true, gunakan false untuk localhost
		SameSite: "Lax",                          // Aturan SameSite untuk cookie
	})

	// Set role in cookies
	c.Cookie(&fiber.Cookie{
		Name:     "Role",                         // Nama cookie
		Expires:  time.Now().Add(24 * time.Hour), // Masa berlaku cookie
		HTTPOnly: false,                          // Akses hanya melalui HTTP (tidak dapat diakses oleh JS)
		Secure:   false,                          // Gunakan HTTPS jika true, gunakan false untuk localhost
		SameSite: "Lax",                          // Aturan SameSite untuk cookie
	})

	// Return response sukses
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"role":    user.Role,
	})

}
