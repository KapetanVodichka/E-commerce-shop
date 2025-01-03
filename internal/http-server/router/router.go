package router

import (
	"database/sql"
	"e-commerce-shop/internal/http-server/handlers/good"
	"e-commerce-shop/internal/http-server/handlers/user"
	"e-commerce-shop/internal/http-server/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Main page если чо!"))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", user.Register(db))
			r.Post("/login", user.Login(db))
		})

		r.Route("/good", func(r chi.Router) {
			r.Get("/list", good.GetGoodList(db))
			r.Get("/{goodID}", good.GetGoodDetail(db))

			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuth)
				r.Post("/", good.CreateGood(db))
				r.Put("/{goodID}", good.ChangeGood(db))
				r.Delete("/{goodID}", good.DeleteGood(db))
			})
		})
	})

	return r
}
