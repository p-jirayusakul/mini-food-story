package http

import (
	"bufio"
	"context"
	"fmt"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreatePaymentTransaction godoc
// @Summary Create payment transaction
// @Description Create a new payment transaction for an order
// @Tags Payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payment body Payment true "Payment transaction details"
// @Success 201 {object} middleware.SuccessResponse{data=createResponse}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router / [post]
func (s *Handler) CreatePaymentTransaction(c *fiber.Ctx) error {
	body := new(Payment)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	orderID, err := utils.StrToInt64(body.OrderID)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	method, err := utils.StrToInt64(body.Method)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.CreatePaymentTransaction(c.Context(), domain.Payment{
		OrderID: orderID,
		Method:  method,
		Note:    body.Note,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseCreated(c, "create payment transaction success", createResponse{
		TransactionID: result,
	})
}

// CallbackPaymentTransaction godoc
// @Summary Handle payment transaction callback
// @Description Process callback for payment transaction
// @Tags Payment
// @Accept json
// @Produce json
// @Param callback body CallbackPayment true "Payment callback details"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /callback [post]
func (s *Handler) CallbackPaymentTransaction(c *fiber.Ctx) error {
	body := new(CallbackPayment)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError := s.useCase.CallbackPaymentTransaction(c.Context(), body.TransactionID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "payment callback processed successfully", nil)
}

// ListPaymentMethods godoc
// @Summary List payment methods
// @Description Get list of available payment methods
// @Tags Payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} middleware.SuccessResponse{data=[]domain.PaymentMethod}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /methods [get]
func (s *Handler) ListPaymentMethods(c *fiber.Ctx) error {
	result, customError := s.useCase.ListPaymentMethods(c.Context())
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get list payment methods success", result)
}

// StreamPaymentStatusByTransaction godoc
// @Summary List payment methods
// @Description Get list of available payment methods
// @Tags Payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} middleware.SuccessResponse{data=[]domain.PaymentMethod}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /methods [get]
func (s *Handler) StreamPaymentStatusByTransaction(c *fiber.Ctx) error {
	txID := c.Params("transactionID")
	if txID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("missing transactionID")
	}
	if s == nil || s.useCase == nil {
		return c.Status(fiber.StatusInternalServerError).SendString("handler/useCase is nil")
	}

	// SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache, no-transform")
	c.Set("Connection", "keep-alive")
	// (ถ้ามี Nginx) ปิดการ buffer:
	c.Set("X-Accel-Buffering", "no")

	// 🟢 copy ค่าจาก c ออกมาก่อนเข้า writer
	rc := c.Context()          // fasthttp.RequestCtx (เพื่อเอา Done())
	notify := rc.Done()        // <-chan struct{}, ห้ามเรียก c.Context() ใน writer อีก
	baseCtx := c.UserContext() // stdlib context อาจเป็น nil
	if baseCtx == nil {
		baseCtx = context.Background()
	}
	txIDCopy := txID // ป้องกัน capture ตัวแปร outer

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// กัน panic ทั้งหมดใน stream แล้ว log stack
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[SSE panic] tx=%s err=%v\n%s\n", txIDCopy, r, debug.Stack())
				_ = writeSSE(w, "error", "", `{"message":"internal panic"}`)
			}
		}()

		// context สำหรับ DB + lifecycle ของ SSE (อ้างจาก baseCtx ที่ copy ไว้แล้ว)
		ctx, cancel := context.WithCancel(baseCtx)
		defer cancel()

		// --- init event ครั้งแรก ---
		last, err := s.useCase.GetPaymentLastStatusCodeByTransaction(ctx, txIDCopy)
		if err != nil {
			_ = writeSSE(w, "error", "0", `{"message":"get current status failed"}`)
			return
		}
		if err := writeSSE(w, "init", "0", fmt.Sprintf(`{"id":"%s","status":"%s"}`, txIDCopy, last)); err != nil {
			return
		}
		if finalStatus[last] {
			return
		}

		tick := time.NewTicker(2 * time.Second)
		defer tick.Stop()

		heartbeat := time.NewTicker(15 * time.Second)
		defer heartbeat.Stop()

		evID := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-notify:
				// client ปิด → cancel แล้วจบ
				cancel()
				return
			case <-heartbeat.C:
				if err := writeSSE(w, "ping", fmt.Sprint(evID), fmt.Sprintf(`{"ts":%d}`, time.Now().Unix())); err != nil {
					cancel()
					return
				}
			case <-tick.C:
				cur, err := s.useCase.GetPaymentLastStatusCodeByTransaction(ctx, txIDCopy)
				if err != nil {
					_ = writeSSE(w, "error", fmt.Sprint(evID), `{"message":"poll failed"}`)
					cancel()
					return
				}
				if cur != last {
					last = cur
					evID++
					if err := writeSSE(w, "update", fmt.Sprint(evID),
						fmt.Sprintf(`{"id":"%s","status":"%s"}`, txIDCopy, cur)); err != nil {
						cancel()
						return
					}
					if finalStatus[cur] {
						return
					}
				}
			}
		}
	})

	return nil
}

// PaymentTransactionQR godoc
// @Summary List payment methods
// @Description Get list of available payment methods
// @Tags Payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} middleware.SuccessResponse{data=[]domain.PaymentMethod}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /methods [get]
func (s *Handler) PaymentTransactionQR(c *fiber.Ctx) error {
	txID := c.Params("transactionID")
	if txID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("missing transactionID")
	}

	result, customError := s.useCase.PaymentTransactionQR(c.Context(), txID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get payment transaction QR success", result)
}

// ชุดสถานะที่ถือว่า final (ปรับตาม md_payment_statuses ของคุณ)
var finalStatus = map[string]bool{
	"SUCCESS":   true,
	"FAILED":    true,
	"CANCELLED": true,
	"REFUNDED":  true,
}

// ส่ง event ในรูปแบบ SSE มาตรฐาน
func writeSSE(w *bufio.Writer, event, id, data string) error {
	if event != "" {
		if _, err := w.WriteString("event: " + event + "\n"); err != nil {
			return err
		}
	}
	if id != "" {
		if _, err := w.WriteString("id: " + id + "\n"); err != nil {
			return err
		}
	}
	if _, err := w.WriteString("data: " + data + "\n\n"); err != nil {
		return err
	}
	return w.Flush()
}
