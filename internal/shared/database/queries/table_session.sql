-- name: CreateTableSession :one
INSERT INTO public.table_session
(id, table_id, number_of_people, status, started_at, expire_at, ended_at)
VALUES($1, $2, $3, 'active', NOW(), $4, NULL)
RETURNING session_id;

