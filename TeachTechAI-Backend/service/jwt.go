package service

import (
	"fmt"
	"os"
	"teach-tech-ai/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService interface {
	GenerateToken(userID uuid.UUID, role string) (string, string, time.Time, time.Time, error)
	ValidateToken(token string) (*jwt.Token, error)
	InvalidateToken(token string) error
	ValidateTokenWithDB(token string) (bool, error)
	RefreshToken(refreshToken string) (string, string, time.Time, time.Time, error)
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

func (j *jwtService) GenerateToken(userID uuid.UUID, role string) (string, string, time.Time, time.Time, error) {
	// Access token claims
	atx := time.Now().Add(time.Minute * 120)
	accessClaims := &jwtCustomClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(atx),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessT, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	// Refresh token claims
	rtx := time.Now().Add(time.Hour * 24 * 30)
	refreshClaims := &jwtCustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(rtx),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshT, err := refreshToken.SignedString([]byte(j.refreshSecretKey))
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return accessT, refreshT, atx, rtx, nil
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (j *jwtService) InvalidateToken(token string) error {
	userID, err := j.GetUserIDByToken(token)
	if err != nil {
		return err
	}

	err = j.userRepository.InvalidateUserToken(userID)
	if err != nil {
		return err
	}

	return nil
}

func (j *jwtService) ValidateTokenWithDB(token string) (bool, error) {
	userID, err := j.GetUserIDByToken(token)
	if err != nil {
		return false, err
	}

	dbToken, err := j.userRepository.GetUserSessionToken(userID)
	if err != nil {
		return false, err
	}

	if dbToken != token {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}

func (j *jwtService) RefreshToken(refreshToken string) (string, string, time.Time, time.Time, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &jwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.refreshSecretKey), nil
	})

	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	if claims, ok := token.Claims.(*jwtCustomClaim); ok && token.Valid {
		roleID, err := j.userRepository.FindUserRoleIDByID(claims.UserID)
		if err != nil {
			return "", "", time.Time{}, time.Time{}, err
		}
		role, err := j.roleRepository.FindRoleNameByID(roleID)
		if err != nil {
			return "", "", time.Time{}, time.Time{}, err
		}
		return j.GenerateToken(claims.UserID, role)
	}

	return "", "", time.Time{}, time.Time{}, fmt.Errorf("invalid refresh token")
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