package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const CtxRequestIDKey = "requestid"

func RequestIDMiddleware() fiber.Handler {
	return requestid.New(requestid.Config{
		Header:     "X-Request-ID",
		ContextKey: CtxRequestIDKey,
	})
}
