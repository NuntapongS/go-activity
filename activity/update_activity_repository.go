package activity

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateActivity(db *mongo.Database) func(context.Context, Activity) error {
	return func(ctx context.Context, activity Activity) error {
		collection := getActivityCollection(db)

		var result Activity
		filter := bson.M{"employee_id": bson.M{"$eq": (activity.ActivityID)}}
		findResult := collection.FindOne(ctx, filter)
		findResult.Decode(&result)

		_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": activity})
		return err

	}
}
