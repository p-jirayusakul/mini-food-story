-- name: CreateOrder :one
INSERT INTO public.orders
(id, session_id, table_id, status_id)
VALUES(sqlc.arg(id)::bigint, sqlc.arg(session_id)::uuid, sqlc.arg(table_id)::bigint, (SELECT id FROM public.md_order_statuses WHERE code = 'CONFIRMED' LIMIT 1))
RETURNING id;

-- name: IsOrderExist :one
SELECT COUNT(id) > 0
FROM public.orders WHERE id = $1;

-- name: GetOrderByID :one
SELECT o.id, o.session_id as "sessionID", o.table_id as "tableID", t.table_number as "tableNumber", t.table_number as "tableNumber", o.status_id as "statusID", mos.name as "statusName", mos.name_en as "statusNameEN"
FROM public.orders as o
JOIN public.md_order_statuses as mos ON o.status_id = mos.id
JOIN public.tables as t ON o.table_id = t.id
WHERE o.id = sqlc.arg(id)::bigint;

-- name: UpdateOrderStatus :exec
UPDATE public.orders
SET status_id = (SELECT id FROM public.md_order_statuses WHERE code = sqlc.arg(status_code)::text LIMIT 1)
WHERE id = sqlc.arg(id)::bigint;

-- name: GetOrderWithItems :many
SELECT o.id  AS "orderID",
       oi.id AS "id",
       oi.product_id as "productID",
       oi.product_name as "productName",
       oi.product_name_en as "productNameEN",
       oi.quantity,
       (oi.price * oi.quantity) as "price",
       oi.status_id as "statusID",
       mos.name as "statusName",
       mos.name_en as "statusNameEN",
       mos.code as "statusCode",
       oi.note as "note"
FROM public.orders o
JOIN public.order_items oi ON oi.order_id = o.id
JOIN public.md_order_statuses mos ON oi.status_id = mos.id
WHERE o.id = sqlc.arg(order_id)::bigint
order by oi.created_at DESC;

-- name: GetOrderWithItemsByID :one
SELECT o.id  AS "orderID",
       oi.id AS "id",
       oi.product_id as "productID",
       oi.product_name as "productName",
       oi.product_name_en as "productNameEN",
       oi.quantity,
       (oi.price * oi.quantity) as "price",
       oi.status_id as "statusID",
       mos.name as "statusName",
       mos.name_en as "statusNameEN",
       mos.code as "statusCode",
       oi.note as "note"
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
         JOIN public.md_order_statuses mos ON oi.status_id = mos.id
WHERE o.id = sqlc.arg(order_id)::bigint AND oi.id = sqlc.arg(order_items_id)::bigint LIMIT 1;

-- name: IsOrderWithItemsExists :one
SELECT COUNT(*) > 0
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
WHERE o.id = sqlc.arg(order_id)::bigint AND oi.id = sqlc.arg(order_items_id)::bigint LIMIT 1;


