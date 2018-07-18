package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tkanos/gonfig"
	"github.com/tushar9989/url-short/write-server/internal/controllers"
	"github.com/tushar9989/url-short/write-server/internal/database"
	"github.com/tushar9989/url-short/write-server/internal/pkg/counters"
)

type Configuration struct {
	Port              int
	ServerID          string
	DbServers         []string
	DbKeySpace        string
	DbPersistInterval int
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

	data, dbErr := db.LoadServerMeta(config.ServerID)

	if dbErr != nil {
		log.Fatal(dbErr)
	}

	counter, err := counters.NewBigInt(data.Start, data.End, data.Current)

	if err != nil {
		log.Fatal(err)
	}

	go startPeriodicDbUpdate(db, config.ServerID, counter, time.Minute*time.Duration(config.DbPersistInterval))

	http.HandleFunc("/shorten", controllers.Write(counter, db))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}

func startPeriodicDbUpdate(db database.Database, serverId string, counter counters.BigInt, interval time.Duration) {
	time.Sleep(interval)
	for range time.Tick(interval) {
		db.UpdateServerCount(serverId, counter.Value())
	}
}
