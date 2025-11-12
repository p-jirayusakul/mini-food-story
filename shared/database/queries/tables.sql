-- name: CreateTable :one
INSERT INTO public.tables(id, table_number, status_id, seats)
VALUES ($1, $2, (SELECT id FROM public.md_table_statuses WHERE code = 'DISABLED'), $3)
RETURNING id;

-- name: GetTableNumber :one
SELECT table_number as "tableNumber" FROM public.tables WHERE id = $1;

-- name: UpdateTables :exec
UPDATE public.tables
SET table_number=$2, seats=$3, updated_at = NOW()
WHERE id=$1;

-- name: IsTableExists :one
SELECT COUNT(id) > 0 as "isExists" FROM public.tables WHERE id = $1;

-- name: IsTableAvailableOrReserved :one
SELECT COUNT(id) > 0 as "isAvailable" FROM public.tables WHERE id = sqlc.arg(id)::bigint
AND (status_id = (select id from public.md_table_statuses WHERE code = 'AVAILABLE') OR status_id = (select id from public.md_table_statuses WHERE code = 'RESERVED'));

-- name: UpdateTablesStatus :exec
UPDATE public.tables
SET status_id=$2, updated_at = NOW()
WHERE id=$1;

-- name: UpdateTablesStatusWaitingToBeServed :exec
UPDATE public.tables
SET status_id=(select id from public.md_table_statuses WHERE code = 'WAIT_SERVE'), updated_at = NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: UpdateTablesStatusAvailable :exec
UPDATE public.tables
SET status_id=(select id from public.md_table_statuses WHERE code = 'AVAILABLE'), updated_at = NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: UpdateTablesStatusReserved :exec
UPDATE public.tables
SET status_id=(select id from public.md_table_statuses WHERE code = 'RESERVED'), updated_at = NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: UpdateTablesStatusDisabled :exec
UPDATE public.tables
SET status_id=(select id from public.md_table_statuses WHERE code = 'DISABLED'), updated_at = NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: UpdateTablesStatusFoodServed :exec
UPDATE public.tables
SET status_id=(select id from public.md_table_statuses WHERE code = 'FOOD_SERVED'), updated_at = NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: UpdateTablesStatusWaitToOrder :exec
UPDATE public.tables
SET status_id=(select id from public.md_table_statuses WHERE code = 'WAIT_ORDER'), updated_at = NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: UpdateTablesStatusWaitingForPayment :exec
UPDATE public.tables
SET status_id=(select id from public.md_table_statuses WHERE code = 'WAITING_PAYMENT'), updated_at = NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: UpdateTablesStatusCleaning :exec
UPDATE public.tables
SET status_id=(select id from public.md_table_statuses WHERE code = 'CLEANING'), updated_at = NOW()
WHERE id=sqlc.arg(id)::bigint;

-- name: SearchTables :many
SELECT
    t.id,
    t.table_number                         AS "tableNumber",
    s.name                                  AS status,
    s.name_en                               AS "statusEN",
    s.code                                  AS "statusCode",
    t.seats,
    CASE WHEN sess.order_id IS NOT NULL THEN sess.order_id  END AS "orderID",
    sess.expires_at                         AS "expiresAt",
    CASE WHEN sess.extend_total_minutes IS NOT NULL THEN sess.extend_total_minutes  END AS "extendTotalMinutes"
FROM public.tables t
         JOIN public.md_table_statuses s
              ON t.status_id = s.id
         LEFT JOIN LATERAL (
    SELECT
        ts.session_id,
        ts.expires_at,
        ts.extend_total_minutes,
        (
            SELECT o.id
            FROM public.orders o
            WHERE o.session_id = ts.session_id
            ORDER BY o.created_at DESC
            LIMIT 1
        ) AS order_id
    FROM public.table_session ts
    WHERE ts.table_id = t.id
        AND ts.status   = 'active' or ts.status = 'expired'
    ORDER BY ts.started_at DESC
    LIMIT 1
    ) AS sess ON TRUE
WHERE (sqlc.narg(table_number)::int IS NULL OR t.table_number = sqlc.narg(table_number)::int)
  AND (sqlc.narg(seats)::int IS NULL OR t.seats = sqlc.narg(seats)::int)
  AND (
        sqlc.narg(status_code)::varchar[] IS NULL
        OR array_length(sqlc.narg(status_code)::varchar[], 1) = 0
        OR s.code = ANY (sqlc.narg(status_code)::varchar[])
    )
