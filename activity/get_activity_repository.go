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
