package middleware

import (
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"

	"hello-world-go/internal/auth"
)

// Authenticate validates the Bearer token and stores claims in the request context.
// Returns 401 if the token is missing or invalid.
func Authenticate(verifier *oidc.IDTokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r)
			if token == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			idToken, err := verifier.Verify(r.Context(), token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			var claims auth.Claims
			if err := idToken.Claims(&claims); err != nil {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithClaims(r.Context(), claims)))
		})
	}
}

// RequireRole returns 403 if the authenticated user does not have the given role.
// Must be used after Authenticate.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := auth.FromContext(r.Context())
			if !ok || !claims.HasRole(role) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}
