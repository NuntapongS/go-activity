package activity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ActivityID string             `bson:"activity_id" json:"activity_id"`
	Name       string             `bson:"name" json:"name"`
	Zone       string             `bson:"zone" json:"zone"`
	StartDate  *time.Time         `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate    *time.Time         `bson:"end_date,omitempty" json:"end_date,omitempty"`
	StartTime  string             `bosn:"start_time,omitempty" json:"start_time,omitempty"`
	EndTime    string             `bson:"end_time,omitempty" json:"end_time,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type Activities struct {
	TotalRecords int64      `json:"total_records"`
	Rows         int64      `json:"rows"`
	Page         int64      `json:"page"`
	Data         []Activity `json:"data"`
}

type ActivityRequest struct {
	ActivityID string     `json:"activity_id,omitempty"`
	Name       string     `json:"name"`
	Zone       string     `json:"zone"`
	StartDate  string     `json:"start_date,omitempty"`
	EndDate    string     `json:"end_date,omitempty"`
	StartTime  string     `json:"start_time,omitempty"`
	EndTime    string     `json:"end_time,omitempty"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

func (reqActivity *ActivityRequest) mapToActivity(startDate, endDate *time.Time) Activity {
	return Activity{
		ActivityID: reqActivity.ActivityID,
		Name:       reqActivity.Name,
		Zone:       reqActivity.Zone,
		StartDate:  startDate,
		EndDate:    endDate,
		StartTime:  reqActivity.StartTime,
		EndTime:    reqActivity.EndTime,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

type ActivityResponse struct {
	ID         primitive.ObjectID `json:"id,omitempty"`
	ActivityID string             `bson:"activity_id" json:"activity_id"`
	Name       string             `bson:"name" json:"name"`
	Zone       string             `bson:"zone" json:"zone"`
	StartDate  string             `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate    string             `bson:"end_date,omitempty" json:"end_date,omitempty"`
	StartTime  string             `bson:"start_time,omitempty" json:"start_time"`
	EndTime    string             `bson:"end_time,omitempty" json:"end_time,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

func (Activity *Activity) mapToActivityResponse(startDate, endDate string) ActivityResponse {
	return ActivityResponse{
		ID:         Activity.ID,
		ActivityID: Activity.ActivityID,
		Name:       Activity.Name,
		Zone:       Activity.Zone,
		StartDate:  startDate,
		EndDate:    endDate,
		StartTime:  Activity.StartTime,
		EndTime:    Activity.EndTime,
		CreatedAt:  Activity.CreatedAt,
		UpdatedAt:  Activity.UpdatedAt,
	}
}
