-- name: ListCategory :many
SELECT id, "name", name_en as "nameEN", icon_name as "icon" FROM public.md_categories ORDER BY sort_order;