package handler

import (
	"github.com/VladPetriv/tg_scanner/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Manager
}

func NewHandler(service *service.Manager) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/hello", h.sayHello).Methods("GET")

	return router
}
