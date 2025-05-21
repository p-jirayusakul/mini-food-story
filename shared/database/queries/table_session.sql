-- name: CreateTableSession :exec
INSERT INTO public.table_session
(id, table_id, number_of_people, session_id, status, started_at, expire_at, ended_at)
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