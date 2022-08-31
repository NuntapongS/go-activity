package login

import (
	"myapp/staff"

	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	UserID   string     `json:"id"`
	Username string     `json:"username"`
	Role     staff.Role `json:"role"`
	ClientId string     `json:"client_id"`
	Sub      string     `json:"sub"`
	AuthTime int64      `json:"auth_time"`
	Idp      string     `json:"ldp"`
	Iat      int64      `json:"iat"`
	Scope    []string   `json:"scope"`
	Arm      []string   `json:"amr"`
	jwt.StandardClaims
}
