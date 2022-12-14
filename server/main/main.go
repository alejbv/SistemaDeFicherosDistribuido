package main

import (
	"github.com/alejbv/SistemaDeFicherosDistribuido/server/services"
)

func main() {

	rsaPrivateKeyPath := "pv.pem"
	rsaPublicteKeyPath := "pub.pem"
	network := "tcp"
	services.StartServer(network, rsaPrivateKeyPath, rsaPublicteKeyPath)
}
