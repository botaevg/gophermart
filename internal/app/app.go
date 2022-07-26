package app

import (
	"github.com/botaevg/gophermart/internal/apperror"
	"github.com/botaevg/gophermart/internal/config"
	"github.com/botaevg/gophermart/internal/handlers"
	"github.com/botaevg/gophermart/internal/repositories"
	"github.com/botaevg/gophermart/internal/service"
	"github.com/botaevg/gophermart/pkg/postgre"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func StartApp() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Print("config error")
		return
	}
	dbpool, err := postgre.NewClient(cfg.DataBaseDsn)
	if err != nil {
		log.Print("DB connect error")
		return
	}

	storage := repositories.NewDB(dbpool)

	auth := service.NewAuth(storage, cfg.SecretKey)
	log.Print(cfg.SecretKey)
	gophermart := service.NewGophermart(storage)
	r := chi.NewRouter()

	h := handlers.NewHandler(cfg, auth, gophermart)

	authcookie := apperror.NewAuthMiddleware(cfg.SecretKey)
	r.Use(authcookie.AuthCookie)

	r.Use(middleware.Logger)

	r.Post("/api/user/register", h.RegisterUser)
	r.Post("/api/user/login", h.Login)

	r.Post("/api/user/orders", h.LoadOrder)
	r.Get("/api/user/orders", h.GetListOrders)

	//r.Get("/api/user/balance", h.Login)
	//r.Post("/api/user/balance/withdraw", h.Login)
	//r.Get("/api/user/balance/withdrawals", h.Login)

	log.Fatal(http.ListenAndServe(cfg.RunAddress, r))
}