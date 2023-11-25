package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type SubscriberAuth struct {
	Token      string        // Token
	Expiration time.Duration // 过期时间
	SecretKey  string        // 密钥
}

// GenToken 生成JWT token
func (sa *SubscriberAuth) GenToken(username string) (token string, err error) {
	// 创建一个新的令牌对象
	tk := jwt.New(jwt.SigningMethodHS256)

	// 设置令牌的claims
	claims := tk.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(sa.Expiration).Unix()

	// 生成令牌
	token, err = tk.SignedString([]byte(sa.SecretKey))
	if err != nil {
		return "", err
	} else {
		sa.Token = token
	}
	return token, nil
}

// ValidateToken 验证和检查是否过期
func (sa *SubscriberAuth) ValidateToken(tokenString string) (bool, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(sa.SecretKey), nil
	})

	if err != nil {
		return false, fmt.Errorf("token validation error: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().UTC().After(expirationTime) {
			return false, errors.New("token has expired")
		}
		fmt.Println("过期时间: ", expirationTime)
		return true, nil
	}

	return false, errors.New("invalid access token")
}

// GenUUID 生成UUID
func GenUUID(username string) string {
	ns := uuid.NameSpaceDNS
	userUUID := uuid.NewSHA1(ns, []byte(username))
	return userUUID.String()
}
