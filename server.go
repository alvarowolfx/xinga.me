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

func getRandomFromFile(file string) string {
	rand.Seed(time.Now().UnixNano())
	all, _ := ioutil.ReadFile(file)
	list := bytes.Split(all, []byte("\n"))
	return string(list[rand.Intn(len(list))])
}

func Get() Xingamento {
	return Xingamento{
		Value: fmt.Sprintf("%s %s",
			getRandomFromFile("data/subjects.txt"),
			getRandomFromFile("data/predicates.txt"),
		),
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	xingamento, _ := json.Marshal(Get())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(xingamento)
}

func slackHandler(w http.ResponseWriter, r *http.Request) {
	xingamento := Get()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(
		fmt.Sprintf(`{"response_type": "in_channel","text": "%s"}`, xingamento.Value),
	))
}

type Xingamento struct {
	Value string `json:"xingamento"`
}
