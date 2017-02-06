package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"os"
	"time"
	"log"
	"encoding/json"
	"net/http"
)

type BFSStats struct {
	ActiveStations int `json:"betriebsbereit"`
	Average        struct{
				Value float64 `json:"mw"`
				Time  string  `json:"t"`
		             } `json:"mwavg"`
	Minimum        struct{
				StationID string  `json:"kenn"`
				Value     float64 `json:"mw"`
		             } `json:"mwmin"`
	Maximum        struct{
				StationID string  `json:"kenn"`
				Value     float64 `json:"mw"`
		             } `json:"mwmax"`
}

type BFSValues struct{
	Station struct{
			City       string `json:"ort"`
			Id         string `json:"kenn"`
		        PLZ        string `json:"plz"`
		        State         int `json:"status"`
			AreaID        int `json:"kid"`
		        Elevation     int `json:"hoehe"`
			Longitude float64 `json:"lon"`
			Latitude  float64 `json:"lat"`
			Average   float64 `json:"mw"`
		      }`json:"stamm"`
	Average1H struct{
			Time      []string  `json:"t"`
		        Average   []float64 `json:"mw"`
			Status    []int     `json:"ps"`
		        TimeRain  []string  `json:"tr"`
		        RainProp  []float64 `json:"r"`
		        Cosmic    []float64 `json:"cos"`
		        Terrestic []float64 `json:"ter"`
		  }`json:"mw1h"`
	Average24H struct{
			Time      []string  `json:"t"`
		        Average   []float64 `json:"mw"`
		        Cosmic    []float64 `json:"cos"`
		        Terrestic []float64 `json:"ter"`
		  }`json:"mw1h"`
}

func getStatistic (url, username, passwd string) (BFSStats) {
	var stat *BFSStats;
	req_url := fmt.Sprintf("%s%s", url, "stat.json")
	log.Println(req_url)

    	client := &http.Client{}
	req, err := http.NewRequest("GET", req_url, nil)
    	req.SetBasicAuth(username, passwd)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching data: %s for url: %s", err.Error(), url)
		return *stat
	}

	decoder := json.NewDecoder(resp.Body);
	if err = decoder.Decode(&stat); err != nil {
		log.Println("Error decoding feed: %s", err.Error())
		return *stat
	}
	return *stat
}

func getValues (url, username, passwd, station string) (BFSValues){
	var values *BFSValues;
	req_url := fmt.Sprintf("%s%s%s", url, station, "ct.json")
	log.Println(req_url)

    	client := &http.Client{}
	req, err := http.NewRequest("GET", req_url, nil)
    	req.SetBasicAuth(username, passwd)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching data: %s for url: %s", err.Error(), url)
		return *values
	}

	decoder := json.NewDecoder(resp.Body);
	if err = decoder.Decode(&values); err != nil {
		log.Println("Error decoding feed: %s", err.Error())
		return *values
	}

	return *values
}

func main() {
	fmt.Println("Radiation updater")

	bfsUrl     := "https://odlinfo.bfs.de/daten/json/"
	bfsUser    := os.Getenv("BFS_USER")
	bfsPasswd  := os.Getenv("BFS_PASSWD")
	bfsStation := os.Getenv("BFS_STATION")
	mqttServer := os.Getenv("MQTT_SERVER")

	if bfsUser == ""{
		fmt.Print("BFS_USER not set")
		os.Exit(1)
	}

	if bfsPasswd == ""{
		fmt.Print("BFS_PASSWD not set")
		os.Exit(1)
	}

	if bfsStation == ""{
		fmt.Print("BFS_STATION not set")
		os.Exit(1)
	}

	if mqttServer == ""{
		fmt.Print("MQTT_SERVER not set")
		os.Exit(1)
	}

	// MQTT
	opts := mqtt.NewClientOptions().AddBroker(mqttServer).SetClientID("goradiation")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
		os.Exit(2)
	}

	// Statistics
	stat := getStatistic(bfsUrl, bfsUser, bfsPasswd)

	text  := fmt.Sprintf("%d", stat.ActiveStations)
	token := c.Publish("radiation/activeStations", 0, false, text)
	token.Wait()

	text  = fmt.Sprintf("%.3f", stat.Average.Value)
	token = c.Publish("radiation/average", 0, false, text)
	token.Wait()

	text  = fmt.Sprintf("%.3f", stat.Minimum.Value)
	token = c.Publish("radiation/min", 0, false, text)
	token.Wait()

	text  = fmt.Sprintf("%s", stat.Minimum.StationID)
	token = c.Publish("radiation/minStation", 0, false, text)
	token.Wait()

	text  = fmt.Sprintf("%.3f", stat.Maximum.Value)
	token = c.Publish("radiation/max", 0, false, text)
	token.Wait()

	text  = fmt.Sprintf("%s", stat.Maximum.StationID)
	token = c.Publish("radiation/maxStation", 0, false, text)
	token.Wait()

	// LocalStation
	values := getValues(bfsUrl, bfsUser, bfsPasswd, bfsStation)
	text    = fmt.Sprintf("%.3f", values.Station.Average)
	token   = c.Publish("radiation/localAverage", 0, false, text)
	token.Wait()

	c.Disconnect(0)
}
