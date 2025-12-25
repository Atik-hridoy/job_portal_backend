package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"job_portal/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository handles persistence for user accounts.
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a repository bound to the users collection.
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// CreateUser inserts a new user document with timestamps.
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	now := primitive.NewDateTimeFromTime(time.Now())
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return nil
}

// GetByEmail finds a user by email address.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var raw bson.M
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&raw)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return decodeUser(raw)
}

// GetByID finds a user by their ObjectID.
func (r *UserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var raw bson.M
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&raw)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return decodeUser(raw)
}

func decodeUser(doc bson.M) (*models.User, error) {
	if doc == nil {
		return nil, nil
	}

	user := &models.User{}

	if v, ok := doc["_id"].(primitive.ObjectID); ok {
		user.ID = v
	}
	if v, ok := doc["email"].(string); ok {
		user.Email = v
	}
	if v, ok := doc["password"].(string); ok {
		user.Password = v
	}
	if v, ok := doc["role"].(string); ok {
		user.Role = v
	}
	if phoneVal, ok := doc["phone"]; ok {
		phone, err := normalizePhone(phoneVal)
		if err != nil {
			return nil, err
		}
		user.Phone = phone
	}
	if v, ok := doc["is_verified"].(bool); ok {
		user.IsVerified = v
	}
	if v, ok := doc["created_at"].(primitive.DateTime); ok {
		user.CreatedAt = v
	}
	if v, ok := doc["updated_at"].(primitive.DateTime); ok {
		user.UpdatedAt = v
	}

	return user, nil
}

func normalizePhone(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return strconv.FormatInt(int64(v), 10), nil
	case primitive.Decimal128:
		return v.String(), nil
	case nil:
		return "", nil
	default:
		return fmt.Sprint(v), nil
	}
}
