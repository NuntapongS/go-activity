package login

import (
	"context"
	"myapp/staff"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type getStaffByUsername func(context.Context, string) (*staff.Staff, error)

func (fn getStaffByUsername) GetStaffByUsername(ctx context.Context, username string) (*staff.Staff, error) {
	return fn(ctx, username)
}

func LoginHandler(svs getStaffByUsername) echo.HandlerFunc {
	return func(c echo.Context) error {
		var login Login
		if err := c.Bind(&login); err != nil {
			return c.JSON(http.StatusBadRequest, bson.M{"error": err.Error()})
		}
		staff, err := svs.GetStaffByUsername(c.Request().Context(), login.Username)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, bson.M{"error": err.Error()})
		}
		if !checkPasswordHash(login.Password, staff.HashedPassword) {
			return c.JSON(http.StatusUnauthorized, bson.M{"error": "invalid password"})
		}
		claims := JwtCustomClaims{
			UserID:         staff.ID.Hex(),
			Username:       staff.Username,
			Role:           staff.Role,
			StandardClaims: jwt.StandardClaims{},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, bson.M{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, bson.M{"token": t})
	}
}

func checkPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
