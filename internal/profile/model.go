package profile

import "go.mongodb.org/mongo-driver/bson/primitive"

// Profile captures additional personal information for a verified user.
type Profile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name      string             `bson:"name" json:"name"`
	Gender    string             `bson:"gender" json:"gender"`
	Phone     string             `bson:"phone" json:"phone"`
	Bio       string             `bson:"bio" json:"bio"`
	Image     string             `bson:"image" json:"image"`
	Completed bool               `bson:"completed" json:"completed"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
