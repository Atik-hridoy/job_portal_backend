package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PendingSignup struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Role      string             `bson:"role" json:"role" validate:"oneof=job_seeker hirer"`
	Phone     string             `bson:"phone" json:"phone"`
	OTP       string             `bson:"otp" json:"-"`
	OTPExpiry primitive.DateTime `bson:"otp_expiry" json:"otp_expiry"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
