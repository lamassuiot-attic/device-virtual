package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"sync"

	"github.com/lamassuiot/device-virtual/pkg/client"

	"github.com/pkg/errors"
)

const (
	certificatePEMBlockType = "CERTIFICATE"
)

type Service interface {
	Health(ctx context.Context) bool
	PostSendMessage(ctx context.Context, message string, topic string) error
	PostConnect(ctx context.Context, authKey string, authCRT string, brokerURL string, clientID string) error
	PostDisconnect(ctx context.Context)
}

type deviceService struct {
	mtx    sync.RWMutex
	client client.Client
	CAPath string
}

func NewDeviceService(CAPath string, client client.Client) Service {
	return &deviceService{CAPath: CAPath, client: client}
}

var (
	ErrSendMessage    = errors.New("error sending message")
	ErrDeviceAuth     = errors.New("error authenticating device")
	ErrCACertLoading  = errors.New("unable to read CA certificate")
	ErrTLSConfLoading = errors.New("unable to read client TLS configuration")
	ErrBrokerURLEmpty = errors.New("invalid empty broker URL")
	ErrClientIDEmpty  = errors.New("invalid empty client ID")
	ErrTopicEmpty     = errors.New("invalid empty topic")
)

func (s *deviceService) Health(ctx context.Context) bool {
	return true
}

func (s *deviceService) PostSendMessage(ctx context.Context, message string, topic string) error {
	if topic == "" {
		return ErrTopicEmpty
	}

	err := s.client.SendMessage(message, topic)
	if err != nil {
		return ErrSendMessage
	}
	return nil
}

func (s *deviceService) PostConnect(ctx context.Context, authKey string, authCRT string, brokerURL string, clientID string) error {
	if brokerURL == "" {
		return ErrBrokerURLEmpty
	}

	if clientID == "" {
		return ErrClientIDEmpty
	}

	conf, err := newTLSConfig(s.CAPath, authKey, authCRT)
	if err != nil {
		return err
	}

	err = s.client.Connect(brokerURL, clientID, conf)
	if err != nil {
		return ErrDeviceAuth
	}
	return nil
}

func (s *deviceService) PostDisconnect(ctx context.Context) {
	s.client.Disconnect()
}

func newTLSConfig(CAPath string, authKey string, authCRT string) (*tls.Config, error) {
	caCertPool, err := createCACertPool(CAPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair([]byte(authCRT), []byte(authKey))
	if err != nil {
		return nil, ErrTLSConfLoading
	}

	conf := &tls.Config{
		RootCAs:            caCertPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
	}

	return conf, nil
}

func createCACertPool(CAPath string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(CAPath)
	if err != nil {
		return nil, ErrCACertLoading
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM([]byte(caCert))
	if !ok {
		return nil, ErrCACertLoading
	}
	return caCertPool, nil
}
