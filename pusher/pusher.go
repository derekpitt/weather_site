package main

import (
	"encoding/json"
	"flag"
	"github.com/derekpitt/weather_site/postdata"
	"github.com/derekpitt/weather_station"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	stationDevice = flag.String("dev", "/dev/tty.SLAB_USBtoUART", "the tty where the weather station console is hooked up")
	server        = flag.String("server", "http://localhost:8080/data", "the server post location to post the data to")
	serverKey     = flag.String("key", "powpow", "the key to sign each request with")
)

// TODO: figure out how to make a mutex around this...
var packetRetries = make(map[time.Time]postdata.PostData)

// prevent us from posting more than one thing at a time (if we want to call this from go routines)
var serverPostMutex = &sync.Mutex{}

func sendPacket(data postdata.PostData) error {
	jsonData, _ := json.Marshal(data)

	log.Printf("sending json %s", jsonData)

	serverPostMutex.Lock()
	_, err := http.Post(*server, "application/json", strings.NewReader(string(jsonData)))
	serverPostMutex.Unlock()

	if err != nil {
		// append to retries
		packetRetries[time.Now()] = data

		log.Printf("couldn't send...  saving for later: %d\n", len(packetRetries))
	}

	return err
}

func sendRetries() {
	if len(packetRetries) == 0 {
		return
	}

	for key, data := range packetRetries {
		log.Printf("Retrying from %s...", key)
		go sendPacket(data)

		delete(packetRetries, key)
	}
}

func main() {
	flag.Parse()

	ws, err := weather_station.New(*stationDevice)
	if err != nil {
		panic(err)
	}

	sampleTickChan := time.Tick(5 * time.Second)
	retryTickChan := time.Tick(5 * time.Second)

	log.Printf("Starting tick")

	// retries
	go func() {
		for {
			<-retryTickChan
			sendRetries()
		}
	}()

	// sample fetching
	for {
		<-sampleTickChan

		sample, err := ws.GetSample()

		if err != nil {
			log.Printf("some sort of error")
		} else {
			data := postdata.NewData(time.Now(), sample, *serverKey)
			sendPacket(data)
		}
	}

}
