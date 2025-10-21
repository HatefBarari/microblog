package repository

import (
	"context"
	"errors"
	"time"

	"github.com/HatefBarari/microblog-media/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMediaRepository struct {
	collection *mongo.Collection
}

func NewMongoMediaRepository(db *mongo.Database) *MongoMediaRepository {
	return &MongoMediaRepository{
		collection: db.Collection("media"),
	}
}

func (r *MongoMediaRepository) Create(ctx context.Context, media *domain.Media) error {
	// Generate ObjectID
	if media.ID == "" {
		media.ID = primitive.NewObjectID().Hex()
	}
	
	// Set timestamps
	now := time.Now()
	media.CreatedAt = now
	media.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, media)
	return err
}

func (r *MongoMediaRepository) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid media ID")
	}

	var media domain.Media
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&media)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &media, nil
}

func (r *MongoMediaRepository) GetByURL(ctx context.Context, url string) (*domain.Media, error) {
	var media domain.Media
	err := r.collection.FindOne(ctx, bson.M{"url": url}).Decode(&media)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &media, nil
}

func (r *MongoMediaRepository) ListByUploader(ctx context.Context, uploaderID string, limit, offset int) ([]*domain.Media, int, error) {
	// Build filter
	filter := bson.M{"uploader_id": uploaderID}
	
	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build options
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var media []*domain.Media
	if err = cursor.All(ctx, &media); err != nil {
		return nil, 0, err
	}

	return media, int(total), nil
}

func (r *MongoMediaRepository) Update(ctx context.Context, media *domain.Media) error {
	objectID, err := primitive.ObjectIDFromHex(media.ID)
	if err != nil {
		return errors.New("invalid media ID")
	}

	media.UpdatedAt = time.Now()
	
	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, media)
	return err
}

func (r *MongoMediaRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid media ID")
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoMediaRepository) DeleteByUploader(ctx context.Context, uploaderID, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid media ID")
	}

	filter := bson.M{
		"_id":         objectID,
		"uploader_id": uploaderID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("media not found or not owned by user")
	}

	return nil
}
