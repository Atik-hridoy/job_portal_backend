package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte("your-secret-key"), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			// Check if user role is allowed
			isAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return RoleMiddleware("job_seeker", "hirer")(next)
}

func RequireJobSeeker(next http.Handler) http.Handler {
	return RoleMiddleware("job_seeker")(next)
}

func RequireEmployer(next http.Handler) http.Handler {
	return RoleMiddleware("hirer")(next)
}

func RequireEmployerOrAdmin(next http.Handler) http.Handler {
	return RoleMiddleware("hirer")(next)
}
