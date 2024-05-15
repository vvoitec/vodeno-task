package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/vvoitec/vodeno-task/internal/config"
	"github.com/vvoitec/vodeno-task/internal/handler"
	customer_jobqueue "github.com/vvoitec/vodeno-task/internal/mailing/jobqueue"
	"github.com/vvoitec/vodeno-task/internal/mailing/smtp"
	"github.com/vvoitec/vodeno-task/internal/mailing/sqlc"
	"github.com/vvoitec/vodeno-task/pkg/jobqueue"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	logger.Out = os.Stdout
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()
	logger.Info("starting server, hello")

	appConfig := config.MustProvide()
	worker := jobqueue.NewWorker(ctx, log.New(logger.Writer(), "", 0))

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", appConfig.WebApiPort),
		Handler:        registerHandlers(ctx, appConfig, worker),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       log.New(logger.Writer(), "", 0),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.WithError(err).Fatal("failed to serve")
		}
	}()

	go func() {
		worker.Run()
	}()

	select {
	case <-ctx.Done():
		logger.Info("gracefully exiting, bye")
	}
}

func registerHandlers(ctx context.Context, appConfig config.App, worker *jobqueue.Worker) http.Handler {
	dbPool, err := pgxpool.New(ctx, appConfig.DatabaseURL)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize database connection")
	}
	customerSQLCRepository := sqlc.NewCustomerSQLCRepository(dbPool)
	customerSQLCLocker := sqlc.NewSQLCMailingLocker(dbPool)
	customerSMTPSender := smtp.NewSender(customerSQLCRepository, customerSQLCLocker, smtp.Config{SMTPURL: appConfig.SMTPURL}, logger, nil)
	customerController := handler.NewCustomerController(customerSQLCRepository, customerSMTPSender, logger, worker.GetQueue())

	customer_jobqueue.RegisterHandler(customerSMTPSender, logger, worker)

	return handler.Setup(customerController, logger)
}
