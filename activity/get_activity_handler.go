package activity

import (
	"context"
	"myapp/pagination"
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

type getActivityFunc func(context.Context, pagination.Pagination) ([]Activity, int64, error)

func (fn getActivityFunc) GetActivity(ctx context.Context, pag pagination.Pagination) ([]Activity, int64, error) {
	return fn(ctx, pag)
}

func GetActivityHandler(svc getActivityFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		pag := pagination.Pagination{
			Search:     c.QueryParam("search"),
			TargetDate: c.QueryParam("target_date"),
			Type:       c.QueryParam("type"),
			Sort:       c.QueryParam("sort"),
			Direction:  c.QueryParam("direction"),
			Page:       c.QueryParam("page"),
			PageSize:   c.QueryParam("page_size"),
			IsDownload: c.QueryParam("is_download") == "true",
		}

		activities, total, err := svc.GetActivity(c.Request().Context(), pag)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "not found",
			})
		}

		return c.JSON(http.StatusOK, Activities{
			total,
			int64(len(activities)),
			pag.GetPage(),
			activities,
		})
	}
}

type getActivityByActivityID func(context.Context, string) (Activity, error)

func (fn getActivityByActivityID) GetActivityByActivityID(ctx context.Context, ActivityID string) (Activity, error) {
	return fn(ctx, ActivityID)
}

func GetActivityByActivityIDHandler(svc getActivityByActivityID) echo.HandlerFunc {
	return func(c echo.Context) error {
		activityID := c.Param("activityid")

		activity, err := svc.GetActivityByActivityID(c.Request().Context(), activityID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "error",
			})
		}

		startDate := activity.StartDate
		endDate := activity.EndDate
		rs := activity.mapToActivityResponse(activity.DateToString(startDate), activity.DateToString(endDate))

		return c.JSON(http.StatusOK, rs)
	}
}
