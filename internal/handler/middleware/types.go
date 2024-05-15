package middleware

import "net/http"

type Func func(next http.Handler) http.Handler
