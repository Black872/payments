package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"payments/config"
	"payments/models"
	"payments/repository"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/zeebo/errs"
)

var authErr = errs.Class("authorization service")

var (
	// The method to sign a token https://pkg.go.dev/github.com/golang-jwt/jwt/v5#SigningMethod
	tokenSigningMethod = jwt.SigningMethodHS256

	// The name of the alg identifier for the tokenSigningMethod.
	tokenSigningMethodAlgName = "HS256"
)

// The custom claims with the user ID.
type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

// AuthorizationService contains repository interface Authorization.
type AuthorizationService struct {
	repo repository.Authorization
}

// NewAuthorizationService creates a new authorization service which contains the underlying
// repository Authorization interface.
func NewAuthorizationService(r *repository.Repository) *AuthorizationService {
	return &AuthorizationService{repo: r}
}

// CreateUser receives the user details from the (*AuthorizationHandlers) SignUp() method,
// replaces the password with a password hash and sends the user details to the
// (*AuthorizationDB) CreateUser() method.
func (s *AuthorizationService) CreateUser(ctx context.Context, user models.User) (id int, err error) {

	user.PasswordHash = generatePasswordHash(user.PasswordHash)
	if id, err = s.repo.CreateUser(ctx, user); err != nil {
		return 0, authErr.Wrap(err)
	}
	return id, nil
}

// GenerateToken receives the user details (email, password) from the (*AuthorizationHandlers) Login() method,
// and returns the signed token with claims containing user ID.
func (s *AuthorizationService) GenerateToken(ctx context.Context, user models.User) (signedToken string, err error) {
	id, err := s.repo.GetUserID(ctx, user.Email, generatePasswordHash(user.PasswordHash))
	if err != nil {
		return "", authErr.Wrap(err)
	}

	claims := tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenTTL())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: id,
	}
	token := jwt.NewWithClaims(tokenSigningMethod, claims)

	if signedToken, err = token.SignedString([]byte(config.TokenSignature())); err != nil {
		return "", authErr.Wrap(err)
	}

	return signedToken, nil
}

// ParseToken receives a signed token string, parses and validates it and returns the user ID.
func (s *AuthorizationService) ParseToken(accessToken string) (id int, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != tokenSigningMethodAlgName {
			return 0, authErr.New("ParseToken(): invalid signing method of token")
		}
		return []byte(config.TokenSignature()), nil
	})
	if err != nil {
		return 0, authErr.Wrap(err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, authErr.New("ParseToken(): token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

// Makes the password hash from the raw password.
func generatePasswordHash(password string) (passwordHash string) {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(config.HashSalt())))
}
