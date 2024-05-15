package jobqueue

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vvoitec/vodeno-task/internal/mailing"
	"github.com/vvoitec/vodeno-task/pkg/jobqueue"
)

func RegisterHandler(sender mailing.Sender, logger *logrus.Logger, worker *jobqueue.Worker) {
	worker.AddHandler(jobqueue.JobLabel(mailing.SendMailingJobLabel), func(ctx context.Context, payload interface{}) error {
		p := payload.(mailing.Mailing)
		if err := sender.SendTo(ctx, p); err != nil {
			logger.WithError(err).Error("failed to send customer emails")
			return err
		}
		return nil
	})
}
