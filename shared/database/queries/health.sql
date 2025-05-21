-- name: Health :one
SELECT 1 as "healthy";

-- name: GetTimeNow :one
select NOW() as "today";