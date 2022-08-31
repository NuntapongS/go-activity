package login

import (
	"myapp/staff"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProfileHander() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)

		claims := user.Claims.(*JwtCustomClaims)

		UserID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "invalid user id",
			})
		}

		username := claims.Username

		role := claims.Role

		return c.JSON(http.StatusOK, staff.Staff{
			ID:       UserID,
			Username: username,
			Role:     role,
		})

	}
}
