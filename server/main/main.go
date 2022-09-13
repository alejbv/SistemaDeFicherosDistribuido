package main

import (
	"github.com/alejbv/SistemaDeFicherosDistribuido/server/logging"
	"github.com/alejbv/SistemaDeFicherosDistribuido/server/services"
	log "github.com/sirupsen/logrus"
)

func main() {

	logging.SettingLogger(log.InfoLevel, "server")
	rsaPrivateKeyPath := "pv.pem"
	rsaPublicteKeyPath := "pub.pem"
	network := "tcp"
	services.StartServer(network, rsaPrivateKeyPath, rsaPublicteKeyPath)

}
