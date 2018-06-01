package main

import (
	"flag"
	"fmt"
	"./api"
	"./dyndns"
	"./log"
	"os"
)

func main() {
	logFile := flag.String("log", "", "Specify a logfile. If no file is provided, uses stdout.")
	flag.Parse()

	var file *os.File
	defer file.Close()

	if *logFile == "" {
		file = os.Stdout
	} else {
		var err error
		file, err = os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Could not open log for reading")
			os.Exit(1)
		}
	}

	log.Init(file)

	config := api.LoadConfig()


	if config.Domain == "" || len(config.Hostnames) == 0 {
		log.Logger.Fatalf("Empty configuration detected. Exiting.")
	}

	log.Logger.Printf("Detected configuration for %s", config.Domain)
	dyndns.Run(config)
}
