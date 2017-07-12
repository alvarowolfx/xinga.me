package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	staticHandler := http.FileServer(http.Dir("static"))

	http.Handle("/", staticHandler)
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/slack/", slackHandler)

	log.Printf("iniciando o aplicativo %s\n", Get().Value)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func GetRandomFromFile(file string) string {
	rand.Seed(time.Now().UnixNano())
	all, _ := ioutil.ReadFile(file)
	list := bytes.Split(all, []byte("\n"))
	return string(list[rand.Intn(len(list))])
}

func Get() Xingamento {
	return Xingamento{
		Value: fmt.Sprintf("%s %s",
			GetRandomFromFile("data/subjects.txt"),
			GetRandomFromFile("data/predicates.txt"),
		),
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(Get())
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func slackHandler(w http.ResponseWriter, r *http.Request) {
	message := Get()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{
    "response_type": "in_channel",
    "attachments": [{
          "text": "%s",
          "fallback": "%s"
	  }
	]}`, message.Value, message.Value)))
}

type Xingamento struct {
	Value string `json:"xingamento"`
}
