package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/ArtemGontar/d-wallet/docs"
	"github.com/ArtemGontar/d-wallet/wallet"
	walletstore "github.com/ArtemGontar/d-wallet/wallet/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

func Start(config *Config) error {
	srv := newServer()
	return http.ListenAndServe(config.BindAddr, srv)
}

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  *walletstore.Store
}

type ctxKey int8

func newServer() *server {
	store, err := walletstore.InitialiseStore("wallets1111")
	if err != nil {
		return &server{}
	}
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/hello", s.handleHello).Methods("GET")
	s.router.HandleFunc("/wallets/{id}", s.getWalletInfo).Methods("GET")
	s.router.HandleFunc("/wallets", s.createWallet).Methods("POST")
	s.router.HandleFunc("/wallets/import", s.importWallet).Methods("POST")
	s.router.HandleFunc("/wallets", s.deleteWallet).Methods("DELETE")
	s.router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		next.ServeHTTP(rw, r)
		logger.Infof("completed in %v", time.Now().Sub(start))

	})
}

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

// GetWallet godoc
// @Summary Get wallet info
// @Description Method for get wallet info
// @Tags wallets
// @Accept  json
// @Produce  json
// @Param  id         path   string     true   "Wallet ID"
// @Param  passphrase query  string  false  "Pass phrase"
// @Success 200 {object} wallet.GetWalletInfoResponse
// @Router /wallets/{id} [get]
func (s *server) getWalletInfo(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletId := vars["id"]
	passphrase := r.URL.Query().Get("passphrase")
	s.logger.Infof("test %s %s", walletId, passphrase)
	resp, err := wallet.GetWalletInfo(s.store, &wallet.GetWalletInfoRequest{
		Wallet:     walletId,
		Passphrase: passphrase,
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
	resp, err := wallet.CreateWallet(s.store, createWalletRequest)
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
	resp, err := wallet.ImportWallet(s.store, importWalletRequest)
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
	err := wallet.DeleteWallet(s.store, deleteWalletRequest)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusOK, nil)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
