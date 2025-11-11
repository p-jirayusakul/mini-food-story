-- name: ListSessionExtensionReason :many
select
    id,
    code,
    name,
    name_en as "nameEN",
    category,
    mode_code as "modeCode"
from public.md_session_extension_reason order by sort_order;

-- name: GetSessionExtensionModeByReasonCode :one
select m.id, sem.id as "sessionExtensionModeID"
from public.md_session_extension_reason as m
         left join public.md_session_extension_mode as sem ON sem.code = m.mode_code
where m.code = $1 limit 1;

-- name: IsSessionExtensionReasonExist :one
select id from public.md_session_extension_reason where id = $1 limit 1;
