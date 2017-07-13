package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"html/template"
	"image"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path"
	"runtime"
	"time"
)

// Xingamento 0800
type Xingamento struct {
	Value string `json:"xingamento"`
}

func getRandomFromFile(file string) string {
	rand.Seed(time.Now().UnixNano())
	all, _ := ioutil.ReadFile(path.Join(path.Dir(packageDirectory), file))
	list := bytes.Split(all, []byte("\n"))
	return string(list[rand.Intn(len(list))])
}

// NewRandomXingamento 0800
func NewRandomXingamento() Xingamento {
	return Xingamento{
		Value: fmt.Sprintf("%s %s",
			getRandomFromFile("data/subjects.txt"),
			getRandomFromFile("data/predicates.txt"),
		),
	}
}

// NewRandomXingamentoImage 0800
func NewRandomXingamentoImage(client *http.Client) io.Reader {
	context := gg.NewContextForImage(getPlaceholder(client))
	const width = 1024
	const height = 768

	// TODO: improve better loading fontface
	if err := context.LoadFontFace("/Library/Fonts/Arial.ttf", 96); err != nil {
		panic(err)
	}

	context.SetRGB(0, 0, 0)
	xingamento := NewRandomXingamento().Value
	strokeSize := 6
	for dy := -strokeSize; dy <= strokeSize; dy++ {
		for dx := -strokeSize; dx <= strokeSize; dx++ {
			if dx*dx+dy*dy >= strokeSize*strokeSize {
				continue
			}
			x := width/2 + float64(dx)
			y := height/2 + float64(dy)
			context.DrawStringWrapped(xingamento, x, y, 0.5, 0.5, width/2, 1.5, gg.AlignCenter)
		}
	}
	context.SetRGB(1, 1, 1)
	context.DrawStringWrapped(xingamento, width/2, height/2, 0.5, 0.5, width/2, 1.5, gg.AlignCenter)

	buff := new(bytes.Buffer)
	context.EncodePNG(buff)

	reader := bytes.NewReader(buff.Bytes())
	return reader
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	xingamento, _ := json.Marshal(NewRandomXingamento())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(xingamento)
}

type slackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func slackHandler(w http.ResponseWriter, r *http.Request) {
	xingamento := NewRandomXingamento()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(slackResponse{
		ResponseType: "in_channel",
		Text:         xingamento.Value,
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("templates/index.html")
	tpl, _ := template.New("index").Parse(string(content))

	xingamento := NewRandomXingamento()
	data := struct {
		Value string
	}{
		Value: xingamento.Value,
	}
	tpl.Execute(w, data)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	xingamento := NewRandomXingamentoImage(client)
	w.Header().Set("Content-Type", "image/png")

	response, _ := ioutil.ReadAll(xingamento)
	w.Write(response)
}

var packageDirectory string

func init() {
	_, packageDirectory, _, _ = runtime.Caller(0)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/slack", slackHandler)
	http.HandleFunc("/image", imageHandler)
	http.HandleFunc("/", indexHandler)
}

func getPlaceholder(client *http.Client) image.Image {
	placeholder := "https://placeimg.com/1024/768/nature"
	requester := http.Get
	if client != nil {
		requester = client.Get
	}

	response, err := requester(placeholder)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var wallpaper image.Image
	wallpaper, _, err = image.Decode(response.Body)

	return wallpaper
}
