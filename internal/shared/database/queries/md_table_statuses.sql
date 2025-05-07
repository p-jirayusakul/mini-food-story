-- name: ListTableStatus :many
SELECT id, code, name, name_en FROM public.md_table_statuses
ORDER BY id;

-- name: CreateTableStatus :one
INSERT INTO public.md_table_statuses(
    id, code, name, name_en)
VALUES ($1, $2, $3, $4)
RETURNING id;