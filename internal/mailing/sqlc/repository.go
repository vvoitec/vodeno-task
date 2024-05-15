package sqlc

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvoitec/vodeno-task/internal/dataaccess"
	"github.com/vvoitec/vodeno-task/internal/mailing"
)

type CustomerSQLCRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewCustomerSQLCRepository(pool *pgxpool.Pool) *CustomerSQLCRepository {
	return &CustomerSQLCRepository{queries: sqlc.New(pool), pool: pool}
}

func (r *CustomerSQLCRepository) SaveCustomer(ctx context.Context, c mailing.Customer) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)

	isLocked, err := qtx.IsMailingLocked(ctx, int64(c.Mailing.ID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			isLocked = false
		} else {
			return err
		}
	}
	if isLocked {
		return mailing.ErrMailingLocked
	}

	if err = qtx.UpsertMailing(ctx, int64(c.Mailing.ID)); err != nil {
		return err
	}

	if err = qtx.InsertCustomer(ctx, sqlc.InsertCustomerParams{
		Email:         c.Email,
		Title:         pgtype.Text{String: c.Title, Valid: true},
		Content:       pgtype.Text{String: c.Title, Valid: true},
		MailingID:     pgtype.Int8{Int64: int64(c.Mailing.ID), Valid: true},
		InsertionTime: pgtype.Timestamp{Time: c.InsertionTime, Valid: true},
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *CustomerSQLCRepository) FindManyByMailing(ctx context.Context, m mailing.Mailing, limit int, offset int) ([]mailing.Customer, error) {
	sqlcCustomers, err := r.queries.SelectCustomersByMailingID(ctx, sqlc.SelectCustomersByMailingIDParams{
		MailingID: pgtype.Int8{Int64: int64(m.ID), Valid: true},
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, err
	}

	customers := make([]mailing.Customer, len(sqlcCustomers))
	for i, s := range sqlcCustomers {
		customers[i] = mailing.Customer{
			ID:            uint(s.ID),
			Email:         s.Email,
			Title:         s.Title.String,
			Content:       s.Content.String,
			Mailing:       mailing.Mailing{ID: uint(s.MailingID.Int64)},
			InsertionTime: s.InsertionTime.Time,
		}
	}

	return customers, nil
}

func (r *CustomerSQLCRepository) CountByMailing(ctx context.Context, mailing mailing.Mailing) (uint, error) {
	count, err := r.queries.CountCustomersByMailingID(ctx, pgtype.Int8{Int64: int64(mailing.ID), Valid: true})
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

func (r *CustomerSQLCRepository) DeleteManyCustomers(ctx context.Context, customers []mailing.Customer) error {
	ids := make([]int64, len(customers))
	for i, c := range customers {
		ids[i] = int64(c.ID)
	}

	return r.queries.DeleteManyCustomers(ctx, ids)
}
