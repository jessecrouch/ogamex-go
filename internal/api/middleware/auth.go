package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

const UserIDKey = "user_id"
const UserKey = "user"

func AuthMiddleware(authService interface {
	ValidateToken(ctx interface{}, token string) (interface{}, error)
}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "authorization required"})
		}

		parts := []string{}
		for _, p := range strings.Split(authHeader, " ") {
			if p != "" {
				parts = append(parts, p)
			}
		}
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid authorization format"})
		}

		token := parts[1]
		user, err := authService.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals(UserIDKey, getUserID(user))
		c.Locals(UserKey, user)

		return c.Next()
	}
}

func getUserID(user interface{}) uint {
	switch u := user.(type) {
	case interface{ GetID() uint }:
		return u.GetID()
	default:
		return 0
	}
}

func NewRateLimiter(maxRequests int, expiration int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		KeyGenerator: func(c *fiber.Ctx) string {
			if c.Locals(UserIDKey) != nil {
				return c.Locals(UserIDKey).(string)
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "rate limit exceeded",
				"retry_after": expiration,
			})
		},
	})
}
