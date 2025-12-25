package profile

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository persists profile data to MongoDB.
type Repository struct {
	collection *mongo.Collection
}

// NewRepository constructs a repository targeting the profiles collection.
func NewRepository(db *mongo.Database) *Repository {
	return &Repository{collection: db.Collection("profiles")}
}

// UpsertByUserID creates or updates a profile for the provided user.
func (r *Repository) UpsertByUserID(ctx context.Context, userID primitive.ObjectID, doc *Profile) (*Profile, error) {
	now := primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{
		"$set": bson.M{
			"name":       doc.Name,
			"gender":     doc.Gender,
			"phone":      doc.Phone,
			"bio":        doc.Bio,
			"image":      doc.Image,
			"completed":  true,
			"updated_at": now,
		},
		"$setOnInsert": bson.M{
			"user_id":    userID,
			"created_at": now,
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var out Profile
	err := r.collection.FindOneAndUpdate(ctx, bson.M{"user_id": userID}, update, opts).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			if err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&out); err != nil {
				return nil, err
			}
			return &out, nil
		}
		return nil, err
	}

	return &out, nil
}

// GetByUserID retrieves the profile for a user if it exists.
func (r *Repository) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*Profile, error) {
	var result Profile
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
