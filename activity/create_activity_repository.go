package activity

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateActivity(db *mongo.Database) func(context.Context, Activity) error {
	return func(ctx context.Context, activity Activity) error {
		collection := getActivityCollection(db)
		_, err := collection.InsertOne(ctx, activity)
		return err
	}
}
