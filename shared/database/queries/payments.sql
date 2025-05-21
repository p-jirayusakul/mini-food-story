-- name: CreatePayment :one
INSERT INTO public.payments
(id, order_id, amount, "method", transaction_id, ref_code, note)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: UpdateStatusPaymentPaidByID :exec
UPDATE public.payments
SET status='paid'::payment_status, paid_at=NOW(), updated_at=NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: UpdateStatusPaymentPaidByTransactionID :exec
UPDATE public.payments
SET status='paid'::payment_status, paid_at=NOW(), updated_at=NOW()
WHERE transaction_id=sqlc.arg(transaction_id)::text;

-- name: UpdateStatusPaymentFail :exec
UPDATE public.payments
SET status='failed'::payment_status, paid_at=NULL, updated_at=NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: GetPaymentOrderIDByTransaction :one
SELECT order_id as "orderID" FROM public.payments WHERE transaction_id =sqlc.arg(transaction_id)::text LIMIT 1;