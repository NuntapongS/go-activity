package activity

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateActivity(db *mongo.Database) func(context.Context, Activity) error {
	return func(ctx context.Context, activity Activity) error {
		collection := getActivityCollection(db)

		filter := bson.M{
			"$and": []bson.M{
				{"activity_id": bson.M{"$eq": activity.ActivityID}},
			},
		}
		rs, err := collection.ReplaceOne(ctx, filter, activity)
		if rs.ModifiedCount == 0 {
			return errors.New("activity can not update")
		}
		return err
	}
}
