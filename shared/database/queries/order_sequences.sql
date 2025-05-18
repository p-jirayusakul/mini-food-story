-- name: GetOrCreateOrderSequence :one
INSERT INTO public.order_sequences (order_date, current_number)
VALUES (sqlc.arg(order_date)::date, 1)
ON CONFLICT (order_date) DO UPDATE
    SET current_number = order_sequences.current_number + 1
RETURNING current_number;