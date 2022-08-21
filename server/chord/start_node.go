package chord

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

/*
Se inicializa el servidor como un nodo del chord, se inicializan los servicios de la
capa de transporte y los hilos que periodicamente estabilizan el servidor
*/
// Start the node server, by registering the server of this node as a chord server, starting
// the transport layer services and the periodically threads that stabilizes the server.
func (node *Node) Start() error {
	log.Info("Iniciando el servidor...")

	// Si este nodo servidor está actualemnte en ejecución se reporta un error
	if IsOpen(node.shutdown) {
		log.Error("Error iniciando el servidor: este nodo esta actualemnte en ejecución.")
		return errors.New("error iniciando el servidor: este nodo esta actualemnte en ejecución")
	}

	node.shutdown = make(chan struct{}) // Reporta que el nodo esta corriendo

	ip := GetOutboundIP()                    // Obtiene la IP de este nodo.
	address := ip.String() + ":" + node.Port // Get the address of this node.
	log.Infof("Dirección del nodo en %s.", address)

	id, err := HashKey(address, node.config.Hash) // Obtain the ID relative to this address.
	if err != nil {
		log.Error("Error iniciando el nodo: no se puede obtener el hash de la dirección del nodo.")
		return errors.New("error iniciando el nodo: no se puede obtener el hash de la dirección del nodo\n" + err.Error())
	}

	node.ID = id          // Se le asigna al nodo el ID correspondiente.
	node.IP = ip.String() //Se le asigna al nodo el IP correspondiente.

	// Empezando a escuchar en la dirección correspondiente.
	log.Debug("Tratando de escuchar en la dirección correspondiente.")
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("Error iniciando el servidor: no se puede escuchar en la dirección %s.", address)
		return errors.New(
			fmt.Sprintf("error iniciando el servidor: no se puede escuchar en la dirección %s\n%s", address, err.Error()))
	}
	log.Infof("Listening at %s.", address)

	node.successors = NewQueue[chord.Node](node.config.StabilizingNodes) // Se crea la cola de sucesores.
	node.successors.PushBack(node.Node)                                  // Se establece a este nodo como su propio sucesor.
	node.predecessor = node.Node                                         // Se establece este nodo como su propio predecesor.
	node.fingerTable = NewFingerTable(node.config.HashSize)              // Se crea la finger table.
	node.server = grpc.NewServer(node.config.ServerOpts...)              // Se establece  el nodo como un servidor.
	node.sock = listener.(*net.TCPListener)                              // Se almacena el socket.
	err = node.dictionary.Clear()                                        // Se limpia el diccionario
	if err != nil {
		log.Errorf("Error iniciando el servidor: el diccionario no es valido.")
		return errors.New("errror iniciando el servidor: el diccionario no es valido\n" + err.Error())
	}

	chord.RegisterChordServer(node.server, node) // Se registra a este servidor como un servidor chord.
	log.Debug("Registrado nodo como servidor chord.")

	err = node.RPC.Start() // Se empiezan los servicios RPC (capa de transporte).
	if err != nil {
		log.Error("Error iniciando el servidor: no se puede arrancar la capa de transporte.")
		return errors.New("error iniciando el servidor: no se puede arrancar la capa de transporte\n" + err.Error())
	}

	// Empezando el servicio en el socket abierto
	go node.Listen()

	discovered, err := node.NetDiscover(ip) // Descubre la red chord, de existir
	if err != nil {
		log.Error("Error iniciando el servidor:  no se pudo descubrir una red para conectarse.")
		return errors.New("error iniciando el servidor:  no se pudo descubrir una red para conectarse\n" + err.Error())
	}

	// Si existe una red de chord existente.
	if discovered != "" {
		err = node.Join(&chord.Node{IP: discovered, Port: node.Port}) // Se une a la red descubierta.
		if err != nil {
			log.Error("Error uniendose a la red de servidores de chord.")
			return errors.New("error uniendose a la red de servidores de chord.\n" + err.Error())
		}
	} else {
		// En otro caso, Se crea.
		log.Info("Creando el anillo del chord.")
	}

	// Start periodically threads.
	go node.PeriodicallyCheckPredecessor()
	go node.PeriodicallyCheckSuccessor()
	go node.PeriodicallyStabilize()
	go node.PeriodicallyFixSuccessor()
	go node.PeriodicallyFixFinger()
	go node.PeriodicallyFixStorage()
	go node.BroadListen()

	log.Info("Servidor Iniciado.")
	return nil
}

// NetDiscover trata de encontrar si existe en la red algun otro anillo chord.
func (node *Node) NetDiscover(ip net.IP) (string, error) {
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
			log.Error("Error establecientdo eltiempo de espera para mensajes entrantes.")
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
