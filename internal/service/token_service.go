package service

import (
	"fmt"
	"justn0w-bot/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type CustomClaims struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(viper.GetString("jwt.secret"))

func GenerateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	// 声明JWT的载荷claims
	claims := &CustomClaims{
		UserId:   user.Id,
		UserName: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), //过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "justn0w-bot",
			Subject:   fmt.Sprintf("%d", user.Id),
		},
	}
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名token
	return token.SignedString(jwtSecret)
}

// ParseToken 解析token
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 1 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
