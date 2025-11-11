-- name: ListCategory :many
SELECT id, "name", name_en as "nameEN", icon_name as "icon" FROM public.md_categories WHERE is_visible IS TRUE ORDER BY sort_order;