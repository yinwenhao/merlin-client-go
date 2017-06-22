package socket

type SocketClient interface {
	Connect()
	Send(interface{}) error
	Close() error
}
