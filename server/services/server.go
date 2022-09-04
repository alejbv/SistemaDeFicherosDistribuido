package services

import (
	"fmt"
	"log"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord"
)

var (
	node       *chord.Node
	rsaPrivate string
	rsaPublic  string
)

func StartServer(network string, rsaPrivateKeyPath string, rsaPublicteKeyPath string) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("El servicio de Sistema de Ficheros a empezado")

	var err error
	rsaPrivate = rsaPrivateKeyPath
	rsaPublic = rsaPublicteKeyPath

	//NewServerListening("0.0.0.0:50051")
	//fmt.Println("Ending the service")
	// Completar esto
	node, err = chord.DefaultNode("50050")
	if err != nil {
		log.Fatalf("No se pudo crear el nodo")
	}

	err = node.Start()

	if err != nil {
		log.Fatalf("No se pudo iniciar el nodo")
	}

}
