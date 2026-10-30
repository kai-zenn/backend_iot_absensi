package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type JwtInterface interface {
  GenerateToken(userId uuid.UUID, role string) (string, error)
  VerifyToken(tokenStr string) (*Claims, error)
}

type JsonWebToken struct {
  SecretKey []byte
  ExpiredTime time.Duration
}

func NewJWTToken(secretKey string, expiredTime time.Duration) JwtInterface {
  return &JsonWebToken{
    SecretKey: []byte(secretKey),
    ExpiredTime: expiredTime,
  }
}

func (j *JsonWebToken) GenerateToken(userId uuid.UUID, role string) (string, error) {
  claims := Claims{
    UserId: userId,
    Role: role,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiredTime)),
      IssuedAt: jwt.NewNumericDate(time.Now()),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(j.SecretKey)
}

func (j *JsonWebToken) VerifyToken(tokenStr string) (*Claims, error) {
  token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func (t *jwt.Token) (interface{}, error) {
    return j.SecretKey, nil
  })

  if token == nil {
    return nil, errors.New("Token is null")
  }

  if err != nil {
    return nil, err
  }

  claims, ok := token.Claims.(*Claims)
  if !ok || !token.Valid {
   return nil, errors.New("Invalid Token")
  } 

 return claims, nil 
}
