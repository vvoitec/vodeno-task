package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/vvoitec/vodeno-task/internal/mailing"
	"github.com/vvoitec/vodeno-task/pkg/jobqueue"
	"net/http"
	"time"
)

type (
	CustomerController struct {
		customerRepository mailing.Repository
		sender             mailing.Sender
		logger             *logrus.Logger
		queue              *jobqueue.Queue
	}

	customerCreateRequest struct {
		Email         string `validate:"required,email"`
		Title         string `validate:"required"`
		Content       string `validate:"required"`
		MailingID     uint   `json:"mailing_id" validate:"required"`
		InsertionTime string `json:"insertion_time" validate:"required,RFC3339Date"`
	}

	customerSendEmailRequest struct {
		MailingID uint `json:"mailing_id" validate:"required"`
	}
)

func NewCustomerController(customerRepository mailing.Repository, sender mailing.Sender, logger *logrus.Logger, queue *jobqueue.Queue) *CustomerController {
	return &CustomerController{
		customerRepository: customerRepository,
		sender:             sender,
		logger:             logger,
		queue:              queue,
	}
}

func (c *CustomerController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto customerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			http.Error(w, "invalid json string", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(&dto); err != nil {
			errMsg := ""
			for _, err := range err.(validator.ValidationErrors) {
				errMsg += err.Error() + ", "
			}
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		if err := c.customerRepository.SaveCustomer(r.Context(), mailing.Customer{
			Email:         dto.Email,
			Title:         dto.Title,
			Content:       dto.Content,
			Mailing:       mailing.Mailing{ID: dto.MailingID},
			InsertionTime: func() time.Time { t, _ := time.Parse(time.RFC3339, dto.InsertionTime); return t }(),
		}); err != nil {
			if errors.Is(err, mailing.ErrMailingLocked) {
				http.Error(w, "", http.StatusLocked)
				return
			}
			c.logger.WithError(err).Error("failed to create customer")
			http.Error(w, "failed to create customer", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *CustomerController) SendEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto customerSendEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			http.Error(w, "invalid json string", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(&dto); err != nil {
			errMsg := ""
			for _, err := range err.(validator.ValidationErrors) {
				errMsg += err.Error() + ", "
			}
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		c.queue.Enqueue(jobqueue.JobLabel(mailing.SendMailingJobLabel), mailing.Mailing{ID: dto.MailingID})

		w.WriteHeader(http.StatusAccepted)
	}
}
