package auth

import (
	"errors"
	"notex/api/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey     = []byte("your-secret-key") // 在生产环境中应该从配置文件或环境变量中读取
	tokenDuration = 24 * time.Hour            // 访问令牌有效期
)

// GenerateToken 生成 JWT 令牌
func GenerateToken(claims *dto.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"role":     claims.Role,
		"exp":      time.Now().Add(tokenDuration).Unix(),
		"iat":      time.Now().Unix(),
	})

	return token.SignedString(secretKey)
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 刷新令牌有效期7天
		"iat":     time.Now().Unix(),
	})

	return token.SignedString(secretKey)
}

// ParseToken 解析 JWT 令牌
func ParseToken(tokenString string) (*dto.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &dto.TokenClaims{
			UserID:   uint(claims["user_id"].(float64)),
			Username: claims["username"].(string),
			Role:     claims["role"].(string),
		}, nil
	}

	return nil, errors.New("invalid token")
}

// ParseRefreshToken 解析刷新令牌
func ParseRefreshToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return uint(claims["user_id"].(float64)), nil
	}

	return 0, errors.New("invalid refresh token")
}
