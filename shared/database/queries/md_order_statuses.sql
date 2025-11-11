-- name: ListOrderStatus :many
SELECT id, code, "name", name_en
FROM public.md_order_statuses order by sort_order ASC;

-- name: IsOrderStatusExist :one
SELECT (COUNT(id) > 0)  as isExist FROM public.md_order_statuses WHERE code = sqlc.arg(code)::varchar;

-- name: GetOrderStatusPreparing :one
SELECT id FROM public.md_order_statuses WHERE code = 'PREPARING' LIMIT 1;

-- name: GetOrderStatusCompleted :one
SELECT id FROM public.md_order_statuses WHERE code = 'COMPLETED' LIMIT 1;

-- name: IsOrderStatusFinal :one
SELECT COUNT(id) > 0 as "isFinal" FROM public.md_order_statuses WHERE code = sqlc.arg(code)::varchar AND is_final IS TRUE LIMIT 1;