ORDER BY CASE
             WHEN sqlc.arg(order_by_type)::text = 'asc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN t.id::text
                     WHEN sqlc.arg(order_by)::text = 'tableNumber' THEN t.table_number::text
                     WHEN sqlc.arg(order_by)::text = 'seats' THEN t.seats::text
                     WHEN sqlc.arg(order_by)::text = 'status' THEN t.status_id::text
                     ELSE t.table_number::text
                     END
             END,
         CASE
             WHEN sqlc.arg(order_by_type)::text = 'desc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN t.id::text
                     WHEN sqlc.arg(order_by)::text = 'tableNumber' THEN t.table_number::text
                     WHEN sqlc.arg(order_by)::text = 'seats' THEN t.seats::text
                     WHEN sqlc.arg(order_by)::text = 'status' THEN t.status_id::text
                     ELSE t.table_number::text
                     END
             END DESC
OFFSET sqlc.arg(page_number) LIMIT sqlc.arg(page_size);

-- name: GetTotalPageSearchTables :one
SELECT COUNT(*)
FROM public.tables t
         INNER JOIN public.md_table_statuses s ON t.status_id = s.id
WHERE (sqlc.narg(table_number)::int IS NULL OR t.table_number = sqlc.narg(table_number)::int)
  AND (sqlc.narg(seats)::int IS NULL OR t.seats = sqlc.narg(seats)::int)
  AND (
    sqlc.narg(status_code)::varchar[] IS NULL
        OR array_length(sqlc.narg(status_code)::varchar[], 1) = 0
        OR s.code = ANY (sqlc.narg(status_code)::varchar[])
    );

-- name: QuickSearchTables :many
SELECT
    t.id,
    t.table_number                         AS "tableNumber",
    s.name                                  AS status,
    s.name_en                               AS "statusEN",
    s.code                                  AS "statusCode",
    t.seats,
    CASE WHEN sess.order_id IS NOT NULL THEN sess.order_id  END AS "orderID",
    sess.expires_at                         AS "expiresAt",
    CASE WHEN sess.extend_total_minutes IS NOT NULL THEN sess.extend_total_minutes  END AS "extendTotalMinutes"
FROM public.tables t
         JOIN public.md_table_statuses s
              ON t.status_id = s.id
         LEFT JOIN LATERAL (
    SELECT
        ts.session_id,
        ts.expires_at,
        ts.extend_total_minutes,
        (
            SELECT o.id
            FROM public.orders o
            WHERE o.session_id = ts.session_id
            ORDER BY o.created_at DESC
            LIMIT 1
        ) AS order_id
    FROM public.table_session ts
    WHERE ts.table_id = t.id
        AND ts.status   = 'active' or ts.status = 'expired'
    ORDER BY ts.started_at DESC
    LIMIT 1
    ) AS sess ON TRUE
WHERE t.seats >= sqlc.arg(number_of_people)::integer AND s.code = 'AVAILABLE'
ORDER BY CASE
             WHEN sqlc.arg(order_by_type)::text = 'asc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN t.id::text
                     WHEN sqlc.arg(order_by)::text = 'tableNumber' THEN t.table_number::text
                     WHEN sqlc.arg(order_by)::text = 'seats' THEN t.seats::text
                     WHEN sqlc.arg(order_by)::text = 'status' THEN t.status_id::text
                     ELSE t.table_number::text
                     END
             END,
         CASE
             WHEN sqlc.arg(order_by_type)::text = 'desc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN t.id::text
                     WHEN sqlc.arg(order_by)::text = 'tableNumber' THEN t.table_number::text
                     WHEN sqlc.arg(order_by)::text = 'seats' THEN t.seats::text
                     WHEN sqlc.arg(order_by)::text = 'status' THEN t.status_id::text
                     ELSE t.table_number::text
                     END
             END DESC
OFFSET sqlc.arg(page_number) LIMIT sqlc.arg(page_size);

-- name: GetTotalPageQuickSearchTables :one
SELECT COUNT(*) as "totalItems"
FROM public.tables t
         INNER JOIN public.md_table_statuses s ON t.status_id = s.id
WHERE t.seats >= sqlc.arg(number_of_people)::integer AND s.code = 'AVAILABLE';