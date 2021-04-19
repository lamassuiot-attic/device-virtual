package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/lamassuiot/device-virtual/pkg/client"
	"github.com/lamassuiot/device-virtual/pkg/configs"
	"github.com/lamassuiot/device-virtual/pkg/mocks"
)

type serviceSetUp struct {
	client client.Client
	CAPath string
}

func TestPostConnect(t *testing.T) {
	stu := setup(t)
	srv := NewDeviceService(stu.CAPath, stu.client)
	ctx := context.Background()

	stu.client.(*mocks.MockClient).ConnectFn = func(URL string, clientID string, conf *tls.Config) error {
		return nil
	}

	validKey, err := ioutil.ReadFile("testdata/valid.key")
	if err != nil {
		t.Fatal("Unable to read valid key")
	}
	validCert, err := ioutil.ReadFile("testdata/valid.crt")
	if err != nil {
		t.Fatal("Unable to read valid certificate")
	}

	testCases := []struct {
		name      string
		authKey   string
		authCRT   string
		brokerURL string
		clientID  string
		ret       error
	}{
		{"Authentication key invalid", "thisIsNotAKey", string(validCert), "ssl://mosquitto:1883", "lamassu-client", ErrTLSConfLoading},
		{"Authentication certificate invalid", string(validKey), "thisIsNotACert", "ssl://mosquitto:1883", "lamassu-client", ErrTLSConfLoading},
		{"Broker URL empty", string(validKey), string(validCert), "", "lamassu-client", ErrBrokerURLEmpty},
		{"ClientID empty", string(validKey), string(validCert), "ssl://mosquitto:1883", "", ErrClientIDEmpty},
		{"Valid request", string(validKey), string(validCert), "ssl://mosquitto:1883", "lamassu-client", nil},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing %s", tc.name), func(t *testing.T) {
			err := srv.PostConnect(ctx, tc.authKey, tc.authCRT, tc.brokerURL, tc.clientID)
			if tc.ret != err {
				t.Errorf("Got result is %s; want %s", err, tc.ret)
			}
		})
	}
}

func TestPostSendMessage(t *testing.T) {
	stu := setup(t)
	srv := NewDeviceService(stu.CAPath, stu.client)
	ctx := context.Background()

	stu.client.(*mocks.MockClient).SendMessageFn = func(message string, topic string) error {
		return nil
	}

	testCases := []struct {
		name    string
		message string
		topic   string
		ret     error
	}{
		{"Topic empty", "this is a message", "", ErrTopicEmpty},
		{"Correct topic", "this is a message", "lamassu-sample", nil},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing %s", tc.name), func(t *testing.T) {
			err := srv.PostSendMessage(ctx, tc.message, tc.topic)
			if tc.ret != err {
				t.Errorf("Got result is %s; want %s", err, tc.ret)
			}
		})
	}
}

func TestPostDisconnect(t *testing.T) {
	stu := setup(t)
	srv := NewDeviceService(stu.CAPath, stu.client)
	ctx := context.Background()

	stu.client.(*mocks.MockClient).DisconnectFn = func() {}

	srv.PostDisconnect(ctx)
}

func setup(t *testing.T) *serviceSetUp {
	t.Helper()

	cfg, err := configs.NewConfig("devicetest")
	if err != nil {
		t.Fatal("Unable to get configuration variables")
	}
	client := &mocks.MockClient{}

	return &serviceSetUp{CAPath: cfg.CAPath, client: client}
}
