package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/tkanos/gonfig"
	"github.com/tushar9989/url-short/auth-server/internal/pkg/cache"
)

type Page struct {
	ClientID string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", &Page{ClientID: clientId})
}

var authCache *cache.TTLMap = cache.New(60 * 5)

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recover() != nil {
			w.WriteHeader(http.StatusOK)
		}
	}()
	cookie, err := r.Cookie("X-Auth-Token")

	if err == nil {
		result, ok := authCache.Get(cookie.Value)

		if !ok {
			resp, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + cookie.Value)
			if err == nil && resp.StatusCode == http.StatusOK {
				decoder := json.NewDecoder(resp.Body)
				err := decoder.Decode(&result)
				if err == nil && result["aud"] == clientId {
					authCache.Put(cookie.Value, result)
					ok = true
				}
			}
		}

		if ok {
			w.Header().Set("X-User-ID", result["sub"])
			w.Header().Set("X-User-Email", result["email"])
		}
	}
	w.WriteHeader(http.StatusOK)
}

var templates = template.Must(template.ParseFiles("static/index.html"))
var clientId = ""

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type Configuration struct {
	Port     int
	ClientID string
}

func main() {
	config := Configuration{}
	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		log.Fatal("Could not load configuration")
	}
	clientId = config.ClientID
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/verify", verifyHandler)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
