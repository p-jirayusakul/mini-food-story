-- name: GetDurationMinutesByProductID :one
select duration_minutes from public.product_time_extension WHERE products_id = sqlc.arg(products_id)::bigint limit 1;