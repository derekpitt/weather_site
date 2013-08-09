package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/derekpitt/weather_site/postdata"
	"github.com/derekpitt/weather_site/server/data"
	"github.com/derekpitt/weather_site/server/gziper"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// flags
var (
	port      = flag.Int("port", 8080, "port to listen on")
	serverKey = flag.String("key", "powpow", "the key to sign each request with")
	username  = flag.String("user", "root", "mysql username")
	password  = flag.String("password", "dev", "mysql password")
	database  = flag.String("database", "weather", "mysql database name")
)

// cache stuff

type mainData struct {
	Latest data.SampleFormat
	Trends struct {
		OutsideTemerature []data.Trend
		OutsideHumidity   []data.Trend
		Barometer         []data.Trend
	}
  HighLows struct {
		OutsideTemerature []data.HighLow
		OutsideHumidity   []data.HighLow
		Barometer         []data.HighLow
  }
}

var (
	cacheMutex = &sync.RWMutex{}
	cacheData  = mainData{}
)

func fillCache() {
	latestData, err := data.GetLatestSample()

	outsideTempTrend, err := data.Get3HourTrend("OutsideTemerature", latestData.Time)
	outsideHumTrend, err := data.Get3HourTrend("OutsideHumidity", latestData.Time)
	barTrend, err := data.Get3HourTrend("Barometer", latestData.Time)

  outsideTempHighLow, err := data.Get7DayHighLow("OutsideTemerature", latestData.Time)
  outsideHumHighLow, err := data.Get7DayHighLow("OutsideHumidity", latestData.Time)
  barHighLow, err := data.Get7DayHighLow("Barometer", latestData.Time)

	if err != nil {
		return
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cacheData.Latest = latestData
	cacheData.Trends.OutsideTemerature = outsideTempTrend
	cacheData.Trends.OutsideHumidity = outsideHumTrend
	cacheData.Trends.Barometer = barTrend

  cacheData.HighLows.OutsideTemerature = outsideTempHighLow
  cacheData.HighLows.OutsideHumidity = outsideHumHighLow
  cacheData.HighLows.Barometer = barHighLow
}

func postData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	postData := &postdata.PostData{}

	err := json.Unmarshal(body, &postData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if postdata.VerifyData(*postData, *serverKey) == false {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("bad sig!!")
		return
	}
	data.WriteSample(*postData)
	log.Println("Good Times!!")
}

func getData() mainData {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return cacheData
}

func latest(w http.ResponseWriter, r *http.Request) {
	data := getData()
	dataJSON, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(dataJSON))
}

var indexTemplate, _ = template.ParseFiles("views/index.html")

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	indexTemplate.Execute(w, getData())
}

func poll(minutes time.Duration) {
	for {
		fillCache()
		time.Sleep(minutes * time.Minute)
	}
}

func main() {
	flag.Parse()

	// open database
	err := data.OpenDatabase(*username, *password, *database)
	if err != nil {
		log.Fatal("Error opening database")
	}
	defer data.CloseDatabase()

	go poll(1)

	// static dir
	cwd, _ := os.Getwd()
	staticDir := cwd + "/static"
	staticHandler := http.StripPrefix("/static", http.FileServer(http.Dir(staticDir)))
	http.HandleFunc("/static/", gziper.MakeGzipHandler(func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	}))

	// data handler
	http.HandleFunc("/data", postData)
	http.HandleFunc("/latest", gziper.MakeGzipHandler(latest))
	http.HandleFunc("/", gziper.MakeGzipHandler(index))

	log.Println("Listening on port", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
