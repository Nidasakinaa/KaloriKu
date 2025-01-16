package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID    	primitive.ObjectID 	`bson:"user_id,omitempty" json:"user_id,omitempty"`
	Username  	string             	`bson:"username,omitempty" json:"username,omitempty"`
	Email     	string             	`bson:"email,omitempty" json:"email,omitempty"`
	PhoneNumber string           	`bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Password  	string             	`bson:"password,omitempty" json:"password,omitempty"`
	Role     	int 				`bson:"role" json:"role"` // 0 = user, 1 = admin
}