package transport

type Transport interface {
	SendRequest(message []byte) ([]byte, error)
}
