-- name: CreateOrder :one
INSERT INTO public.orders
(id, order_number, session_id, table_id, status_id)
VALUES(sqlc.arg(id)::bigint, sqlc.arg(order_number)::varchar, sqlc.arg(session_id)::uuid, sqlc.arg(table_id)::bigint, (SELECT id FROM public.md_order_statuses WHERE code = 'CONFIRMED' LIMIT 1))
RETURNING id;

-- name: IsOrderExist :one
SELECT COUNT(id) > 0
FROM public.orders WHERE id = $1;

-- name: GetOrderByID :one
SELECT o.id, o.session_id as "sessionID", o.table_id as "tableID", t.table_number as "tableNumber", o.status_id as "statusID", mos.name as "statusName", mos.name_en as "statusNameEN", mos.code as "statusCode"
FROM public.orders as o
JOIN public.md_order_statuses as mos ON o.status_id = mos.id
JOIN public.tables as t ON o.table_id = t.id
WHERE o.id = sqlc.arg(id)::bigint;

-- name: UpdateOrderStatus :exec
UPDATE public.orders
SET status_id = (SELECT id FROM public.md_order_statuses WHERE code = sqlc.arg(status_code)::text LIMIT 1)
WHERE id = sqlc.arg(id)::bigint;

-- name: UpdateOrderStatusWaitForPayment :exec
UPDATE public.orders
SET status_id = (SELECT id FROM public.md_order_statuses WHERE code = 'WAITING_PAYMENT' LIMIT 1)
WHERE id = sqlc.arg(id)::bigint;

-- name: UpdateOrderStatusCompletedAndAmount :exec
UPDATE public.orders
SET total_amount =sqlc.arg(amount)::numeric, status_id = (SELECT id FROM public.md_order_statuses WHERE code = 'COMPLETED' LIMIT 1)
WHERE id = sqlc.arg(id)::bigint;

-- name: SearchOrderItems :many
SELECT o.id  AS "orderID",
       o.order_number as "orderNumber",
       oi.id AS "id",
       oi.product_id as "productID",
       oi.product_name as "productName",
       oi.product_name_en as "productNameEN",
       t.table_number as "tableNumber",
       oi.quantity,
       (oi.price * oi.quantity) as "price",
       oi.status_id as "statusID",
       mos.name as "statusName",
       mos.name_en as "statusNameEN",
       mos.code as "statusCode",
       oi.note as "note",
       oi.created_at
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
         JOIN public.md_order_statuses mos ON oi.status_id = mos.id
         JOIN public.tables t ON o.table_id = t.id
WHERE  DATE(oi.created_at) = CURRENT_DATE
  AND ((sqlc.narg(product_name)::varchar IS NULL OR oi."product_name" ILIKE '%' || sqlc.narg(product_name)::varchar || '%') OR (sqlc.narg(product_name)::varchar IS NULL OR oi.product_name_en ILIKE '%' || sqlc.narg(product_name)::varchar || '%'))
  AND (
    sqlc.narg(table_number)::int[] IS NULL
        OR array_length(sqlc.narg(table_number)::int[], 1) = 0
        OR t.table_number = ANY (sqlc.narg(table_number)::int[])
    )
  AND (
    sqlc.narg(status_code)::varchar[] IS NULL
        OR array_length(sqlc.narg(status_code)::varchar[], 1) = 0
        OR mos.code = ANY (sqlc.narg(status_code)::varchar[])
    )
ORDER BY CASE
             WHEN sqlc.arg(order_by_type)::text = 'asc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN oi.id::text
                     WHEN sqlc.arg(order_by)::text = 'tableNumber' THEN t."table_number"::text
                     WHEN sqlc.arg(order_by)::text = 'statusCode' THEN mos."code"::text
                     WHEN sqlc.arg(order_by)::text = 'productName' THEN oi."product_name"::text
                     WHEN sqlc.arg(order_by)::text = 'quantity' THEN oi."quantity"::text
                     ELSE oi.id::text
                     END
             END,
         CASE
             WHEN sqlc.arg(order_by_type)::text = 'desc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN oi.id::text
                     WHEN sqlc.arg(order_by)::text = 'tableNumber' THEN t."table_number"::text
                     WHEN sqlc.arg(order_by)::text = 'statusCode' THEN mos."code"::text
                     WHEN sqlc.arg(order_by)::text = 'productName' THEN oi."product_name"::text
                     WHEN sqlc.arg(order_by)::text = 'quantity' THEN oi."quantity"::text
                     ELSE oi.id::text
                     END
             END DESC
OFFSET sqlc.arg(page_number) LIMIT sqlc.arg(page_size);

-- name: GetTotalSearchOrderItems :one
SELECT COUNT(*)
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
         JOIN public.md_order_statuses mos ON oi.status_id = mos.id
         JOIN public.tables t ON o.table_id = t.id
WHERE  DATE(oi.created_at) = CURRENT_DATE
  AND ((sqlc.narg(product_name)::varchar IS NULL OR oi."product_name" ILIKE '%' || sqlc.narg(product_name)::varchar || '%') OR (sqlc.narg(product_name)::varchar IS NULL OR oi.product_name_en ILIKE '%' || sqlc.narg(product_name)::varchar || '%'))
  AND (
    sqlc.narg(table_number)::int[] IS NULL
        OR array_length(sqlc.narg(table_number)::int[], 1) = 0
        OR t.table_number = ANY (sqlc.narg(table_number)::int[])
    )
  AND (
    sqlc.narg(status_code)::varchar[] IS NULL
        OR array_length(sqlc.narg(status_code)::varchar[], 1) = 0
        OR mos.code = ANY (sqlc.narg(status_code)::varchar[])
    );

