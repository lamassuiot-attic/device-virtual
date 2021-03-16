package client

import "crypto/tls"

type Client interface {
	Connect(URL string, clientID string, conf *tls.Config) error
	Disconnect()
	SendMessage(message string, topic string) error
}
