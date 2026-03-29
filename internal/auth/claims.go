package auth

import "context"

// Claims holds the validated token claims we care about.
type Claims struct {
	Subject  string   `json:"sub"`
	Name     string   `json:"name"`
	Email    string   `json:"preferred_username"`
	Roles    []string `json:"roles"`
}

// HasRole returns true if the claims include the given role.
func (c Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

type contextKey struct{}

// WithClaims stores claims in the request context.
func WithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, contextKey{}, claims)
}

// FromContext retrieves claims from the request context.
func FromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(contextKey{}).(Claims)
	return claims, ok
}
