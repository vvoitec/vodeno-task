package mailing

import (
	"context"
	"errors"
	"time"
)

var (
	ErrMailingLocked    = errors.New("mailing is locked")
	ErrMailingEmpty     = errors.New("mailing is empty")
	SendMailingJobLabel = "mailing-send-job"
)

type (
	Customer struct {
		ID            uint
		Email         string
		Title         string
		Content       string
		Mailing       Mailing
		InsertionTime time.Time
	}

	Mailing struct {
		ID uint
	}

	Repository interface {
		// SaveCustomer persists customer to the database.
		SaveCustomer(ctx context.Context, customer Customer) error
		// FindManyByMailing finds paginated customers for a mailing list.
		FindManyByMailing(ctx context.Context, mailing Mailing, limit int, offset int) ([]Customer, error)
		// CountByMailing counts customers belonging to a mailing list.
		CountByMailing(ctx context.Context, mailing Mailing) (uint, error)
		// DeleteManyCustomers deletes multiple customers.
		DeleteManyCustomers(ctx context.Context, customers []Customer) error
	}

	Locker interface {
		// Lock locks mailing to prevent data races.
		Lock(context.Context, Mailing) error
		// UnLock unlocks mailing.
		UnLock(context.Context, Mailing) error
	}

	Sender interface {
		// SendTo sends emails to customers on mailing list.
		SendTo(context.Context, Mailing) error
	}
)
