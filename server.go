package backend

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var (
	fontSize  float64 = 42
	imageFont font.Face
)

var imageProviders = []string{
	"https://placeimg.com/%d/%d/nature",
	"http://lorempixel.com/%d/%d/nature/",
}

// Xingamento 0800
type Xingamento struct {
	Value string `json:"xingamento"`
}

type slackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
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
func NewRandomXingamentoImage(client *http.Client, text string) io.Reader {
	const width = 1024 / 2
	const height = 768 / 2
	image := getPlaceholder(client, width, height)

	context := gg.NewContextForImage(image)
	context.SetFontFace(imageFont)

	context.SetRGB(0, 0, 0)
	strokeSize := 3
	for dy := -strokeSize; dy <= strokeSize; dy++ {
		for dx := -strokeSize; dx <= strokeSize; dx++ {
			if dx*dx+dy*dy >= strokeSize*strokeSize {
				continue
			}

			x := width/2 + float64(dx)
			y := height/2 + float64(dy)
			context.DrawStringWrapped(text, x, y, 0.5, 0.5, width/2, 1.5, gg.AlignCenter)
		}
	}
	context.SetRGB(1, 1, 1)
	context.DrawStringWrapped(text, width/2, height/2, 0.5, 0.5, width/2, 1.5, gg.AlignCenter)

	buff := new(bytes.Buffer)
	imgOut := context.Image()
	jpeg.Encode(buff, imgOut, &jpeg.Options{Quality: jpeg.DefaultQuality})

	return buff
}

func getPlaceholder(client *http.Client, width, height int) image.Image {
	return getPlaceholders(client, imageProviders[rand.Intn(len(imageProviders))], width, height)
}

func getPlaceholders(client *http.Client, url string, width, height int) image.Image {
	placeholder := fmt.Sprintf(url, width, height)
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

func getRandomFromFile(file string) string {
	rand.Seed(time.Now().UnixNano())
	all, _ := ioutil.ReadFile(file)
	list := bytes.Split(all, []byte("\n"))
	return string(list[rand.Intn(len(list))])
}

func loadFontFace() {
	// TODO: improve better loading fontface
	var err error
	imageFont, err = gg.LoadFontFace("fonts/Impact.ttf", fontSize)
	if err != nil {
		panic(err)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	xingamento, _ := json.Marshal(NewRandomXingamento())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(xingamento)
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

	hash := strings.TrimPrefix(r.URL.Path, "/image/")
	text, err := base64.StdEncoding.DecodeString(hash)
	if err != nil || len(text) == 0 {
		http.Redirect(w, r, "/image", 301)
		return
	}

	imageReader := NewRandomXingamentoImage(client, string(text))
	w.Header().Set("Content-Type", "image/jpeg")

	io.Copy(w, imageReader)
}

func randomImageHandler(w http.ResponseWriter, r *http.Request) {
	xingamento := NewRandomXingamento()
	text := strings.ToUpper(xingamento.Value)

	hash := base64.StdEncoding.EncodeToString([]byte(text))

	permURL := fmt.Sprintf("/image/%s", hash)
	http.Redirect(w, r, permURL, 301)
}

func init() {

	loadFontFace()

	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/slack", slackHandler)
	http.HandleFunc("/image", randomImageHandler)
	http.HandleFunc("/image/", imageHandler)
	http.HandleFunc("/", indexHandler)
}
