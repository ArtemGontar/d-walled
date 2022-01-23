package apiserver

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"github.com/ArtemGontar/d-wallet/wallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
)

// GetTransactions by block godoc
// @Summary Get transactions
// @Description Method for get transactions
// @Tags transactions
// @Accept  json
// @Produce  json
// @Success 200 {int} int
// @Router /transactions/{blockNumber} [get]
func (s *server) getTransactions(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockNumberString := vars["blockNumber"]
	n := new(big.Int)
	blockNumber, _ := n.SetString(blockNumberString, 10)
	block, err := s.ethclient.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}
	count, err := s.ethclient.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusCreated, count)
}

// CreateTransaction by block godoc
// @Summary Create Transaction
// @Description Method for create transaction
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param data body wallet.CreateTransactionRequest true "The input for create transaction"
// @Success 200 {object} types.Transaction
// @Router /transactions/create [post]
func (s *server) createTransaction(rw http.ResponseWriter, r *http.Request) {
	createTransactionRequest := &wallet.CreateTransactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(createTransactionRequest); err != nil {
		s.error(rw, r, http.StatusBadRequest, err)
		return
	}

	amount := big.NewInt(int64(createTransactionRequest.Amount))
	toAddress := common.HexToAddress(createTransactionRequest.ToAddressHex)

	nonce, err := s.ethclient.PendingNonceAt(context.Background(), toAddress)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}
	gasPrice, err := s.ethclient.SuggestGasPrice(context.Background())
	tx := types.NewTransaction(nonce, toAddress, amount, 21000, gasPrice, nil)
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	s.respond(rw, r, http.StatusCreated, &tx)
}

// SignTransaction by block godoc
// @Summary Sign Transaction
// @Description Method for sign transaction
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param data body wallet.SignTransactionRequest true "The input for sign transaction"
// @Success 200 {object} types.Transaction
// @Router /transactions/sign [post]
func (s *server) signTransaction(rw http.ResponseWriter, r *http.Request) {
	signTransactionRequest := &wallet.SignTransactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(signTransactionRequest); err != nil {
		s.error(rw, r, http.StatusBadRequest, err)
		return
	}

	tx, _, err := s.ethclient.TransactionByHash(context.Background(), common.HexToHash(signTransactionRequest.TransactionHex))
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}
	chainId, err := s.ethclient.NetworkID(context.Background())
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}

	privateKey, err := crypto.HexToECDSA(signTransactionRequest.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	s.respond(rw, r, http.StatusCreated, &tx)
}

// SendTransaction by block godoc
// @Summary Send Transaction
// @Description Method for send transaction
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param data body wallet.SendTransactionRequest true "The input for send transaction"
// @Success 200 {object} types.Transaction
// @Router /transactions/send [post]
func (s *server) sendTransaction(rw http.ResponseWriter, r *http.Request) {
	sendTransactionRequest := &wallet.SendTransactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(sendTransactionRequest); err != nil {
		s.error(rw, r, http.StatusBadRequest, err)
		return
	}

	tx, _, err := s.ethclient.TransactionByHash(context.Background(), common.HexToHash(sendTransactionRequest.TransactionHex))
	if err != nil {
		s.error(rw, r, http.StatusInternalServerError, err)
	}
	s.ethclient.SendTransaction(context.Background(), tx)
	s.respond(rw, r, http.StatusCreated, &tx)
}
