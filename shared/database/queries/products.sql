-- name: CreateProduct :one
INSERT INTO public.products
(id, "name", name_en, categories, description, price, is_available, image_url)
VALUES($1, $2, $3, $4, $5, $6, $7, $8 )
RETURNING id;

-- name: UpdateProduct :exec
UPDATE public.products
SET "name"       = $2,
    name_en      = $3,
    categories   = $4,
    description  = $5,
    price        = $6,
    is_available = $7,
    image_url    = $8,
    updated_at   = NOW()
WHERE id = $1;

-- name: UpdateProductAvailability :exec
UPDATE public.products
SET is_available = $2,
    updated_at   = NOW()
WHERE id = $1;

-- name: SearchProducts :many
SELECT p.id,
       p."name",
       p.name_en,
       p.categories,
       c.name as "categoryName",
       c.name_en as "categoryNameEN",
       p.description,
       p.price,
       p.is_available,
       p.image_url
FROM public.products as p
         INNER JOIN public.md_categories as c ON c.id = p.categories
WHERE ((sqlc.narg(name)::varchar IS NULL OR p."name" ILIKE '%' || sqlc.narg(name)::varchar || '%') OR (sqlc.narg(name)::varchar IS NULL OR p.name_en ILIKE '%' || sqlc.narg(name)::varchar || '%'))
  AND (sqlc.narg(is_available)::boolean IS NULL OR p.is_available = sqlc.narg(is_available)::boolean)
  AND (
    sqlc.narg(category_id)::bigint[] IS NULL
        OR array_length(sqlc.narg(category_id)::bigint[], 1) = 0
        OR p.categories = ANY (sqlc.narg(category_id)::bigint[])
    )
ORDER BY CASE
             WHEN sqlc.arg(order_by_type)::text = 'asc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN p.id::text
                     WHEN sqlc.arg(order_by)::text = 'name' THEN p."name"
                     WHEN sqlc.arg(order_by)::text = 'price' THEN p.price::text
                     ELSE p.id::text
                     END
             END,
         CASE
             WHEN sqlc.arg(order_by_type)::text = 'desc' THEN
                 CASE
                     WHEN sqlc.arg(order_by)::text = 'id' THEN p.id::text
                     WHEN sqlc.arg(order_by)::text = 'name' THEN p."name"
                     WHEN sqlc.arg(order_by)::text = 'price' THEN p.price::text
                     ELSE p.id::text
                     END
             END DESC
OFFSET sqlc.arg(page_number) LIMIT sqlc.arg(page_size);

-- name: GetTotalPageSearchProducts :one
SELECT COUNT(*)
FROM public.products as p
         INNER JOIN public.md_categories as c ON c.id = p.categories
WHERE ((sqlc.narg(name)::varchar IS NULL OR p."name" ILIKE '%' || sqlc.narg(name)::varchar || '%') OR (sqlc.narg(name)::varchar IS NULL OR p.name_en ILIKE '%' || sqlc.narg(name)::varchar || '%'))
  AND (sqlc.narg(is_available)::boolean IS NULL OR p.is_available = sqlc.narg(is_available)::boolean)
  AND (
    sqlc.narg(category_id)::bigint[] IS NULL
        OR array_length(sqlc.narg(category_id)::bigint[], 1) = 0
        OR p.categories = ANY (sqlc.narg(category_id)::bigint[])
    );

-- name: GetProductByID :one
SELECT p.id,
       p."name",
       p.name_en,
       p.categories,
       c.name as "categoryName",
       c.name_en as "categoryNameEN",
       p.description,
       p.price,
       p.is_available,
       p.image_url
FROM public.products as p
         INNER JOIN public.md_categories as c ON c.id = p.categories
WHERE p.id = sqlc.arg(id)::bigint LIMIT 1;

-- name: GetProductAvailableByID :one
SELECT p.id,
       p."name",
       p.name_en,
       p.categories,
       c.name as "categoryName",
       c.name_en as "categoryNameEN",
       p.description,
       p.price,
       p.is_available,
       p.image_url
FROM public.products as p
         INNER JOIN public.md_categories as c ON c.id = p.categories
WHERE p.id = sqlc.arg(id)::bigint AND p.is_available IS TRUE LIMIT 1;

-- name: IsProductExists :one
SELECT count(*) > 0 FROM public.products WHERE id = $1;