package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

const (
	certificatePEMBlockType = "CERTIFICATE"
)

type Service interface {
	PostSendMessage(ctx context.Context, message string) error
	PostConnect(ctx context.Context, authKey string, authCRT string, brokerURL string) error
	PostDisconnect(ctx context.Context)
}

type deviceService struct {
	mtx    sync.RWMutex
	client MQTT.Client
	CAPath string
}

func NewDeviceService(CAPath string) Service {
	return &deviceService{CAPath: CAPath}
}

var (
	ErrBadRequest = errors.New("Bad Request")
)

func (s *deviceService) PostSendMessage(ctx context.Context, message string) error {
	s.client.Publish("lamassu-sample", 0, false, message)
	return nil
}

func (s *deviceService) PostConnect(ctx context.Context, authKey string, authCRT string, brokerURL string) error {
	conf, err := s.newTLSConfig(authKey, authCRT)
	if err != nil {
		return errors.New("Unable to load TLS configuration properties")
	}
	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID("lamassu-client").SetTLSConfig(conf)

	s.client = MQTT.NewClient(opts)
	if token := s.client.Connect(); token.Wait() && token.Error() != nil {
		return errors.New("Unable to connect to the MQTT broker, check your TLS configuration")
	}
	s.client.Subscribe("lamassu-sample", 0, nil)
	return nil
}

func (s *deviceService) PostDisconnect(ctx context.Context) {
	s.client.Disconnect(250)
}

func (s *deviceService) newTLSConfig(authKey string, authCRT string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(s.CAPath)
	if err != nil {
		return nil, ErrBadRequest
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM([]byte(caCert))
	if !ok {
		return nil, ErrBadRequest
	}
	cert, err := tls.X509KeyPair([]byte(authCRT), []byte(authKey))
	if err != nil {
		return nil, ErrBadRequest
	}

	conf := &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: caCertPool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
	// Create tls.Config with desired tls properties

	return conf, nil
}
