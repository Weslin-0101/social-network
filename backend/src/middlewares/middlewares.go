package middlewares

import (
	"backend/src/authentication"
	"backend/src/exceptions"
	"log"
	"net/http"
)


func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Host)
		next(w,r)
	}
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if err := authentication.ValidateToken(r); err != nil {
			exceptions.HandleError(w, r, http.StatusUnauthorized, exceptions.ErrUnauthorized)
			return
		}
		next(w,r)
	}
}