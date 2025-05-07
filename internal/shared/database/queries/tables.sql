-- name: CreateTable :one
INSERT INTO public.tables(id, table_number, status_id, seats)
VALUES ($1, $2, (SELECT id FROM public.md_table_statuses WHERE code = 'DISABLED'), $3)
RETURNING id;

-- name: SearchTables :many
SELECT t.id, t.table_number as "tableNumber", s.name as status, s.name_en as "statusEN", t.seats
FROM public.tables t
         INNER JOIN public.md_table_statuses s ON t.status_id = s.id
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
SELECT t.id, t.table_number as "tableNumber", s.name as status, s.name_en as "statusEN", t.seats
FROM public.tables t
         INNER JOIN public.md_table_statuses s ON t.status_id = s.id
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
SELECT COUNT(*)
FROM public.tables t
         INNER JOIN public.md_table_statuses s ON t.status_id = s.id
WHERE t.seats >= sqlc.arg(number_of_people)::integer AND s.code = 'AVAILABLE';