-- name: GetTableNumberOrderByID :one
SELECT t.table_number
FROM public.orders o
         JOIN public.tables t ON o.table_id = t.id
WHERE o.id = sqlc.arg(order_id)::bigint LIMIT 1;

-- name: GetOrderWithItems :many
SELECT o.id  AS "orderID",
       o.order_number as "orderNumber",
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
       oi.note as "note",
       oi.created_at,
       t.table_number as "tableNumber"
FROM public.orders o
JOIN public.order_items oi ON oi.order_id = o.id
JOIN public.md_order_statuses mos ON oi.status_id = mos.id
JOIN public.tables t ON o.table_id = t.id
WHERE o.id = sqlc.arg(order_id)::bigint
order by oi.id DESC;

-- name: GetOrderWithItemsByID :one
SELECT o.id  AS "orderID",
       o.order_number as "orderNumber",
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
       oi.note as "note",
       oi.created_at,
       t.table_number as "tableNumber"
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
         JOIN public.md_order_statuses mos ON oi.status_id = mos.id
         JOIN public.tables t ON o.table_id = t.table_number
WHERE o.id = sqlc.arg(order_id)::bigint AND oi.id = sqlc.arg(order_items_id)::bigint LIMIT 1;

-- name: GetOrderWithItemsGroupID :many
SELECT o.id  AS "orderID",
       o.order_number as "orderNumber",
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
       oi.note as "note",
       oi.created_at,
       t.table_number as "tableNumber"
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
         JOIN public.md_order_statuses mos ON oi.status_id = mos.id
        JOIN public.tables t ON o.table_id = t.id
WHERE oi.id = ANY(sqlc.arg(order_items_id)::bigint[])
order by oi.id DESC;

-- name: IsOrderWithItemsExists :one
SELECT COUNT(*) > 0
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
WHERE o.id = sqlc.arg(order_id)::bigint AND oi.id = sqlc.arg(order_items_id)::bigint LIMIT 1;

-- name: SearchOrderItemsIsNotFinal :many
SELECT o.id  AS "orderID",
       o.order_number as "orderNumber",
       oi.id AS "id",
       oi.product_id as "productID",
       oi.product_name as "productName",
       oi.product_name_en as "productNameEN",
       t.table_number as "tableNumber",
       oi.quantity,
       (oi.price * oi.quantity) as "price",
       oi.status_id as "statusID",
       mos.name as "statusName",
       mos.name_en as "statusNameEN",
       mos.code as "statusCode",
       oi.note as "note",
       oi.created_at
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
         JOIN public.md_order_statuses mos ON oi.status_id = mos.id
         JOIN public.tables t ON o.table_id = t.id
WHERE o.id = sqlc.arg(order_id)::bigint AND (mos.code != 'SERVED' AND mos.code != 'CANCELLED')
  AND ((sqlc.narg(product_name)::varchar IS NULL OR oi."product_name" ILIKE '%' || sqlc.narg(product_name)::varchar || '%') OR (sqlc.narg(product_name)::varchar IS NULL OR oi.product_name_en ILIKE '%' || sqlc.narg(product_name)::varchar || '%'))
  AND (
    sqlc.narg(status_code)::varchar[] IS NULL
        OR array_length(sqlc.narg(status_code)::varchar[], 1) = 0
        OR mos.code = ANY (sqlc.narg(status_code)::varchar[])
    )
ORDER BY CASE
             WHEN sqlc.arg(order_by_type)::text = 'asc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN oi.id::text
                     WHEN sqlc.arg(order_by)::text = 'statusCode' THEN mos."code"::text
                     WHEN sqlc.arg(order_by)::text = 'productName' THEN oi."product_name"::text
                     WHEN sqlc.arg(order_by)::text = 'quantity' THEN oi."quantity"::text
                     ELSE oi.id::text
                     END
             END,
         CASE
             WHEN sqlc.arg(order_by_type)::text = 'desc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN oi.id::text
                     WHEN sqlc.arg(order_by)::text = 'statusCode' THEN mos."code"::text
                     WHEN sqlc.arg(order_by)::text = 'productName' THEN oi."product_name"::text
                     WHEN sqlc.arg(order_by)::text = 'quantity' THEN oi."quantity"::text
                     ELSE oi.id::text
                     END
             END DESC
OFFSET sqlc.arg(page_number) LIMIT sqlc.arg(page_size);


-- name: GetTotalSearchOrderItemsIsNotFinal :one
SELECT COUNT(*)
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
         JOIN public.md_order_statuses mos ON oi.status_id = mos.id
         JOIN public.tables t ON o.table_id = t.id
WHERE o.id = sqlc.arg(order_id)::bigint AND (mos.code != 'SERVED' AND mos.code != 'CANCELLED')
  AND ((sqlc.narg(product_name)::varchar IS NULL OR oi."product_name" ILIKE '%' || sqlc.narg(product_name)::varchar || '%') OR (sqlc.narg(product_name)::varchar IS NULL OR oi.product_name_en ILIKE '%' || sqlc.narg(product_name)::varchar || '%'))
  AND (
    sqlc.narg(status_code)::varchar[] IS NULL
        OR array_length(sqlc.narg(status_code)::varchar[], 1) = 0
        OR mos.code = ANY (sqlc.narg(status_code)::varchar[])
    );

-- name: IsOrderItemsNotFinal :one
SELECT COUNT(*) > 0
FROM public.orders o
         JOIN public.order_items oi ON oi.order_id = o.id
         JOIN public.md_order_statuses mos ON oi.status_id = mos.id
WHERE o.id = sqlc.arg(order_id)::bigint AND (mos.code != 'SERVED' AND mos.code != 'CANCELLED');