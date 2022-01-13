package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"code.vegaprotocol.io/shared/paths"
	_ "github.com/ArtemGontar/d-wallet/docs"
	netstore "github.com/ArtemGontar/d-wallet/network/store/v1"
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
	router       *mux.Router
	logger       *logrus.Logger
	walletStore  *walletstore.Store
	networkStore *netstore.Store
}

type ctxKey int8

func newServer() *server {
	walletStore, walletErr := walletstore.InitialiseStore("wallets1111")
	networkStore, networkErr := netstore.InitialiseStore(paths.New("network11111"))
	if walletErr != nil || networkErr != nil {
		return &server{}
	}
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		walletStore:  walletStore,
		networkStore: networkStore,
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
	//wallets
	s.router.HandleFunc("/wallets", s.getListWallets).Methods("GET")
	s.router.HandleFunc("/wallets/{id}", s.getWalletInfo).Methods("GET")
	s.router.HandleFunc("/wallets", s.createWallet).Methods("POST")
	s.router.HandleFunc("/wallets/import", s.importWallet).Methods("POST")
	s.router.HandleFunc("/wallets", s.deleteWallet).Methods("DELETE")

	//network
	s.router.HandleFunc("/networks", s.getNetworks).Methods("GET")
	s.router.HandleFunc("/networks/{name}", s.getNetworkInfo).Methods("GET")
	s.router.HandleFunc("/networks/import", s.importNetwork).Methods("POST")
	s.router.HandleFunc("/networks", s.deleteNetwork).Methods("DELETE")

	//swagger
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

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
