package mosquitto

import (
	"crypto/tls"
	"device-virtual/pkg/client"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type mosquitto struct {
	client MQTT.Client
	logger log.Logger
}

func NewClient(logger log.Logger) client.Client {
	return &mosquitto{logger: logger}
}

func (m *mosquitto) Connect(URL string, clientID string, conf *tls.Config) error {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(URL)
	opts.SetClientID(clientID).SetTLSConfig(conf)

	m.client = MQTT.NewClient(opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		err := token.Error()
		level.Error(m.logger).Log("err", err, "msg", "Could not connect with MQTT broker in URL "+URL)
		return err
	}
	level.Info(m.logger).Log("msg", "Client connected with MQTT broker in URL "+URL)
	return nil
}

func (m *mosquitto) Disconnect() {
	m.client.Disconnect(250)
}

func (m *mosquitto) SendMessage(message string, topic string) error {
	if token := m.client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		err := token.Error()
		level.Error(m.logger).Log("err", err, "msg", "Could not send message: "+message+" to MQTT broker in topic: "+topic)
		return err
	}
	level.Info(m.logger).Log("msg", "Message: "+message+" succesfully sent to MQTT broker in topic: "+topic)

	return nil
}
