package activity

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type deleteActivityFunc func(context.Context, string) error

func (fn deleteActivityFunc) DeleteActivity(ctx context.Context, activityID string) error {
	return fn(ctx, activityID)
}

func DeleteActivityHandler(svc deleteActivityFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		activityID := c.Param("activityid")

		err := svc.DeleteActivity(c.Request().Context(), activityID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "error",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Delete successfully",
		})
	}
}
