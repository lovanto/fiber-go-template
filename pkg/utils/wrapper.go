package utils

import (
	"github.com/gofiber/fiber/v2"
)

// SuccessResponse returns a standard success response
func SuccessResponse(ctx *fiber.Ctx, message string, data any) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   message,
		"data":  data,
	})
}

// ErrorResponse returns a standard error response with custom status code
func ErrorResponse(ctx *fiber.Ctx, statusCode int, message string, err error) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"error":   true,
		"msg":     message,
		"details": err.Error(),
	})
}
