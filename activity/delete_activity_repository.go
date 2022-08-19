package activity

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteActivity(db *mongo.Database) func(context.Context, string) error {
	return func(ctx context.Context, activityID string) error {
		collection := getActivityCollection(db)

		filter := bson.M{
			"activity_id": bson.M{"$eq": activityID},
		}

		_, err := collection.DeleteOne(ctx, filter)

		return err

	}
}
