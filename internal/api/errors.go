package api

import (
	"github.com/gofiber/fiber/v2"

	appLogger "ogamex-go/internal/logger"
)

type ErrorResponse struct {
	Error       string      `json:"error"`
	Code        string      `json:"code,omitempty"`
	Details     interface{} `json:"details,omitempty"`
	RetryAfter *int        `json:"retry_after,omitempty"`
}

func HandleError(c *fiber.Ctx, status int, message string, err error) error {
	if err != nil {
		appLogger.Error().Err(err).Int("status", status).Str("path", c.Path()).Msg(message)
	} else {
		appLogger.Warn().Int("status", status).Str("path", c.Path()).Msg(message)
	}

	if status == 429 {
		retry := 60
		return c.Status(status).JSON(ErrorResponse{
			Error:       message,
			Code:        "RATE_LIMITED",
			RetryAfter: &retry,
		})
	}

	return c.Status(status).JSON(ErrorResponse{
		Error: message,
		Code:  getErrorCode(status),
	})
}

func HandleNotFound(c *fiber.Ctx) error {
	return HandleError(c, 404, "Resource not found", nil)
}

func HandleBadRequest(c *fiber.Ctx, message string) error {
	return HandleError(c, 400, message, nil)
}

func HandleUnauthorized(c *fiber.Ctx) error {
	return HandleError(c, 401, "Unauthorized", nil)
}

func HandleForbidden(c *fiber.Ctx) error {
	return HandleError(c, 403, "Forbidden", nil)
}

func HandleInternalError(c *fiber.Ctx, err error) error {
	return HandleError(c, 500, "Internal server error", err)
}

func getErrorCode(status int) string {
	codes := map[int]string{
		400: "BAD_REQUEST",
		401: "UNAUTHORIZED",
		403: "FORBIDDEN",
		404: "NOT_FOUND",
		409: "CONFLICT",
		422: "VALIDATION_ERROR",
		429: "RATE_LIMITED",
		500: "INTERNAL_ERROR",
	}
	if code, ok := codes[status]; ok {
		return code
	}
	return "ERROR"
}
