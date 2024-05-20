package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService interface {
	GenerateToken(userID uuid.UUID, role string) string
	ValidateToken(token string) (*jwt.Token, error)
	GetUserIDByToken(token string) (uuid.UUID, error)
}

type jwtCustomClaim struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "Template",
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "Template"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(userID uuid.UUID, role string) string {
	claims := &jwtCustomClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Println(err)
	}
	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (j *jwtService) GetUserIDByToken(token string) (uuid.UUID, error) {
	tToken, err := j.ValidateToken(token)
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := tToken.Claims.(jwt.MapClaims)
	if !ok || !tToken.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	id := fmt.Sprintf("%v", claims["user_id"])
	userID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
