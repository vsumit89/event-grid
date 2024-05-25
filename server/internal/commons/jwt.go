package commons

import (
	"errors"
	"server/internal/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenNotFound = errors.New("bearer token not found")
	ErrTokenExpired  = errors.New("access token expired")

	ErrTokenInvalid = errors.New("invalid token")
)

type JwtSvc struct {
	SecretKey string
	TTL       time.Duration
}

// CustomClaims includes standard JWT claims and custom user-defined claims
type CustomClaims struct {
	Email  string `json:"username"`
	UserID uint   `json:"userid"`
	jwt.StandardClaims
}

func NewJWTService(cfg *config.JWTConfig) *JwtSvc {
	return &JwtSvc{
		SecretKey: cfg.Secret,
		TTL:       time.Duration(cfg.TTL) * getUnit(cfg.Unit),
	}
}

func getUnit(unit string) time.Duration {
	switch unit {
	case "minutes":
		return time.Minute
	case "hours":
		return time.Hour
	case "days":
		return time.Hour * 24
	default:
		return time.Minute
	}
}

// GenerateToken generates a JWT token with custom claims
func (j *JwtSvc) GenerateToken(email string, user uint) (string, error) {
	claims := CustomClaims{
		Email:  email,
		UserID: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.TTL).Unix(),
			Issuer:    "server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtSvc) GetTTL() time.Duration {
	return j.TTL
}

// ValidateToken validates a JWT token and returns the custom claims if the token is valid
func (j *JwtSvc) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
