package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/tkanos/gonfig"
	"github.com/tushar9989/url-short/read-server/internal/controllers"
	"github.com/tushar9989/url-short/read-server/internal/database"
)

type Configuration struct {
	Port       int
	DbServers  []string
	DbKeySpace string
}

func main() {
	config := Configuration{}
	err := gonfig.GetConf("../../config.json", &config)
	if err != nil {
		log.Fatal("Could not load configuration")
	}

	db, dbErr := database.NewCassandra(config.DbServers, config.DbKeySpace)

	if dbErr != nil {
		log.Fatal(dbErr)
	}

	http.HandleFunc("/s/", controllers.ReadSlug(db))
	http.HandleFunc("/stats", controllers.ReadUserStats(db))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
