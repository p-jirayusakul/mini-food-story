package middleware

import (
	"fmt"
	"food-story/pkg/exceptions"
	"io"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const PrefixAuth = "Bearer "

type AuthInterface interface {
	JWTMiddleware() fiber.Handler
	RequireRole(role []string) fiber.Handler
}

type AuthInstance struct {
	auth keyfunc.Keyfunc
}

func NewAuthInstance(keycloakCertURL string) *AuthInstance {

	// validate keycloakCertURL
	resp, err := http.Get(keycloakCertURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Failed to GET from keycloakCertURL: %s, error: %v", keycloakCertURL, err))
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	auth, err := keyfunc.NewDefault([]string{keycloakCertURL})
	if err != nil {
		panic(err)
	}

	return &AuthInstance{
		auth: auth,
	}
}

func (i *AuthInstance) JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, PrefixAuth) {
			return ResponseError(c, exceptions.Error(exceptions.CodeUnauthorized, "Authorization header is required"))
		}

		tokenString := strings.TrimPrefix(authHeader, PrefixAuth)
		token, err := jwt.Parse(tokenString, i.auth.Keyfunc)
		if err != nil {
			return ResponseError(c, exceptions.Error(exceptions.CodeUnauthorized, "failed to parse the JWT"))
		}

		if !token.Valid {
			return ResponseError(c, exceptions.Error(exceptions.CodeUnauthorized, "the token is not valid"))
		}

		c.Locals("jwt", token)
		return c.Next()
	}
}

func (i *AuthInstance) RequireRole(findRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := c.Locals("jwt").(*jwt.Token)
		claims := t.Claims.(jwt.MapClaims)

		realmAccess, ok := claims["realm_access"].(map[string]interface{})
		if !ok {
			return ResponseError(c, exceptions.Error(exceptions.CodeForbidden, "no roles"))
		}

		rolesData, ok := realmAccess["roles"].([]interface{})
		if !ok {
			return ResponseError(c, exceptions.Error(exceptions.CodeForbidden, "invalid roles"))
		}

		roles := make([]string, len(rolesData))
		for idx, role := range rolesData {
			roles[idx] = role.(string)
		}

		roleMap := make(map[string]bool)
		for _, role := range roles {
			roleMap[role] = true
		}

		for _, fr := range findRoles {
			if roleMap[fr] {
				return c.Next()
			}
		}

		return ResponseError(c, exceptions.Error(exceptions.CodeForbidden, "permission denied"))
	}
}
