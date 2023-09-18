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

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5005
// @BasePath /

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
