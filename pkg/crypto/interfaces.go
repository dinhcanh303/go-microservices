package crypto

type Crypto interface {
	RandomBytesToString() (string, error)
}
