package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/derekpitt/weather_site/postdata"
	"github.com/derekpitt/weather_site/server/data"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// flags
var (
	port      = flag.Int("port", 8080, "port to listen on")
	serverKey = flag.String("key", "powpow", "the key to sign each request with")
	username  = flag.String("user", "root", "mysql username")
	password  = flag.String("password", "dev", "mysql password")
	database  = flag.String("database", "weather", "mysql database name")
)

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

func getLatestJsonString() (string, error) {
	latestData, err := data.GetLatestSample()
	jsonData, err := json.Marshal(latestData)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func latest(w http.ResponseWriter, r *http.Request) {
	latest, err := getLatestJsonString()

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", latest)
}

var indexTemplate, _ = template.ParseFiles("views/index.html")

func index(w http.ResponseWriter, r *http.Request) {
	latestData, _ := data.GetLatestSample()
	indexTemplate.Execute(w, latestData)
}

func main() {
	flag.Parse()

	// open database
	err := data.OpenDatabase(*username, *password, *database)
	if err != nil {
		log.Fatal("Error opening database")
	}
	defer data.CloseDatabase()

	// static dir
	cwd, _ := os.Getwd()
	staticDir := cwd + "/static"
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(staticDir))))

	// data handler
	http.HandleFunc("/data", postData)
	http.HandleFunc("/latest", latest)
	http.HandleFunc("/", index)

	log.Println("Listening on port", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
