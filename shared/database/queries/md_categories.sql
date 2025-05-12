-- name: ListCategory :many
SELECT id, "name", name_en as "nameEN" FROM public.md_categories ORDER BY id DESC;