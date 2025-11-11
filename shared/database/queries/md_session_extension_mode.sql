-- name: IsSessionExtensionModeFree :one
select COUNT(id) > 0 from public.md_session_extension_mode where id = $1 AND code = 'COMP_FREE' limit 1;