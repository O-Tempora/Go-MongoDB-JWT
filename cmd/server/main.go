package main

import (
	"flag"
	"gomongojwt/internal/server"
	"gomongojwt/internal/util"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	configPath string
	resetKeys  string
)

func init() {
	flag.StringVar(&configPath, "config", "configs/default.yaml", "server and db configuration")
	flag.StringVar(&resetKeys, "resetKeys", "y", "reset rs512 keys or not")
}

func main() {
	flag.Parse()
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal("Selected config file doesn't exist")
	}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Unable to read config file")
	}
	config := server.NewConfig()
	if err = yaml.Unmarshal(data, config); err != nil {
		log.Fatal("Failed to parse config")
	}
	if resetKeys == "Y" || resetKeys == "y" {
		util.SeedRS512Keys()
	}
	if err = server.StartServer(config); err != nil {
		log.Fatal(err)
	}
}
