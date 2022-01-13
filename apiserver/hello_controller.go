package apiserver

import "net/http"

// SayHello godoc
// @Summary Method to say hello
// @Description Method to say hello
// @Tags hello
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /hello [get]
func (s *server) handleHello(rw http.ResponseWriter, r *http.Request) {
	s.respond(rw, r, http.StatusCreated, "hello")
}
