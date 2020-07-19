package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"nokia_task/api"
	"nokia_task/rabbitMQ"
	"runtime"
	"time"
)

func main() {
	r := mux.NewRouter()
	rabbitMQ.RabbitMqinit()
	go NewMonitor(300)
	log.Println("server started at port : 8080")
	r.HandleFunc("/add/user", api.AddUserHandler).Methods("POST")
	r.HandleFunc("/users", api.GetUserWithPaginationHandler).Methods("GET")
	r.HandleFunc("/reload", api.ReloadHandler).Methods("GET")

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic("error starting server")
	}

}



type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs uint64

	NumGC        uint32
	NumGoroutine int
}

func NewMonitor(duration int) {
	var m Monitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		// Just encode to json and print
		b, _ := json.Marshal(m)
		log.Println(string(b))
	}
}
