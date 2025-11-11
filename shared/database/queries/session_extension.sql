-- name: CreateSessionExtension :one
insert into public.session_extension (id, session_id, requested_minutes, created_at, mode_id, reason_id)
values ($1, $2, $3, NOW(), $4, $5) RETURNING id;