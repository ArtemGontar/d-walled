package apiserver

import (
	"net/http"
)

// GetNetworks godoc
// @Summary Get networks
// @Description Method for get networks
// @Tags networks
// @Accept  json
// @Produce  json
// @Success 200 {object} network.ListNetworksResponse
// @Router /networks/ [get]
func (s *server) getNetworks(rw http.ResponseWriter, r *http.Request) {

	s.respond(rw, r, http.StatusCreated, nil)
}

// GetNetwork godoc
// @Summary Get network info
// @Description Method for get network info
// @Tags networks
// @Accept  json
// @Produce  json
// @Param  name   path   string   true   "Network name"
// @Success 200 {object} network.DescribeNetworkResponse
// @Router /networks/{name} [get]
func (s *server) getNetworkInfo(rw http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//name := vars["name"]

	s.respond(rw, r, http.StatusCreated, nil)
}

// ImportNetwork godoc
// @Summary Import network info
// @Description Method for import network
// @Tags networks
// @Accept  json
// @Produce  json
// @Param data body network.ImportNetworkFromSourceRequest true "The input for import network"
// @Success 200 {object} network.ImportNetworkFromSourceResponse
// @Router /networks/import [post]
func (s *server) importNetwork(rw http.ResponseWriter, r *http.Request) {
	s.respond(rw, r, http.StatusCreated, nil)
}

// DeleteNetwork godoc
// @Summary Delete network info
// @Description Method for Delete network
// @Tags networks
// @Accept  json
// @Produce  json
// @Param data body network.DeleteNetworkRequest true "The input for delete network"
// @Success 200 {object} nil
// @Router /networks [delete]
func (s *server) deleteNetwork(rw http.ResponseWriter, r *http.Request) {
	s.respond(rw, r, http.StatusOK, nil)
}
