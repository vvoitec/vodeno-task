package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/vvoitec/vodeno-task/internal/handler/middleware"
	"net/http"
	"time"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("RFC3339Date", IsRFC3339Date)
}

func Setup(customerController *CustomerController, logger *logrus.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/messages", middleware.RecoveryFuncProvider(logger)(middleware.HTTPMethodValidatorProvider(http.MethodPost)(customerController.Create())))
	mux.Handle("/api/messages/send", middleware.RecoveryFuncProvider(logger)(middleware.HTTPMethodValidatorProvider(http.MethodPost)(customerController.SendEmail())))

	return mux
}

func IsRFC3339Date(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	if err != nil {
		return false
	}
	return true
}
