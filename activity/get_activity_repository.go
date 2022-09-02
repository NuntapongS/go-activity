package activity

import (
	"context"
	"myapp/pagination"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetActivity(db *mongo.Database) func(context.Context, pagination.Pagination) ([]Activity, int64, error) {
	return func(ctx context.Context, pag pagination.Pagination) ([]Activity, int64, error) {
		collection := getActivityCollection(db)

		var filterSearchSXID,
			filterFullName,
			filterCheckinType,
			filterPhone,
			filterTargetDate,
			filterSearch primitive.E

		if pag.Sort == "" {
			pag.Sort = "created_at"
		}

		if pag.TargetDate != "" {
			start, err := time.Parse("2006-01-02", pag.TargetDate)
			if err == nil {
				end := start.AddDate(0, 0, 1).Add(-time.Second)
				filterTargetDate = bson.E{Key: "created_at", Value: bson.M{
					"$gte": primitive.NewDateTimeFromTime(start),
					"$lt":  primitive.NewDateTimeFromTime(end),
				}}
			}

		}

		if pag.Type != "" {
			filterCheckinType = bson.E{
				Key: "checkin_type", Value: primitive.Regex{
					Pattern: pag.Type,
					Options: "im",
				},
			}
		}

		if pag.Search != "" {
			filterSearchSXID = bson.E{
				Key: "sx_id", Value: primitive.Regex{
					Pattern: pag.Search,
					Options: "im",
				},
			}
			filterFullName = bson.E{
				Key: "full_name", Value: primitive.Regex{
					Pattern: pag.Search,
					Options: "im",
				},
			}
			filterPhone = bson.E{
				Key: "phone", Value: primitive.Regex{
					Pattern: pag.Search,
					Options: "im",
				}}
			filterSearch = bson.E{
				Key: "$or", Value: []bson.D{
					{filterSearchSXID},
					{filterFullName},
					{filterPhone},
				},
			}
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{filterSearch, filterTargetDate, filterCheckinType}}}

		cur, err := collection.Aggregate(context.Background(), mongo.Pipeline{matchStage})
		if err != nil {
			return nil, 0, err
		}
		var total int64
		for cur.Next(context.Background()) {
			total++
		}
		if err := cur.Err(); err != nil {
			return nil, 0, err
		}
		cur.Close(context.Background())

		limitStage := bson.D{{
			Key: "$limit", Value: pag.GetPageSize(),
		}}
		skipStage := bson.D{{
			Key: "$skip", Value: (pag.GetPage() - 1) * pag.GetPageSize(),
		}}

		var sortStage primitive.D
		if pag.Sort != "created_at" {
			sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: pag.Sort, Value: pag.GetDirection()}}}}
		} else {
			sortStage = bson.D{{Key: "$sort", Value: bson.D{
				{Key: "created_at", Value: pag.GetDirection()},
			}}}
		}

		results := []Activity{}
		var pipeline mongo.Pipeline
		if pag.IsDownload {
			pipeline = mongo.Pipeline{matchStage, sortStage}
		} else {
			pipeline = mongo.Pipeline{matchStage, skipStage, limitStage, sortStage}
		}

		var cursor *mongo.Cursor
		cursor, err = collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			return nil, 0, err
		}

		if err = cursor.All(context.Background(), &results); err != nil {
			panic(err)
		}

		cursor.Close(context.Background())

		return results, total, err
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
