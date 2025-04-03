package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/services"
	"github.com/MdSadiqMd/Broadcast-API/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

type AuthConfig struct {
	JWTSecret     string
	TokenDuration time.Duration
	UserService   *services.UserService
}

type Auth struct {
	config AuthConfig
}

func NewAuth(config AuthConfig) *Auth {
	return &Auth{
		config: config,
	}
}

func (a *Auth) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicRoute(r) {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := extractToken(r)
			if tokenString == "" {
				utils.RespondError(w, http.StatusUnauthorized, "unauthorized: no token provided")
				return
			}

			claims, err := a.validateToken(tokenString)
			if err != nil {
				utils.RespondError(w, http.StatusUnauthorized, fmt.Sprintf("unauthorized: %v", err))
				return
			}

			ctx := utils.SetUserInContext(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublicRoute(r *http.Request) bool {
	path := r.URL.Path
	return path == "/api/login" || path == "/api/register" ||
		strings.HasPrefix(path, "/api/public/") ||
		strings.HasPrefix(path, "/api/health")
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:7]) == "BEARER " {
		return bearerToken[7:]
	}
	if len(bearerToken) > 6 && strings.ToUpper(bearerToken[0:6]) == "BEARER" {
		return bearerToken[6:]
	}

	return ""
}

func (a *Auth) validateToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		if a.config.UserService != nil {
			exists, err := a.config.UserService.UserExists(claims.UserID)
			if err != nil {
				return nil, err
			}
			if !exists {
				return nil, errors.New("user no longer exists")
			}
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (a *Auth) GenerateToken(user *models.User) (string, error) {
	claims := models.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.config.TokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "listmonk-clone",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.config.JWTSecret))
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := utils.GetUserFromContext(r.Context())
			if !ok {
				utils.RespondError(w, http.StatusUnauthorized, "unauthorized: no user in context")
				return
			}

			if claims.Role != role && claims.Role != "admin" {
				utils.RespondError(w, http.StatusForbidden, "forbidden: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
