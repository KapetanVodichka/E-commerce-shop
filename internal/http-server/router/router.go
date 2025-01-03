package router

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Main page если чо!"))
		})
	})

	return r
}
