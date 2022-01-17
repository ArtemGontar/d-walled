package wallet

type Wallet interface {
	Name() string
	SetName(newName string)
	Id() string
}
