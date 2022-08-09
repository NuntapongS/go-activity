package user

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getUserFunc func(context.Context) ([]User, error)

func (fn getUserFunc) GetUser(ctx context.Context) ([]User, error) {
	return fn(ctx)
}

func GetUserHandler(svc getUserFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := svc.GetUser(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, user)
	}
}

type getUserByIDFunc func(context.Context, string) (*User, error)

func (fn getUserByIDFunc) GetUserByID(ctx context.Context, UserID string) (*User, error) {
	return fn(ctx, UserID)
}

func GetUserByIDHandler(svc getUserByIDFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		UserID := c.Param("UserID")
		user, err := svc.GetUserByID(c.Request().Context(), UserID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, user)
	}
}
