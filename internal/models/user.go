package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email      string             `bson:"email" json:"email" validate:"required,email"`
	Password   string             `bson:"password" json:"-" validate:"required,min=8"`
	Role       string             `bson:"role" json:"role" validate:"required,oneof=job_seeker hirer"`
	Phone      string             `bson:"phone" json:"phone" validate:"required"`
	IsVerified bool               `bson:"is_verified" json:"is_verified"`
	CreatedAt  primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt  primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
