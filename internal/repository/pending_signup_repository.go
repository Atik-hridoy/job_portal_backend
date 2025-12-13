package repository

import (
	"context"

	"job_portal/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PendingSignupRepository manages pending signup records awaiting OTP verification.
type PendingSignupRepository struct {
	collection *mongo.Collection
}

// NewPendingSignupRepository binds the repository to the pending_signups collection.
func NewPendingSignupRepository(db *mongo.Database) *PendingSignupRepository {
	return &PendingSignupRepository{
		collection: db.Collection("pending_signups"),
	}
}

// Upsert stores or replaces a pending signup for the provided email.
func (r *PendingSignupRepository) Upsert(ctx context.Context, pending *models.PendingSignup) error {
	update := bson.M{
		"$set": bson.M{
			"email":      pending.Email,
			"password":   pending.Password,
			"role":       pending.Role,
			"phone":      pending.Phone,
			"otp":        pending.OTP,
			"otp_expiry": pending.OTPExpiry,
			"updated_at": pending.UpdatedAt,
		},
	}

	if pending.ID.IsZero() {
		update["$setOnInsert"] = bson.M{
			"created_at": pending.CreatedAt,
		}
	}

	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(ctx, bson.M{"email": pending.Email}, update, opts)
	if err != nil {
		return err
	}

	if result.UpsertedID != nil {
		if oid, ok := result.UpsertedID.(primitive.ObjectID); ok {
			pending.ID = oid
		}
	}

	return nil
}

// Delete removes a pending signup by email.
func (r *PendingSignupRepository) Delete(ctx context.Context, email string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"email": email})
	return err
}

// GetByEmail fetches a pending signup by email address.
func (r *PendingSignupRepository) GetByEmail(ctx context.Context, email string) (*models.PendingSignup, error) {
	var pending models.PendingSignup
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&pending)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &pending, nil
}
