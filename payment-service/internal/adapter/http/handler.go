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

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache, no-transform")
	c.Set("Connection", "keep-alive")

	// 🟢 จับ fasthttp.RequestCtx ไว้ก่อน ห้ามเรียก c.Context() ภายใน writer อีก
	rc := c.Context()
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer func() {
			if r := recover(); r != nil {
				// log ฝั่ง server ให้เห็นจุดพังจริง
				fmt.Printf("[SSE panic] tx=%s err=%v\n%s\n", txID, r, debug.Stack())
				// ส่งข้อความเบา ๆ ไปให้ client
				_ = writeSSE(w, "error", "", `{"message":"internal panic"}`)
			}
		}()

		// base context (อาจเป็น nil)
		base := c.UserContext()
		if base == nil {
			base = context.Background()
		}
		ctx, cancel := context.WithCancel(base)
		defer cancel()

		// 🔔 channel จาก fasthttp ที่จะแจ้งเมื่อ client ปิดการเชื่อมต่อ
		notify := rc.Done()

		// init status
		last, err := s.useCase.GetPaymentLastStatusCodeByTransaction(ctx, txID)
		if err != nil {
			_ = writeSSE(w, "error", "0", `{"message":"get current status failed"}`)
			return
		}
		_ = writeSSE(w, "init", "0", fmt.Sprintf(`{"id":"%s","status":"%s"}`, txID, last))
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
				// client ปิด → ยกเลิก context DB แล้วจบ
				cancel()
				return
			case <-heartbeat.C:
				if err := writeSSE(w, "ping", fmt.Sprint(evID), fmt.Sprintf(`{"ts":%d}`, time.Now().Unix())); err != nil {
					cancel()
					return
				}
			case <-tick.C:
				cur, err := s.useCase.GetPaymentLastStatusCodeByTransaction(ctx, txID)
				if err != nil {
					_ = writeSSE(w, "error", fmt.Sprint(evID), `{"message":"poll failed"}`)
					cancel()
					return
				}
				if cur != last {
					last = cur
					evID++
					if err := writeSSE(w, "update", fmt.Sprint(evID),
						fmt.Sprintf(`{"id":"%s","status":"%s"}`, txID, cur)); err != nil {
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
