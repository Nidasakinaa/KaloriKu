package repository

import (
	"KaloriKu/config"
	"KaloriKu/model"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByUsername(username string) (*model.User, error) {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("Users")
	var user model.User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Jika tidak ada pengguna dengan username ini, return nil
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *model.User) (string, error) {
	// Akses koleksi 'users' dalam database
	collection := config.GetMongoClient().Database("kaloriKu").Collection("Users")

	// Cek apakah username sudah ada
	existingUser, err := GetUserByUsername(user.Username)
	if err != nil {
		log.Println("Error checking username:", err)
		return "", err
	}
	if existingUser != nil {
		return "", fmt.Errorf("username '%s' already exists", user.Username)
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Error inserting user:", err)
		return "", err
	}
	
	// Konversi InsertedID ke ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("Failed to convert InsertedID to ObjectID")
		return "", fmt.Errorf("failed to convert inserted ID")
	}
	
	return objectID.Hex(), nil
	
}

func GetUserByID(id string) (*model.User, error) {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("Users")
	var user model.User

	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Jika tidak ada user dengan ID ini, return nil
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetAllUser() ([]model.User, error) {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("Users")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("Error fetching users:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []model.User
	for cursor.Next(context.Background()) {
		var user model.User
		err := cursor.Decode(&user)
		if err != nil {
			log.Println("Error decoding user:", err)
			return nil, err
		}

		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return users, nil
}

func UpdateUser(id string, user *model.User) (string, error) {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("Users")

	// Cek apakah user sudah ada
	existingUser, err := GetUserByID(id)
	if err != nil {
		log.Println("Error checking user:", err)
		return "", err
	}
	if existingUser == nil {
		return "", fmt.Errorf("user with ID '%s' does not exist", id)
	}

	// Hapus _id dari struct user untuk mencegah pengubahan _id
	updateData := bson.M{
		"username": user.Username,
		"email":    user.Email,
		"phone":    user.PhoneNumber,
		"role":     user.Role,
	}

	// Update user di MongoDB
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateData}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating user:", err)
		return "", err
	}

	return id, nil
}

func DeleteUser(id string) error {
	collection := config.GetMongoClient().Database("kaloriKu").Collection("Users")

	// Cek apakah user sudah ada
	existingUser, err := GetUserByID(id)
	if err != nil {
		log.Println("Error checking user:", err)
		return err
	}
	if existingUser == nil {
		return fmt.Errorf("user with ID '%s' does not exist", id)
	}

	// Hapus user dari MongoDB
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}

	return nil
}
