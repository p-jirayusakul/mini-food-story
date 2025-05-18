-- name: ListPaymentMethods :many
SELECT id, code, "name", name_en as "nameEN"
FROM public.payment_methods WHERE enable IS TRUE ORDER BY id;