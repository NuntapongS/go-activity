package user

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type createUserFunc func(context.Context, User) error

func (fn createUserFunc) CreateUser(ctx context.Context, reqUser User) error {
	return fn(ctx, reqUser)
}

func CreateUserHandler(svc createUserFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqUser User

		if err := c.Bind(&reqUser); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

			err := svc.CreateUser(c.Request().Context(), reqUser)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
			}

			return c.JSON(http.StatusOK, map[string]string{
				"meesage": "OK",
			})
		}
	}

