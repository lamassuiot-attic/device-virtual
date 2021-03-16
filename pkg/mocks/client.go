package mocks

import "crypto/tls"

type MockClient struct {
	ConnectFn      func(URL string, clientID string, conf *tls.Config) error
	ConnectInvoked bool

	DisconnectFn      func()
	DisconnectInvoked bool

	SendMessageFn      func(message string, topic string) error
	SendMessageInvoked bool
}

func (mc *MockClient) Connect(URL string, clientID string, conf *tls.Config) error {
	mc.ConnectInvoked = true
	return mc.ConnectFn(URL, clientID, conf)
}

func (mc *MockClient) Disconnect() {
	mc.DisconnectInvoked = true
	mc.DisconnectFn()
}

func (mc *MockClient) SendMessage(message string, topic string) error {
	mc.SendMessageInvoked = true
	return mc.SendMessageFn(message, topic)
}
