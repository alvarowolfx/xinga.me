package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// Xingamento 0800
type Xingamento struct {
	Value string `json:"xingamento"`
}

func getRandomFromFile(file string) string {
	rand.Seed(time.Now().UnixNano())
	all, _ := ioutil.ReadFile(file)
	list := bytes.Split(all, []byte("\n"))
	return string(list[rand.Intn(len(list))])
}

func newRandomXingamento() Xingamento {
	return Xingamento{
		Value: fmt.Sprintf("%s %s",
			getRandomFromFile("data/subjects.txt"),
			getRandomFromFile("data/predicates.txt"),
		),
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	xingamento, _ := json.Marshal(newRandomXingamento())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(xingamento)
}

type slackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func slackHandler(w http.ResponseWriter, r *http.Request) {
	xingamento := newRandomXingamento()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(slackResponse{
		ResponseType: "in_channel",
		Text:         xingamento.Value,
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("templates/index.html")
	tpl, _ := template.New("index").Parse(string(content))

	xingamento := newRandomXingamento()
	data := struct {
		Value string
	}{
		Value: xingamento.Value,
	}

	tpl.Execute(w, data)

}

func init() {
	//staticHandler := http.FileServer(http.Dir("static"))

	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/slack", slackHandler)
	//http.Handle("/", staticHandler)
	http.HandleFunc("/", indexHandler)
}
