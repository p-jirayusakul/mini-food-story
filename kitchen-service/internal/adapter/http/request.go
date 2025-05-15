package http

type StatusOrderItems struct {
	StatusCode string `json:"statusCode" validate:"required,oneof=SERVED CANCELLED"`
}
