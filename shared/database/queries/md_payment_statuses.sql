-- name: GetPaymentStatusPending :one
SELECT id FROM public.md_payment_statuses WHERE code = 'PENDING' LIMIT 1;