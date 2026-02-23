package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const TraceIDKey = "trace_id"

func TraceIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Locals(TraceIDKey, traceID)
		c.Set("X-Trace-ID", traceID)
		return c.Next()
	}
}

func GetTraceID(c *fiber.Ctx) string {
	if traceID, ok := c.Locals(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}
