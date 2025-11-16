package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Job represents a job listing in the system
type Job struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Company     string             `bson:"company" json:"company"`
	Location    string             `bson:"location" json:"location"`
	Salary      string             `bson:"salary" json:"salary"`
	Description string             `bson:"description" json:"description"`
	Skills      []string           `bson:"skills" json:"skills"`
	Type        string             `bson:"type" json:"type"` // full-time, part-time, contract, etc.
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
