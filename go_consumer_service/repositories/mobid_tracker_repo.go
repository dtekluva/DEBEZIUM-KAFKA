package repositories

import (
	"context"
	"go_consumer_service/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MobidTrackerRepo struct {
	database *mongo.Database
}

func NewMobidTrackerRepo(database *mongo.Database) *MobidTrackerRepo {
	return &MobidTrackerRepo{database: database}
}

func (r *MobidTrackerRepo) GetAll(ctx context.Context, limit, offset int) (*[]types.MobidTracker, int, error) {
	collection := r.database.Collection("mobid_tracker")
	// Count total docs
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	// Query with pagination
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "date_created", Value: -1}})

	// Find all documents in the collection
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var trackers []types.MobidTracker
	if err := cursor.All(ctx, &trackers); err != nil {
		return nil, 0, err
	}

	return &trackers, int(total), nil
}
