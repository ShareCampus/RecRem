package utils

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// 常量
var (
	ErrTokenExpired     = errors.New("令牌已过期")
	ErrTokenNotValidYet = errors.New("令牌未激活")
	ErrTokenMalformed   = errors.New("令牌格式有误")
	ErrTokenInvalid     = errors.New("无效的令牌")
	SignKey             = "aries-open-source-blog" // 签名
)

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// CustomClaims 载荷
type CustomClaims struct {
	Username string `json:"username"`
	UserImg  string `json:"user_img"`
	jwt.StandardClaims
}

// NewJWT 创建一个 JWT 实例
func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(GetSignKey()),
	}
}

// GetSignKey 获取 SignKey
func GetSignKey() string {
	return SignKey
}

// CreateToken 创建 Token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析 Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil
		})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
