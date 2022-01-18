package wallet

type Wallet interface {
	privateKey() string
	publicKey() string
	address() string
}
