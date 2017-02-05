package main

import (
	"github.com/ansgarschmidt/go-weather/weather"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/peterjliu/gowu"
	"fmt"
	"os"
	"time"
)

func main() {

	fmt.Println("Weather underground updater")

	wuKey      := os.Getenv("WU_API_KEY")
	mqttserver := os.Getenv("MQTT_SERVER")

	if wuKey == ""{
		fmt.Print("WU_API_KEY not set")
		os.Exit(1)
	}

	if mqttserver == ""{
		fmt.Print("MQTT_SERVER not set")
		os.Exit(1)
	}

	// MQTT
	opts := mqtt.NewClientOptions().AddBroker(mqttserver).SetClientID("goweatherunderground")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
		os.Exit(2)
	}

	// Current
	gowuClient    := gowu.NewClient(wuKey)
	cond, err     := gowuClient.GetConditions("Berlin", "Germany")
	if err != nil {
		fmt.Print("Error getting current conditions")
		fmt.Println(err)
		os.Exit(2)
	}
	text := fmt.Sprintf("%.1f", cond.CurrentObservation.TempC)
	token := c.Publish("outside/temperature", 0, false, text)
	token.Wait()

	text = fmt.Sprintf("%s", cond.CurrentObservation.RelativeHumidity[:len(cond.CurrentObservation.RelativeHumidity)-1])
	token = c.Publish("outside/humidity", 0, false, text)
	token.Wait()

	text = fmt.Sprintf("%.1f", cond.CurrentObservation.DewpointC)
	token = c.Publish("outside/dewpoint", 0, false, text)
	token.Wait()

	text = fmt.Sprintf("%s", cond.CurrentObservation.FeelslikeC)
	token = c.Publish("outside/feelslike", 0, false, text)
	token.Wait()

	// Forecast
	forecast      := weather.CreateClient(wuKey).Get10DayForecast("Germany/Berlin")
	for i:=0; i<len(forecast.Forecast); i++{

		key   := fmt.Sprintf("outside/forecast/%s/temperature", forecast.Forecast[i].WeatherTime.Hour)
		text  := fmt.Sprintf("%s", forecast.Forecast[i].Temperature.Metric)
		token := c.Publish(key, 0, false, text)
		token.Wait()

		key    = fmt.Sprintf("outside/forecast/%s/humidity", forecast.Forecast[i].WeatherTime.Hour)
		text   = fmt.Sprintf("%s", forecast.Forecast[i].Humidity)
		token  = c.Publish(key, 0, false, text)
		token.Wait()

		key    = fmt.Sprintf("outside/forecast/%s/dewpoint", forecast.Forecast[i].WeatherTime.Hour)
		text = fmt.Sprintf("%s", forecast.Forecast[i].Dewpoint.Metric)
		token = c.Publish(key, 0, false, text)
		token.Wait()

		key    = fmt.Sprintf("outside/forecast/%s/feelslike", forecast.Forecast[i].WeatherTime.Hour)
		text = fmt.Sprintf("%s", forecast.Forecast[i].Feelslike.Metric)
		token = c.Publish(key, 0, false, text)
		token.Wait()
	}

	c.Disconnect(0)

	os.Exit(0)
}
