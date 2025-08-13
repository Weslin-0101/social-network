package middlewares

import (
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
		next(w,r)
	}
}