package mosquitto

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/lamassuiot/device-virtual/pkg/configs"

	"github.com/go-kit/kit/log"
)

func TestConnect(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	mq := NewClient(logger)
	cfg, err := configs.NewConfig("devicetest")
	if err != nil {
		t.Fatal("Unable to load configuration")
	}

	testCases := []struct {
		name     string
		URL      string
		clientID string
		conf     *tls.Config
		retErr   bool
	}{
		{"Incorrect URL", "thisIsNotAURL", "lamassu-client", TLSConf(t, cfg.CAPath, "testdata/valid.crt", "testdata/valid.key"), true},
		{"Incorrect Client ID", "ssl://mosquitto:1883", "", TLSConf(t, cfg.CAPath, "testdata/valid.crt", "testdata/valid.key"), true},
		{"Self-signed TLS configuration", "ssl://mosquitto:1883", "lamassu-client", TLSConf(t, cfg.CAPath, "testdata/self.crt", "testdata/self.key"), true},
		{"Correct configuration values", "ssl://mosquitto:1883", "lamassu-client", TLSConf(t, cfg.CAPath, "testdata/valid.crt", "testdata/valid.key"), false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing %s", tc.name), func(t *testing.T) {
			err := mq.Connect(tc.URL, tc.clientID, tc.conf)
			if err != nil && !tc.retErr {
				t.Errorf("Client returned an unexpected error: %s", err)
			}
		})
	}

	mq.Disconnect()
}

func TestDisconnect(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	mq := NewClient(logger)
	cfg, err := configs.NewConfig("devicetest")
	if err != nil {
		t.Fatal("Unable to load configuration")
	}

	validConf := TLSConf(t, cfg.CAPath, "testdata/valid.crt", "testdata/valid.key")
	err = mq.Connect("ssl://mosquitto:1883", "lamassu-client", validConf)
	if err != nil {
		t.Fatal("Unable to connect to the broker")
	}

	mq.Disconnect()

	err = mq.SendMessage("this is a message", "lamassu-sample")
	if err == nil {
		t.Errorf("Client was expected to return an error")
	}
}

func TestSendMessage(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	mq := NewClient(logger)
	cfg, err := configs.NewConfig("devicetest")
	if err != nil {
		t.Fatal("Unable to load configuration")
	}

	validConf := TLSConf(t, cfg.CAPath, "testdata/valid.crt", "testdata/valid.key")
	err = mq.Connect("ssl://mosquitto:1883", "lamassu-client", validConf)
	if err != nil {
		t.Fatal("Unable to connect to the broker")
	}

	testCases := []struct {
		name    string
		message string
		topic   string
		retErr  bool
	}{
		{"Topic empty", "this is a message", "", false},
		{"Message empty", "", "lamassu-test", false},
		{"Correct message values", "this is a message", "lamassu-test", false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing %s", tc.name), func(t *testing.T) {
			err := mq.SendMessage(tc.message, tc.topic)
			if err != nil && !tc.retErr {
				t.Errorf("Client returned an unexpected error: %s", err)
			}
		})
	}

	mq.Disconnect()
}

func TLSConf(t *testing.T, CAPath string, certPath string, keyPath string) *tls.Config {
	t.Helper()

	caCert, err := ioutil.ReadFile(CAPath)
	if err != nil {
		t.Fatal("Unable to read CA certificate file")
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM([]byte(caCert))
	if !ok {
		t.Fatal("CA certificate file does not contain any certificate")
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		t.Fatal("Unable to load certificate and/or key files")
	}

	conf := &tls.Config{
		RootCAs:            caCertPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
	}

	return conf

}
