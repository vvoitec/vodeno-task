package middleware

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RecoveryFuncProvider(logger *logrus.Logger) Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					logger.WithField("panic", fmt.Sprintf("%+v\n", err)).Errorf("panic!")
					http.Error(w, "internal server error :'(", http.StatusInternalServerError)
				}

			}()
			next.ServeHTTP(w, r)
		})
	}
}
