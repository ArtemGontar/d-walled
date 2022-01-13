package main

import (
	"flag"
	"log"

	"github.com/ArtemGontar/d-wallet/apiserver"
	"github.com/zannen/toml"
)

// @title Crypto Wallet API
// @version 1.0
// @description This is a crypto wallet api
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email info@info.info
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
