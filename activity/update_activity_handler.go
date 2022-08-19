package activity

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type updateActivityFunc func(context.Context, Activity) error

func (fn updateActivityFunc) UpdateActivity(ctx context.Context, activity Activity) error {
	return fn(ctx, activity)
}

func UpdateActivityHandler(svc updateActivityFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqActivity ActivityRequest

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

		activity := reqActivity.mapToActivity(startDate, endDate)
		activity.SetUpdatedAt(time.Now())

		err = svc.UpdateActivity(c.Request().Context(), activity)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "error",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "update successfully",
		})

	}
}
