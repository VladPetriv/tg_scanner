package handler

import (
	"net/http"
)

func (h *Handler) sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
