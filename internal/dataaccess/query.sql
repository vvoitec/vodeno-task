-- name: InsertCustomer :exec
INSERT INTO customers (email, title, content, mailing_id, insertion_time) VALUES ($1, $2, $3, $4, $5);

-- name: UpsertMailing :exec
INSERT INTO mailings (id) VALUES($1) ON CONFLICT(id) DO NOTHING;

-- name: SelectCustomersByMailingID :many
SELECT * FROM customers WHERE mailing_id = $1 LIMIT $2 OFFSET $3;

-- name: CountCustomersByMailingID :one
SELECT COUNT(*) FROM customers WHERE mailing_id = $1;

-- name: LockMailing :exec
UPDATE mailings SET is_locked = TRUE WHERE id = $1;

-- name: UnlockMailing :exec
UPDATE mailings SET is_locked = FALSE WHERE id = $1;

-- name: IsMailingLocked :one
SELECT is_locked FROM mailings WHERE id = $1;

-- name: DeleteManyCustomers :exec
DELETE FROM customers WHERE id = ANY($1::bigint[]);