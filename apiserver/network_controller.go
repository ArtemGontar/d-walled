package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemGontar/d-wallet/network"
	"github.com/gorilla/mux"
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
	resp, err := network.ListNetworks(s.networkStore)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusCreated, &resp)
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
	vars := mux.Vars(r)
	name := vars["name"]
	resp, err := network.DescribeNetwork(s.networkStore, &network.DescribeNetworkRequest{
		Name: name,
	})
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusCreated, &resp)
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
	importNetworkRequest := &network.ImportNetworkFromSourceRequest{}
	if err := json.NewDecoder(r.Body).Decode(importNetworkRequest); err != nil {
		s.error(rw, r, http.StatusBadRequest, err)
		return
	}
	resp, err := network.ImportNetworkFromSource(s.networkStore,
		network.NewReaders(), importNetworkRequest)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusCreated, &resp)
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
	deleteNetworkRequest := &network.DeleteNetworkRequest{}
	if err := json.NewDecoder(r.Body).Decode(deleteNetworkRequest); err != nil {
		s.error(rw, r, http.StatusBadRequest, err)
		return
	}
	err := network.DeleteNetwork(s.networkStore, deleteNetworkRequest)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusOK, nil)
}
