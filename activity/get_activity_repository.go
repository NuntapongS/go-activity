package activity

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetActivity(db *mongo.Database) func(context.Context) ([]Activity, error) {
	return func(ctx context.Context) ([]Activity, error) {
		collection := getActivityCollection(db)

		var results []Activity
		cur, err := collection.Find(ctx, bson.D{})

		for cur.Next(ctx) {
			var elem Activity
			err := cur.Decode(&elem)
			if err != nil {
				return nil, err
			}

			results = append(results, elem)

		}

		if cur.Err() != nil {
			return nil, cur.Err()
		}

		cur.Close(ctx)
		return results, err
	}
}

func GetActivityByActivityID(db *mongo.Database) func(context.Context, string) (Activity, error) {
	return func(ctx context.Context, activityID string) (Activity, error) {
		collection := getActivityCollection(db)
		filter := bson.M{"activity_id": bson.M{"$eq": activityID}}
		rs := collection.FindOne(ctx, filter)

		result := Activity{}
		err := rs.Decode(&result)
		return result, err
	}
}
