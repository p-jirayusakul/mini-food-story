-- name: ListPaymentMethods :many
SELECT id, code, "name", name_en as "nameEN"
FROM public.md_payment_methods WHERE enable IS TRUE ORDER BY id;