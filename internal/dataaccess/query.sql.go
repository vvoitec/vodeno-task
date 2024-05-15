// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countCustomersByMailingID = `-- name: CountCustomersByMailingID :one
SELECT COUNT(*) FROM customers WHERE mailing_id = $1
`

func (q *Queries) CountCustomersByMailingID(ctx context.Context, mailingID pgtype.Int8) (int64, error) {
	row := q.db.QueryRow(ctx, countCustomersByMailingID, mailingID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteManyCustomers = `-- name: DeleteManyCustomers :exec
DELETE FROM customers WHERE id = ANY($1::bigint[])
`

func (q *Queries) DeleteManyCustomers(ctx context.Context, dollar_1 []int64) error {
	_, err := q.db.Exec(ctx, deleteManyCustomers, dollar_1)
	return err
}

const insertCustomer = `-- name: InsertCustomer :exec
INSERT INTO customers (email, title, content, mailing_id, insertion_time) VALUES ($1, $2, $3, $4, $5)
`

type InsertCustomerParams struct {
	Email         string
	Title         pgtype.Text
	Content       pgtype.Text
	MailingID     pgtype.Int8
	InsertionTime pgtype.Timestamp
}

func (q *Queries) InsertCustomer(ctx context.Context, arg InsertCustomerParams) error {
	_, err := q.db.Exec(ctx, insertCustomer,
		arg.Email,
		arg.Title,
		arg.Content,
		arg.MailingID,
		arg.InsertionTime,
	)
	return err
}

const isMailingLocked = `-- name: IsMailingLocked :one
SELECT is_locked FROM mailings WHERE id = $1
`

func (q *Queries) IsMailingLocked(ctx context.Context, id int64) (bool, error) {
	row := q.db.QueryRow(ctx, isMailingLocked, id)
	var is_locked bool
	err := row.Scan(&is_locked)
	return is_locked, err
}

const lockMailing = `-- name: LockMailing :exec
UPDATE mailings SET is_locked = TRUE WHERE id = $1
`

func (q *Queries) LockMailing(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, lockMailing, id)
	return err
}

const selectCustomersByMailingID = `-- name: SelectCustomersByMailingID :many
SELECT id, email, title, content, mailing_id, insertion_time FROM customers WHERE mailing_id = $1 LIMIT $2 OFFSET $3
`

type SelectCustomersByMailingIDParams struct {
	MailingID pgtype.Int8
	Limit     int32
	Offset    int32
}

func (q *Queries) SelectCustomersByMailingID(ctx context.Context, arg SelectCustomersByMailingIDParams) ([]Customer, error) {
	rows, err := q.db.Query(ctx, selectCustomersByMailingID, arg.MailingID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Customer
	for rows.Next() {
		var i Customer
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Title,
			&i.Content,
			&i.MailingID,
			&i.InsertionTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unlockMailing = `-- name: UnlockMailing :exec
UPDATE mailings SET is_locked = FALSE WHERE id = $1
`

func (q *Queries) UnlockMailing(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, unlockMailing, id)
	return err
}

const upsertMailing = `-- name: UpsertMailing :exec
INSERT INTO mailings (id) VALUES($1) ON CONFLICT(id) DO NOTHING
`

func (q *Queries) UpsertMailing(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, upsertMailing, id)
	return err
}
