package smtp

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vvoitec/vodeno-task/internal/mailing"
	"sync"
	"testing"
	"time"
)

func TestSender_SendTo(t *testing.T) {
	locker := &LockerMock{}
	repository := &RepositoryMock{}
	mailerMock := &MailerMock{}
	config := Config{}
	logger := &logrus.Logger{}

	sender := NewSender(repository, locker, config, logger, mailerMock)

	t.Run("test that it sends mails", func(t *testing.T) {
		// given
		m := mailing.Mailing{ID: 1}
		require.NoError(t, repository.SaveCustomer(context.Background(), mailing.Customer{
			ID:            1,
			Email:         "example.com",
			Title:         "title",
			Content:       "",
			Mailing:       m,
			InsertionTime: time.Now(),
		}))
		require.Equal(t, 0, mailerMock.count)

		// when
		require.NoError(t, sender.SendTo(context.Background(), m))

		// then
		assert.Equal(t, 1, mailerMock.count)
	})
}

type (
	LockerMock struct {
	}

	RepositoryMock struct {
		storage sync.Map
	}

	MailerMock struct {
		count int
	}
)

func (m *MailerMock) SendMail(to []string, msg []byte) error {
	m.count++
	return nil
}

func (l *LockerMock) Lock(context.Context, mailing.Mailing) error {
	return nil
}

func (l *LockerMock) UnLock(context.Context, mailing.Mailing) error {
	return nil
}

func (r *RepositoryMock) SaveCustomer(ctx context.Context, customer mailing.Customer) error {
	r.storage.Store(customer.ID, customer)
	return nil
}

func (r *RepositoryMock) FindManyByMailing(ctx context.Context, m mailing.Mailing, limit int, offset int) ([]mailing.Customer, error) {
	customers := make([]mailing.Customer, 0)
	r.storage.Range(func(key, value any) bool {
		v := value.(mailing.Customer)
		if v.Mailing.ID == m.ID {
			customers = append(customers, v)
		}
		return true
	})

	return customers, nil
}

func (r *RepositoryMock) CountByMailing(ctx context.Context, m mailing.Mailing) (uint, error) {
	count := 0
	r.storage.Range(func(key, value any) bool {
		v := value.(mailing.Customer)
		if v.Mailing.ID == m.ID {
			count++
		}
		return true
	})

	return uint(count), nil
}

func (r *RepositoryMock) DeleteManyCustomers(ctx context.Context, customers []mailing.Customer) error {
	for _, c := range customers {
		r.storage.Delete(c.ID)
	}
	return nil
}
