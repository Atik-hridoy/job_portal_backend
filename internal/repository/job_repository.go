package repository

import (
	"context"
	"log"
	"time"

	"job_portal/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(db *mongo.Database) *JobRepository {
	return &JobRepository{
		collection: db.Collection("jobs"),
	}
}

func (r *JobRepository) CreateJob(job *models.Job) (*models.Job, error) {
	job.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	job.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := r.collection.InsertOne(context.Background(), job)
	if err != nil {
		return nil, err
	}

	job.ID = result.InsertedID.(primitive.ObjectID)
	return job, nil
}

func (r *JobRepository) GetJobByID(id string) (*models.Job, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var job models.Job
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *JobRepository) GetAllJobs() ([]*models.Job, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error finding jobs: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var jobs []*models.Job
	for cursor.Next(context.Background()) {
		var job models.Job
		err := cursor.Decode(&job)
		if err != nil {
			log.Printf("Error decoding job: %v", err)
			continue
		}
		jobs = append(jobs, &job)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *JobRepository) UpdateJob(id string, job *models.Job) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	job.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": job},
	)

	return err
}

func (r *JobRepository) DeleteJob(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}
