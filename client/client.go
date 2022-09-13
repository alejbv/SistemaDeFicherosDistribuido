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
	ip[3] = 255
	broadcast := ip.String() + ":8830"

	// Tratando de escuchar al puerto en uso.
	pc, err := net.ListenPacket("udp4", ":8830")
	if err != nil {
		log.Errorf("Error al escuchar a la dirección %s.", broadcast)
		return "", err
	}

	//Resuelve la dirección a la que se va a hacer broadcast.
	out, err := net.ResolveUDPAddr("udp4", broadcast)
	if err != nil {
		log.Errorf("Error resolviendo la direccion de broadcast %s.", broadcast)
		return "", err
	}

	log.Info("Resuelta direccion UPD broadcast.")

	// Enviando el mensaje
	_, err = pc.WriteTo([]byte("Chord?"), out)
	if err != nil {
		log.Errorf("Error enviando el mensaje de broadcast a la dirección%s.", broadcast)
		return "", err
	}

	log.Info("Terminado el mensaje broadcast.")
	top := time.Now().Add(10 * time.Second)

	log.Info("Esperando por respuesta.")

	for top.After(time.Now()) {
		// Creando el buffer para almacenar el mensaje.
		buf := make([]byte, 1024)

		// Estableciendo el tiempo de espera para mensajes entrantes.
		err = pc.SetReadDeadline(time.Now().Add(2 * time.Second))
		if err != nil {
			log.Error("Error establecientdo el tiempo de espera para mensajes entrantes.")
			return "", err
		}

		// Esperando por un mensaje.
		n, address, err := pc.ReadFrom(buf)
		if err != nil {
			log.Errorf("Error leyendo el mensaje entrante.\n%s", err.Error())
			continue
		}

		log.Debugf("Mensaje de respuesta entrante. %s enviado a: %s", address, buf[:n])

		if string(buf[:n]) == "Yo soy chord" {
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

func StartClient() (chord.ChordClient, error) {

	// Obtiene la IP del cliente en la red
	ip := GetOutboundIP()
	// Se obtiene la direccion del cliente
	addr := ip.String()
	log.Infof("Dirección IP del cliente en la red %s.", addr)
	addres, err := NetDiscover(ip)

	if err != nil {
		log.Errorln("Hubo problemas encontrando un servidor al que conectarse")
		return nil, fmt.Errorf("Hubo problemas encontrando un servidor al que conectarse" + err.Error())
	}

	conn, err := grpc.Dial(addres, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("No se pudo realizar la conexion al servidor con IP %s\n %v", addres, err)
		return nil, err
	}

	defer conn.Close()
	client := chord.NewChordClient(conn)
	fmt.Printf("El cliente es: %v", client)
	return client, nil

}
