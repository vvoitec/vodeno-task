package smtp

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vvoitec/vodeno-task/internal/mailing"
	"golang.org/x/sync/errgroup"
	"math"
	"net/smtp"
	"runtime"
)

var batchSize = 10

type (
	Sender struct {
		repository mailing.Repository
		locker     mailing.Locker
		config     Config
		logger     *logrus.Logger
		mailer     mailer
	}

	Config struct {
		SMTPURL string
	}

	mailer interface {
		SendMail(to []string, msg []byte) error
	}

	defaultMailer struct {
		config Config
	}
)

func (m *defaultMailer) SendMail(to []string, msg []byte) error {
	return smtp.SendMail(m.config.SMTPURL, nil, "vodeno-task@example.com", to, msg)
}

func NewSender(
	repository mailing.Repository,
	locker mailing.Locker,
	config Config,
	logger *logrus.Logger,
	mailer mailer,
) *Sender {
	if mailer == nil {
		mailer = &defaultMailer{config: config}
	}
	return &Sender{
		repository: repository,
		locker:     locker,
		config:     config,
		logger:     logger,
		mailer:     mailer,
	}
}

// SendTo sends emails to mailing list concurrently in batches configured by batchSize
func (s *Sender) SendTo(ctx context.Context, m mailing.Mailing) error {
	if err := s.locker.Lock(ctx, m); err != nil {
		return err
	}
	defer func() {
		if err := s.locker.UnLock(context.Background(), m); err != nil {
			s.logger.WithError(err).WithField("mailing_id", m.ID).Error("failed to unlock mailing")
		}
	}()

	msgCount, err := s.repository.CountByMailing(ctx, m)
	if err != nil {
		return err
	}
	if msgCount == 0 {
		return mailing.ErrMailingEmpty
	}
	batches := int(math.Ceil(float64(msgCount) / float64(batchSize)))

	errs, ctx := errgroup.WithContext(ctx)
	errs.SetLimit(runtime.GOMAXPROCS(0))

	for i := 0; i < batches; i++ {
		i := i
		errs.Go(func() error {
			msgs, err := s.repository.FindManyByMailing(ctx, m, batchSize, i*batchSize)
			if err != nil {
				return err
			}
			for j := 0; j < len(msgs); j++ {
				msg := msgs[j]
				err := s.mailer.SendMail([]string{msg.Email}, []byte(msg.Content))
				if err != nil {
					s.logger.WithError(err).WithField("email", msg.Email).WithField("mailing_id", m.ID).Error("failed to send email")
					// efficient way to remove item from slice
					msgs[j] = msgs[len(msgs)-1]
					msgs = msgs[:len(msgs)-1]
					j--
				}
			}
			return s.repository.DeleteManyCustomers(ctx, msgs)
		})
	}

	if err = errs.Wait(); err != nil {
		return err
	}
	return err
}
