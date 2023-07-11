package httpgin

import "github.com/gofiber/fiber/v2"

func NewErrResponse(err error) *fiber.Map {
	return &fiber.Map{
		"error": err.Error(),
	}
}
