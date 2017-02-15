package gateways

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/jforman/carbon-golang"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

type tomlConfig struct {
	Title  string
	MQTT   mqttServerConfig   `toml:"mqtt"`
	Carbon carbonServerConfig `toml:"carbon"`
}

type mqttServerConfig struct {
	Host string
	Port int
}

type carbonServerConfig struct {
	Host    string
	Port    int
	Pattern []string
}

var config         tomlConfig

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	for _, b := range config.Carbon.Pattern {
		if b == msg.Topic() {
			key := strings.Replace(msg.Topic(), "/", ".", -1)
			value, err := strconv.ParseFloat(fmt.Sprintf("%s", msg.Payload()), 64)
			if err == nil {
				carbonReceiver, err := carbon.NewCarbon(config.Carbon.Host, config.Carbon.Port, false, false)
				if err != nil {
					panic(err)
					os.Exit(2)
				}
				carbonReceiver.SendMetric(carbon.Metric{Name: key, Value: value, Timestamp: time.Now().Unix()})
			}
		}
	}
}

func Carbon() {
	fmt.Println("Carbon Gateway")

	usr, e := user.Current()
	if e != nil {
		log.Fatal(e)
		os.Exit(2)
	}

	if _, err := toml.DecodeFile(usr.HomeDir+"/.gosomatic/config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.MQTT.Host, config.MQTT.Port))
	opts.SetClientID("carbon")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetDefaultPublishHandler(f)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
		os.Exit(2)
	}

	if token := c.Subscribe("#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

}
