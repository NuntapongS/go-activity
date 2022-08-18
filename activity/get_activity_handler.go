package activity

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (act Activity) DateToString(date *time.Time) string {
	if date == nil {
		return ""
	}
	rs := date.Format("2006-01-02")
	return rs
}

type getActivityFunc func(context.Context) ([]Activity, error)

func (fn getActivityFunc) GetActivity(ctx context.Context) ([]Activity, error) {
	return fn(ctx)
}

func GetActivityHandler(svc getActivityFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		activities, err := svc.GetActivity(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "not found",
			})
		}

		var response []ActivityResponse
		for i := range activities {
			startDate := activities[i].StartDate
			endDate := activities[i].EndDate
			rs := activities[i].mapToActivityResponse(activities[i].DateToString(startDate), activities[i].DateToString(endDate))
			response = append(response, rs)
		}
		return c.JSON(http.StatusOK, response)

	}
}
