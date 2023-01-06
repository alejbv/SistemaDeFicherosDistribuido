package client

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/alejbv/SistemaDeFicherosDistribuido/chord"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NetDiscover trata de encontrar si existe en la red algun otro anillo chord.
func NetDiscover(ip net.IP) (string, error) {
	// La dirección de broadcast.
	ip[2] = 255
	ip[3] = 255
	broadcast := ip.String() + ":8888"

	//Resuelve la dirección a la que se va a hacer broadcast.
	outAddr, err := net.ResolveUDPAddr("udp", broadcast)
	if err != nil {
		log.Errorf("Error resolviendo la direccion de broadcast %s.", broadcast)
		return "", err
	}

	// Tratando de escuchar al puerto en uso.
	conn, err := net.ListenUDP("udp", outAddr)
	if err != nil {
		log.Errorf("Error al escuchar a la dirección %s.", broadcast)
		return "", err
	}
	defer conn.Close()

	log.Infof("Resuelta direccion UPD broadcast:%s.", broadcast)

	//Se crea un objeto para recibir la respuesta del otro servidor
	listAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:8838")
	if err != nil {
		log.Error("Error resolviendo la direccion para escuchar la respuesta 0.0.0.0:8838.\n")
		return "", err
	}

	// Build listining connections
	listen, err := net.ListenUDP("udp", listAddr)
	if err != nil {
		log.Errorf("Error al escuchar a la direccion de respuesta ")
		return "", err
	}
	defer listen.Close()

	// Enviando el mensaje
	message := []byte("Chord?")
	_, err = conn.WriteToUDP(message, outAddr)
	if err != nil {
		log.Errorf("Error enviando el mensaje de broadcast a la dirección%s.", broadcast)
		return "", err
	}

	log.Info("Terminado el mensaje broadcast.")
	top := time.Now().Add(10 * time.Second)

	log.Info("Esperando por respuesta.")

	for top.After(time.Now()) {
		// Creando el buffer para almacenar el mensaje.
		buff := make([]byte, 15000)

		// Estableciendo el tiempo de espera para mensajes entrantes.
		err = listen.SetReadDeadline(time.Now().Add(2 * time.Second))
		if err != nil {
			log.Error("Error establecientdo el tiempo de espera para mensajes entrantes.")
			return "", err
		}

		// Esperando por un mensaje.
		n, address, err := listen.ReadFromUDP(buff)
		if err != nil {
			log.Errorf("Error leyendo el mensaje entrante.\n%s", err.Error())
			continue
		}

		log.Infof("Mensaje de respuesta entrante. %s enviado a: %s", address, buff[:n])

		if string(buff[:n]) == "Yo soy chord" {
			return strings.Split(address.String(), ":")[0], nil
		}
	}

	log.Info("Tiempo de espera por respuesta pasado.")
	return "", nil
}

// GetOutboundIP obtiene la  IP de este nodo en la red.
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {

		log.Fatal(err)
		return nil
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("No se pudo cerrar la conexion al servidor")

		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func StartClient() (chord.ChordClient, *grpc.ClientConn, error) {

	// Obtiene la IP del cliente en la red
	ip := GetOutboundIP()
	// Se obtiene la direccion del cliente
	addr := ip.String()
	log.Infof("Dirección IP del cliente en la red %s.", addr)
	addres, err := NetDiscover(ip)
	addres += ":50050"

	if err != nil {
		log.Errorln("Hubo problemas encontrando un servidor al que conectarse")
		return nil, nil, fmt.Errorf("Hubo problemas encontrando un servidor al que conectarse" + err.Error())
	}

	conn, err := grpc.Dial(addres, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("No se pudo realizar la conexion al servidor con IP %s\n %v", addres, err)
		return nil, nil, err
	}

	//defer conn.Close()
	client := chord.NewChordClient(conn)
	fmt.Printf("El cliente es: %v", client)
	return client, conn, nil

}
