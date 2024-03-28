package controller

import (
	"github.com/gorilla/mux"
	"zatrasz75/go_test/configs"
	"zatrasz75/go_test/internal/redis"
	"zatrasz75/go_test/internal/repository"
	"zatrasz75/go_test/pkg/logger"
	"zatrasz75/go_test/pkg/nats"
)

// NewRouter -.
func NewRouter(cfg *configs.Config, l logger.LoggersInterface, repo *repository.Store, rd *redis.Store, n *nats.Nats) *mux.Router {
	r := mux.NewRouter()
	newEndpoint(r, cfg, l, repo, rd, n)
	return r
}
