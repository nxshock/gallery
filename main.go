package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	previewCache = NewPreviewCache("cache.zkv")
	config       *Config
	semaphore    chan struct{}
)

func init() {
	log.SetFlags(0)
	rand.Seed(time.Now().Unix())

	configPath := defaultConfigPath
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	var err error

	config, err = loadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Chdir(config.WorkingDirectory)
	if err != nil {
		log.Fatalln(err)
	}

	semaphore = make(chan struct{}, config.ProcessCount)
}

func main() {
	go func() {
		http.HandleFunc("/preview/", previewHandler)
		http.HandleFunc("/view/", viewHandler)
		http.HandleFunc("/", rootHandler)
		log.Fatalln(http.ListenAndServe(":8080", nil))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	err := previewCache.Save()
	if err != nil {
		log.Fatalln(err)
	}
}
