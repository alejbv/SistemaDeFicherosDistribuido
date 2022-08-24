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
func (node *Node) Start() error {
	log.Info("Iniciando el servidor...")

	// Si este nodo servidor está actualemnte en ejecución se reporta un error
	if IsOpen(node.shutdown) {
		log.Error("Error iniciando el servidor: este nodo esta actualemnte en ejecución.")
		return errors.New("error iniciando el servidor: este nodo esta actualemnte en ejecución")
	}

	node.shutdown = make(chan struct{}) // Reporta que el nodo esta corriendo

	ip := GetOutboundIP()                    // Obtiene la IP de este nodo.
	address := ip.String() + ":" + node.Port // Se obtiene la direccion del nodo.
	log.Infof("Dirección del nodo en %s.", address)

	id, err := HashKey(address, node.config.Hash) // Obtiene el ID correspondiente a la direccion.
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

func (node *Node) Listen() {
	log.Info("Empezando el servicio en el socket abierto.")
	err := node.server.Serve(node.sock)
	if err != nil {
		log.Errorf("No se puede dar el ser servicio en  %s.\n%s", node.IP, err.Error())
		return
	}
}

/*
 Join es para unir este nodo a un anillo Chord previamente existente a traves de otro nodo
 previamente descubierto en el broadcast. Para unir este nodo al anillo el sucesor inmediato
 de este nodo(Por su ID) es buscado
*/
func (node *Node) Join(knownNode *chord.Node) error {
	log.Info("Uniendose al anillo chord .")
	/*
		Si el nodo es nil entonces se devuelve error. Al menos un nodo del anillo debe ser
		conocido
	*/
	if knownNode == nil {
		log.Error("Error uniendose al anillo chord: el nodo conocido no puede ser nulo.")
		return errors.New("error uniendose al anillo chord: el nodo conocido no puede ser nulo")
	}

	log.Infof("La direccion del nodo es: %s.", knownNode.IP)
	// Busca el sucesor inmediato de esta ID en el anillo
	suc, err := node.RPC.FindSuccessor(knownNode, node.ID)
	if err != nil {
		log.Error("Error uniendose al anillo chord: no se puede encontrar el sucesor de esta ID.")
		return errors.New("error uniendose al anillo chord: no se puede encontrar el sucesor de esta ID.\n" + err.Error())
	}
	// Si la ID obtenida es exactamente la ID de este nodo, entonces este nodo ya existe.
	if Equals(suc.ID, node.ID) {
		log.Error("Error uniendose al anillo chord: un nodo con esta ID ya existe.")
		return errors.New("errror uniendose al anillo chord: un nodo con esta ID ya existe.")
	}

	// Bloquea el sucesor para escribir en el, al terminar se desbloquea.
	node.sucLock.Lock()
	node.successors.PushBeg(suc) // Actualiza el sucesor de este nodo con el obtenido.
	node.sucLock.Unlock()
	log.Infof("Union exitosa. Nodo sucesor en %s.", suc.IP)
	return nil
}

func (node *Node) FindIDSuccessor(id []byte) (*chord.Node, error) {
	log.Trace("Buscando el sucesor de ID .")

	// Si la ID es nula reporta un error.
	if id == nil {
		log.Error("Error buscando el sucesor: ID no puede ser nula.")
		return nil, errors.New("error buscando el sucesor: ID no puede ser nula")
	}

	// Buscando el mas cercano, en esta finger table, con ID menor que esta ID.
	pred := node.ClosestFinger(id)

	// Si el nodo que se esta buscando es exactamente este, se devuelve su sucesor.
	if Equals(pred.ID, node.ID) {
		// Bloquea el sucesor para leer en el, lo desbloque al terminar.
		node.sucLock.RLock()
		suc := node.successors.Beg()
		node.sucLock.RUnlock()
		return suc, nil
	}

	// Si es diferente, busca el sucesor de esta ID
	// del nodo remoto obtenido.
	suc, err := node.RPC.FindSuccessor(pred, id)
	if err != nil {
		log.Errorf("Error buscando el sucesor de la ID en %s.", pred.IP)
		return nil, errors.New(
			fmt.Sprintf("error buscando el sucesor de la ID en %s.\n%s", pred.IP, err.Error()))
	}
	// Devolver el sucesor obtenido.
	log.Trace("Encontrado el sucesor de la ID.")
	return suc, nil
}

/*
LocateKey localiza el nodo correspondiente a una llave determinada.
Para eso, el obtiene el hash de la llave para obtener el correspondiente ID, entonces busca por
el sucesor más cercano a esta ID en el anillo. dado que este es el nodo al que le corresponde la llave
*/

func (node *Node) LocateKey(key string) (*chord.Node, error) {
	log.Tracef("Localizando la llave: %s.", key)

	// Obteniendo el ID relativo a esta llave
	id, err := HashKey(key, node.config.Hash)
	if err != nil {
		log.Errorf("Error localizando la llave: %s.", key)
		return nil, errors.New(fmt.Sprintf("error localizando la llave: %s.\n%s", key, err.Error()))
	}

	// Busca y obtiene el sucesor de esta ID
	suc, err := node.FindIDSuccessor(id)
	if err != nil {
		log.Errorf("Error localizando el sucesor de esta ID: %s.", key)
		return nil, errors.New(fmt.Sprintf("error localizando el sucesor de esta ID: %s.\n%s", key, err.Error()))
	}

	log.Trace("Localizacion de la llave exitosa.")
	return suc, nil
}

// ClosestFinger busca el nodo mas cercano que precede a esta ID.
func (node *Node) ClosestFinger(ID []byte) *chord.Node {
	log.Trace("Buscando el nodo mas cercano que precede a esta ID.")
	defer log.Trace("Nodo mas cercano encontrado.")

	// Itera sobre la finger table en sentido contrario y regresa el nodo mas cercano
	// tal que su ID este entre la ID del nodo actual y la ID dada.
	for i := len(node.fingerTable) - 1; i >= 0; i-- {
		node.fingerLock.RLock()
		finger := node.fingerTable[i]
		node.fingerLock.RUnlock()

		if finger != nil && Between(finger.ID, node.ID, ID) {
			return finger
		}
	}

	// Si ningun otro cumple las condiciones se devuelve este.
	return node.Node
}

// AbsorbPredecessorKeys agrega las llaves replicadas del anterior predecesor a este nodo.
func (node *Node) AbsorbPredecessorKeys(old *chord.Node) {
	// Si el anterior predecesor no es este nodo.
	if !Equals(old.ID, node.ID) {
		log.Debug("Agregando las llaves del predecesor.")

		// Bloquea el sucesor de este nodo para leer de el, despues se desbloquea.
		node.sucLock.RLock()
		suc := node.successors.Beg()
		node.sucLock.RUnlock()

		// De existir el sucesor, transfiere las llaves del anterior predecesor a el, para mantener la replicacion.
		if !Equals(suc.ID, node.ID) {
			// Bloquea el diccionario para leer de el, al terminar se desbloquea.
			node.dictLock.RLock()
			// Obtiene las llaves del predecesor.
			in, out, err := node.dictionary.Partition(old.ID, node.ID)
			node.dictLock.RUnlock()
			if err != nil {
				log.Errorf("Error obtenido las llaves del anterior predecesor.\n%s", err.Error())
				return
			}

			log.Debug("Transferiendo las llaves del predecesor al sucesor.")
			log.Debugf("Llaves a transferir: %s", Keys(out))
			log.Debugf("LLaves restantes: %s", Keys(in))

			// Trasfiere las al nodo sucesor.
			err = node.RPC.Extend(suc, &chord.ExtendRequest{Dictionary: out})
			if err != nil {
				log.Errorf("Error transfiriendo las llaves al sucesor.\n%s", err.Error())
				return
			}
			log.Debug("Transferencia de las llaves al sucesor exitosa.")
		}
	}
}

// DeletePredecessorKeys elimina las llaves replicadas del viejo sucesor en este nodo.
// Devuelve las llaves actuales, las replicadas actuales y las eliminadas.
func (node *Node) DeletePredecessorKeys(old *chord.Node) (map[string][]byte, map[string][]byte, map[string][]byte, error) {
	//Bloquea el predecesor para leer de el, se desbloquea el terminar
	node.predLock.Lock()
	pred := node.predecessor
	node.predLock.Unlock()

	//Bloquea el diccionario para leer y escribir en el, se desbloquea el terminar la funcion
	node.dictLock.Lock()
	defer node.dictLock.Unlock()

	// Obtiene las llaves a transferir
	in, out, err := node.dictionary.Partition(pred.ID, node.ID)
	if err != nil {
		log.Error("Error obteniendo las claves correspondientes al nuevo predecesor.")
		return nil, nil, nil, errors.New("rror obteniendo las claves correspondientes al nuevo predecesor\n" + err.Error())
	}

	// Si el antiguo predecesor es distinto al nodo actual, elimina las llaves del viejo predecesor en este nodo.
	if !Equals(old.ID, node.ID) {
		// Obtiene las llaves a eliminar
		_, deleted, err := node.dictionary.Partition(old.ID, node.ID)
		if err != nil {
			log.Error("Error obteniendo las llaves replicadas del antiguo predecesor en este nodo.")
			return nil, nil, nil, errors.New(
				"error obteniendo las llaves replicadas del antiguo predecesor en este nodo\n" + err.Error())
		}

		// Elimina las llaves del anterior predecesor
		err = node.dictionary.Discard(Keys(out))
		if err != nil {
			log.Error("Error eliminando las llaves del anterior predecesor en este nodo.")
			return nil, nil, nil, errors.New("error eliminando las llaves del anterior predecesor en este nodo\n" + err.Error())
		}

		log.Debug("Eliminacion exitosa de las llaves replicadas del anterior predecesor en este nodo.")
		return in, out, deleted, nil
	}
	return in, out, nil, nil
}

// UpdatePredecessorKeys actualiza el nuevo sucesor con la clave correspondiente.
func (node *Node) UpdatePredecessorKeys(old *chord.Node) {
	//  Bloquea el predecesor para leer de el, lo desbloquea al terminar.
	node.predLock.Lock()
	pred := node.predecessor
	node.predLock.Unlock()

	//  Bloquea el sucesor para escribir en el, lo desbloquea al terminar.
	node.sucLock.Lock()
	suc := node.successors.Beg()
	node.sucLock.Unlock()

	// Si el nuevo predecesor no es el actual nodo.
	if !Equals(pred.ID, node.ID) {
		// Transfiere las claves correspondientes al nuevo predecesor.
		in, out, deleted, err := node.DeletePredecessorKeys(old)
		if err != nil {
			log.Errorf("Error actualizando las claves del nuevo predecesor.\n%s", err.Error())
			return
		}

		log.Debug("Trasnfiriendo las llaves del antiguo predecesor al sucesor.")
		log.Debugf("Llaves a transferir: %s", Keys(out))
		log.Debugf("Llaves restantes: %s", Keys(in))
		log.Debugf("Llaves a eliminar: %s", Keys(deleted))

		// Construye el diccionario del nuevo predecesor, al transferir sus claves correspondientes.
		err = node.RPC.Extend(pred, &chord.ExtendRequest{Dictionary: out})
		if err != nil {
			log.Errorf("Error transfiriendo las claves al nuevo predecesor.\n%s", err.Error())
			/*
			 De existir claves eliminadas, se reinsertan en el almacenamiento de este nodo
			 para prevenir perdida de informacion
			*/
			if deleted != nil {
				//  Bloquea el diccionario para escribir en el, lo desbloquea al terminar.
				node.dictLock.Lock()
				err = node.dictionary.Extend(deleted)
				node.dictLock.Unlock()
				if err != nil {
					log.Errorf("Error reinsertando llaves eliminadas al diccionario.\n%s", err.Error())
				}
			}
			return
		}
		log.Debug("Transferencia exitosa de las llaves al nuevo predecesor.")

		/*
		 De existir el sucesor, y es diferente del nuevo predecesor, se eliminan las llaves transferidas
		 del almacenamiento del sucesor
		*/
		if !Equals(suc.ID, node.ID) && !Equals(suc.ID, pred.ID) {
			err = node.RPC.Discard(suc, &chord.DiscardRequest{Keys: Keys(out)})
			if err != nil {
				log.Errorf("Error eliminando las claves replicadas en el nodo sucesor en %s.\n%s", suc.IP, err.Error())
				/*
				 De haber claves eliminadas, se reinsertan para prevenir la perdida de informacion
				*/
				if deleted != nil {
					// Bloquea el diccionario para leer y escribir en el, al terminar se desbloquea.
					node.dictLock.Lock()
					err = node.dictionary.Extend(deleted)
					node.dictLock.Unlock()
					if err != nil {
						log.Errorf("Error reinsertando las claves en este nodo.\n%s", err.Error())
					}
				}
				return
			}
			log.Debug("Eliminacion exitosa de las llaves replicadas del antiguo predecesor en el nodo sucesor.")
		}
	}
}

// UpdateSuccessorKeys actualiza el nuevo sucesor replicando las llaves del nodo actual.
func (node *Node) UpdateSuccessorKeys() {
	//Bloquea el predecesor para leer de el, al terminar se desbloquea.
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	//Bloquea el sucesor para escribir en el, al terminar se desbloquea.
	node.sucLock.Lock()
	suc := node.successors.Beg()
	node.sucLock.Unlock()

	// Si el sucesor es distinto de este nodo
	if !Equals(suc.ID, node.ID) {
		//Bloquea el diccionario para leer de el, al terminar se desbloquea.
		node.dictLock.RLock()
		in, out, err := node.dictionary.Partition(pred.ID, node.ID) // Obtiene las llaves de este nodo.
		node.dictLock.RUnlock()
		if err != nil {
			log.Errorf("Error obteniendo las llaves de este nodo.\n%s", err.Error())
			return
		}

		log.Debug("Transferiendo las llaves de este nodo a su sucesor.")
		log.Debugf("Llaves a transferir: %s", Keys(out))
		log.Debugf("Llaves restantes: %s", Keys(in))

		//Transfiere las llaves de este nodo a su sucesor, para actualizarlo.
		err = node.RPC.Extend(suc, &chord.ExtendRequest{Dictionary: in})
		if err != nil {
			log.Errorf("Error trasfiriendo las llaves a su sucesor en %s.\n%s", suc.IP, err.Error())
			return
		}
		log.Trace("Transferencia de llaves al sucesor exitosa.")
	}
}
