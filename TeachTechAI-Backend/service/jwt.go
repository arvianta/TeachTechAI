package service

import (
	"fmt"
	"log"
	"os"
	"teach-tech-ai/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService interface {
	GenerateToken(userID uuid.UUID, role string) (string, string, error)
	ValidateToken(token string) (*jwt.Token, error)
	RefreshToken(refreshToken string) (string, string, error)
	GetUserIDByToken(token string) (uuid.UUID, error)
	GetUserRoleByToken(token string) (string, error)
}

type jwtCustomClaim struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey 			string
	refreshSecretKey 	string
	issuer    			string
	userRepository 		repository.UserRepository
	roleRepository 		repository.RoleRepository
}

func NewJWTService(ur repository.UserRepository, rr repository.RoleRepository) JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		refreshSecretKey: getRefreshSecretKey(),
		issuer:    "teachtechai",
		userRepository: ur,
		roleRepository: rr,
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "teachtechai"
	}
	return secretKey
}

func getRefreshSecretKey() string {
	refreshSecretKey := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecretKey == "" {
		refreshSecretKey = "teachtechai_refresh"
	}
	return refreshSecretKey
}

func (j *jwtService) GenerateToken(userID uuid.UUID, role string) (string, string, error) {
	// Access token claims
	accessClaims := &jwtCustomClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessT, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	// Refresh token claims
	refreshClaims := &jwtCustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshT, err := refreshToken.SignedString([]byte(j.refreshSecretKey))
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return accessT, refreshT, nil
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (j *jwtService) RefreshToken(refreshToken string) (string, string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &jwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.refreshSecretKey), nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(*jwtCustomClaim); ok && token.Valid {
		roleID, err := j.userRepository.FindUserRoleIDByID(claims.UserID)
		if err != nil {
			return "", "", err
		}
		role, err := j.roleRepository.FindRoleNameByID(roleID)
		if err != nil {
			return "", "", err
		}
		return j.GenerateToken(claims.UserID, role)
	}

	return "", "", fmt.Errorf("invalid refresh token")
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

func (j *jwtService) GetUserRoleByToken(token string) (string, error) {
	tToken, err := j.ValidateToken(token)
	if err != nil {
		return "", err
	}
	claims, ok := tToken.Claims.(jwt.MapClaims)
	if !ok || !tToken.Valid {
		return "", fmt.Errorf("invalid token")
	}
	role := fmt.Sprintf("%v", claims["role"])
	return role, nil
}

//TODO: invalidate token function