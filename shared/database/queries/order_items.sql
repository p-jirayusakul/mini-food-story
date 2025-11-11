-- name: CreateOrderItems :copyfrom
INSERT INTO public.order_items
(id, order_id, product_id, status_id, product_name, product_name_en, price, quantity, note, created_at, product_image_url, is_visible)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: CreateOrderItemsPerRow :exec
INSERT INTO public.order_items
(id, order_id, product_id, status_id, product_name, product_name_en, price, quantity, note, created_at, product_image_url, is_visible)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);


-- name: GetOrderItemsByID :one
SELECT id, order_id, product_id, status_id, product_name, product_name_en, price, quantity, note
FROM public.order_items WHERE id = $1 AND is_visible IS TRUE;

-- name: IsOrderItemsExist :one
SELECT COUNT(id) > 0 as "isExist"
FROM public.order_items WHERE id = $1;

-- name: UpdateOrderItemsStatus :exec
UPDATE public.order_items
SET status_id = (SELECT id FROM public.md_order_statuses WHERE code = sqlc.arg(status_code)::text LIMIT 1), updated_at = NOW()
WHERE id = sqlc.arg(id)::bigint;

-- name: UpdateOrderItemsStatusServed :exec
UPDATE public.order_items
SET status_id = (SELECT id FROM public.md_order_statuses WHERE code = 'SERVED' LIMIT 1), prepared_at=NOW(), updated_at = NOW()
WHERE id = sqlc.arg(id)::bigint;

-- name: GetTotalAmountToPayForServedItems :one
SELECT SUM(price * quantity) AS "totalAmount"
FROM public.order_items
WHERE order_id = $1 AND status_id = (SELECT id FROM public.md_order_statuses WHERE code = 'SERVED' LIMIT 1);