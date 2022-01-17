package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemGontar/d-wallet/wallet"
	"github.com/gorilla/mux"
)

// // GetListWallets godoc
// // @Summary Get wallets list
// // @Description Method for get wallets list
// // @Tags wallets
// // @Accept  json
// // @Produce  json
// // @Success 200 {object} wallet.ListWalletsResponse
// // @Router /wallets/ [get]
// func (s *server) getListWallets(rw http.ResponseWriter, r *http.Request) {
// 	resp, err := wallet.ListWallets(s.walletStore)
// 	if err != nil {
// 		s.error(rw, r, http.StatusInternalServerError, err)
// 	}

// 	s.respond(rw, r, http.StatusCreated, &resp)
// }

// GetWallet godoc
// @Summary Get wallet info
// @Description Method for get wallet info
// @Tags wallets
// @Accept  json
// @Produce  json
// @Param  address         path   string     true   "Wallet Address"
// @Param  passphrase query  string  false  "Pass phrase"
// @Success 200 {object} wallet.GetWalletInfoResponse
// @Router /wallets/{address} [get]
func (s *server) getWalletInfo(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	resp, err := wallet.GetWalletInfo(s.walletStore, &wallet.GetWalletInfoRequest{
		Address: address,
	})
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusCreated, &resp)
}

// CreateWallet godoc
// @Summary Create wallet info
// @Description Method for create wallet
// @Tags wallets
// @Accept  json
// @Produce  json
// @Param data body wallet.CreateWalletRequest true "The input for create wallet"
// @Success 200 {object} wallet.CreateWalletResponse
// @Router /wallets [post]
func (s *server) createWallet(rw http.ResponseWriter, r *http.Request) {
	createWalletRequest := &wallet.CreateWalletRequest{}
	if err := json.NewDecoder(r.Body).Decode(createWalletRequest); err != nil {
		s.error(rw, r, http.StatusBadRequest, err)
		return
	}
	resp, err := wallet.CreateWallet(s.walletStore, createWalletRequest)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusCreated, &resp)
}

// ImportWallet godoc
// @Summary Import wallet info
// @Description Method for import wallet
// @Tags wallets
// @Accept  json
// @Produce  json
// @Param data body wallet.ImportWalletRequest true "The input for import wallet"
// @Success 200 {object} wallet.ImportWalletResponse
// @Router /wallets/import [post]
func (s *server) importWallet(rw http.ResponseWriter, r *http.Request) {
	importWalletRequest := &wallet.ImportWalletRequest{}
	if err := json.NewDecoder(r.Body).Decode(importWalletRequest); err != nil {
		s.error(rw, r, http.StatusBadRequest, err)
		return
	}
	resp, err := wallet.ImportWallet(s.walletStore, importWalletRequest)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusCreated, &resp)
}

// DeleteWallet godoc
// @Summary Delete wallet info
// @Description Method for Delete wallet
// @Tags wallets
// @Accept  json
// @Produce  json
// @Param data body wallet.DeleteWalletRequest true "The input for delete wallet"
// @Success 200 {object} nil
// @Router /wallets [delete]
func (s *server) deleteWallet(rw http.ResponseWriter, r *http.Request) {
	deleteWalletRequest := &wallet.DeleteWalletRequest{}
	if err := json.NewDecoder(r.Body).Decode(deleteWalletRequest); err != nil {
		s.error(rw, r, http.StatusBadRequest, err)
		return
	}
	err := wallet.DeleteWallet(s.walletStore, deleteWalletRequest)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusOK, nil)
}
