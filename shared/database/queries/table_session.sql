-- name: CreateTableSession :exec
INSERT INTO public.table_session
(id, table_id, number_of_people, session_id, status, started_at, expires_at, ended_at)
VALUES($1, $2, $3, $4, 'active', NOW(), $5, NULL);

-- name: IsTableSessionExists :one
SELECT COUNT(session_id) > 0 as "isExists"
FROM public.table_session
WHERE session_id = sqlc.arg(sessionID)::uuid;

-- name: IsTableSessionActive :one
SELECT COUNT(session_id) > 0 as "isExists"
FROM public.table_session
WHERE session_id = sqlc.arg(sessionID)::uuid
  AND status = 'active';

-- name: GetTableSession :one
SELECT ts.session_id as "sessionID",
       t.id          AS "tableID",
       t.table_number      as "tableNumber",
       ts.status           as "status",
       ts.started_at       as "startedAt",
       o.id                AS "orderID"
FROM public.table_session ts
         JOIN public.tables t ON t.id = ts.table_id
         LEFT JOIN public.orders o ON o.session_id = ts.session_id
WHERE ts.session_id = sqlc.arg(sessionID)::uuid;

-- name: UpdateStatusCloseTableSession :exec
UPDATE public.table_session
SET ended_at=NOW(), status='closed'
WHERE session_id=sqlc.arg(sessionID)::uuid;

-- name: GetSessionIDByTableID :one
select session_id from public.table_session where table_id = sqlc.arg(table_id)::bigint and status = 'active' LIMIT 1;

-- name: UpdateSessionExpireBySessionID :exec
UPDATE public.table_session
SET extend_count =  extend_count + 1,
    extend_total_minutes = extend_total_minutes + sqlc.arg(requested_minutes)::integer,
    last_reason_code = sqlc.arg(last_reason_code)::text,
    lock_version = lock_version + 1,
    expires_at = sqlc.arg(expires_at)::timestamp with time zone
where session_id=sqlc.arg(sessionID)::uuid;

-- name: GetExpiresAtByTableID :one
select expires_at, max_extend_minutes, extend_total_minutes from public.table_session where table_id = sqlc.arg(table_id)::bigint LIMIT 1;
