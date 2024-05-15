package sqlc

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/vvoitec/vodeno-task/internal/dataaccess"
	"github.com/vvoitec/vodeno-task/internal/mailing"
)

type MailingLocker struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewSQLCMailingLocker(pool *pgxpool.Pool) *MailingLocker {
	return &MailingLocker{queries: sqlc.New(pool), pool: pool}
}

func (l *MailingLocker) Lock(ctx context.Context, m mailing.Mailing) error {
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := l.queries.WithTx(tx)

	isLocked, err := qtx.IsMailingLocked(ctx, int64(m.ID))
	if err != nil {
		return err
	}
	if isLocked {
		return mailing.ErrMailingLocked
	}

	if err = qtx.LockMailing(ctx, int64(m.ID)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (l *MailingLocker) UnLock(ctx context.Context, mailing mailing.Mailing) error {
	return l.queries.UnlockMailing(ctx, int64(mailing.ID))
}
