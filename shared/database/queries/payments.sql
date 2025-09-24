-- name: CreatePayment :one
insert into public.payments (id, order_id, amount, method, status, transaction_id, ref_code, note)
values (sqlc.arg(id)::bigint, sqlc.arg(order_id)::bigint, sqlc.arg(amount)::numeric, sqlc.arg(method)::bigint, sqlc.arg(status)::bigint, sqlc.arg(transaction_id)::text, sqlc.arg(ref_code)::varchar, sqlc.narg(note)::text)
RETURNING id;

-- name: UpdateStatusPaymentSuccessByTransactionID :exec
UPDATE public.payments
SET status=(select id from public.md_payment_statuses WHERE code = 'SUCCESS'), paid_at=NOW(), updated_at=NOW()
WHERE transaction_id=sqlc.arg(transaction_id)::text;

-- name: UpdateStatusPaymentPendingByTransactionID :exec
UPDATE public.payments
SET status=(select id from public.md_payment_statuses WHERE code = 'PENDING'), updated_at=NOW()
WHERE transaction_id=sqlc.arg(transaction_id)::text;

-- name: UpdateStatusPaymentConfirmedByTransactionID :exec
UPDATE public.payments
SET status=(select id from public.md_payment_statuses WHERE code = 'CONFIRMED'), updated_at=NOW()
WHERE transaction_id=sqlc.arg(transaction_id)::text;

-- name: UpdateStatusPaymentCancelledByTransactionID :exec
UPDATE public.payments
SET status=(select id from public.md_payment_statuses WHERE code = 'CANCELLED'), updated_at=NOW()
WHERE transaction_id=sqlc.arg(transaction_id)::text;

-- name: UpdateStatusPaymentFailedByTransactionID :exec
UPDATE public.payments
SET status=(select id from public.md_payment_statuses WHERE code = 'FAILED'), updated_at=NOW()
WHERE transaction_id=sqlc.arg(transaction_id)::text;

-- name: GetPaymentOrderIDByTransaction :one
SELECT order_id as "orderID" FROM public.payments WHERE transaction_id =sqlc.arg(transaction_id)::text LIMIT 1;

-- name: GetPaymentAmountByTransaction :one
SELECT amount::numeric FROM public.payments WHERE transaction_id =sqlc.arg(transaction_id)::text LIMIT 1;

-- name: GetPaymentLastStatusCodeByTransaction :one
select mps.code from public.payments as p
LEFT JOIN public.md_payment_statuses as mps ON mps.id = p.status
where p.transaction_id=$1
limit 1;