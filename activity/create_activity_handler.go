package activity

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (reqAct ActivityRequest) GetDate(dateFormat string) (*time.Time, error) {
	date, err := time.Parse("2006-01-02", dateFormat)
	if err != nil && dateFormat != "" {
		return nil, err
	}
	return &date, err
}

func (a *Activity) SetCreatedAt(date time.Time) {
	a.CreatedAt = date
}

func (a *Activity) SetUpdatedAt(date time.Time) {
	a.UpdatedAt = date
}

type CreateActivityFunc func(context.Context, Activity) error

func (fn CreateActivityFunc) CreateActivity(ctx context.Context, activity Activity) error {
	return fn(ctx, activity)
}

func CreateActivityhandler(svc CreateActivityFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqActivity := ActivityRequest{}

		if err := c.Bind(&reqActivity); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "error",
			})
		}

		startDate, err := reqActivity.GetDate(reqActivity.StartDate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "error",
			})
		}

		endDate, err := reqActivity.GetDate(reqActivity.EndDate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "error",
			})
		}

		actvity := reqActivity.mapToActivity(startDate, endDate)
		createdAt, updatedAt := time.Now(), time.Now()
		actvity.SetCreatedAt(createdAt)
		actvity.SetUpdatedAt(updatedAt)

		err = svc.CreateActivity(c.Request().Context(), actvity)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "error",
			})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"message": "create successfully",
		})
	}
}
