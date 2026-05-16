package mqttclient

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func New() (mqtt.Client, error) {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883")
	opts.SetClientID(os.Getenv("EMQX_CLIENTID")).SetPassword(os.Getenv("EMQX_CLIENT_PASSWORD"))

	opts.SetKeepAlive(60 * time.Second)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return c, nil
}
