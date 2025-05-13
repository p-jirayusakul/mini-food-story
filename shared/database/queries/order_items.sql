-- name: CreateOrderItems :copyfrom
INSERT INTO public.order_items
(id, order_id, product_id, status_id, product_name, product_name_en, price, quantity, note)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetOrderItemsByID :one
SELECT id, order_id, product_id, status_id, product_name, product_name_en, price, quantity, note
FROM public.order_items WHERE id = $1;

-- name: IsOrderItemsExist :one
SELECT COUNT(id) > 0
FROM public.order_items WHERE id = $1;

-- name: UpdateOrderItemsStatus :exec
UPDATE public.order_items
SET status_id = (SELECT id FROM public.md_order_statuses WHERE code = sqlc.arg(status_code)::text LIMIT 1)
WHERE id = sqlc.arg(id)::bigint;