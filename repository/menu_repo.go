package repository

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"KaloriKu/config"
	"KaloriKu/model"
)

// var menuCollection *mongo.Collection

func GetMenuByName(name string) (*model.MenuItem, error) {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("MenuItem")

	var menuItem model.MenuItem
	err := collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&menuItem)
	if err != nil {
		log.Printf("Could not find menu item: %v", err)
		return nil, err
	}
	return &menuItem, nil
}

func CreateMenuItem(menuItem *model.MenuItem) (string, error) {
	// Akses koleKsi 'MenuItem' dalam database
	collection := config.GetMongoClient().Database("kaloriKu").Collection("MenuItem")

	// Cek apakah username sudah ada
	existingMenuItem, err := GetMenuByName(menuItem.Name)
	if err != nil {
		log.Println("Error checking menu item name: ", err)
		return "", err
	}
	if existingMenuItem != nil {
		return "", fmt.Errorf("username '%s' already exists", menuItem.Name)
	}

	// Insert menu baru ke dalam MongoDB
	result, err := collection.InsertOne(context.Background(), menuItem)
	if err != nil {
		log.Println("Error inserting menuItem:", err)
		return "", err
	}

	// Mengembalikan ID menuItem yang baru disimpan
	id := result.InsertedID
	return id.(string), nil
}

func GetMenuItemByID(id string) (*model.MenuItem, error) {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("MenuItem")
	var menuItem model.MenuItem

	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&menuItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Jika tidak ada user dengan ID ini, return nil
			return nil, nil
		}
		return nil, err
	}
	return &menuItem, nil
}

func GetAllMenuItem() ([]model.MenuItem, error) {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("MenuItem")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("Error fetching menuItem:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var menuItems []model.MenuItem
	for cursor.Next(context.Background()) {
		var menuItem model.MenuItem
		err := cursor.Decode(&menuItem)
		if err != nil {
			log.Println("Error decoding menuItem:", err)
			return nil, err
		}

		menuItems = append(menuItems, menuItem)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return menuItems, nil
}

func UpdateMenuItem(id string, menuItem *model.MenuItem) (string, error) {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("MenuItem")

	// Cek apakah menuItem sudah ada
	existingMenuItem, err := GetMenuItemByID(id)
	if err != nil {
		log.Println("Error checking menuItem:", err)
		return "", err
	}
	if existingMenuItem == nil {
		return "", fmt.Errorf("menuItem with ID '%s' does not exist", id)
	}

	// Hapus _id dari struct MenuItem untuk mencegah pengubahan _id
	updateData := bson.M{
		"name":        menuItem.Name,
		"description": menuItem.Description,
		"category":    menuItem.Category,
		"image":       menuItem.Image,
		"stock":       menuItem.Stock,
	}

	// Update menuItem di MongoDB
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateData}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating menuItem:", err)
		return "", err
	}

	return id, nil
}

func DeleteMenuItem(id string) error {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("MenuItem")

	// Cek apakah menuItem sudah ada
	existingMenuItem, err := GetMenuItemByID(id)
	if err != nil {
		log.Println("Error checking menuItem:", err)
		return err
	}
	if existingMenuItem == nil {
		return fmt.Errorf("menuItem with ID '%s' does not exist", id)
	}

	// Hapus menuItem dari MongoDB
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Println("Error deleting menuItem:", err)
		return err
	}

	return nil
}
