package chord

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// RemoteServices permite a un nodo a interactuar con otros nodos en el anillo, como un cliente a sus servicios.
type RemoteServices interface {
	Start() error // Empieza el servicio.
	Stop() error  // Para el servicio.

	// GetPredecessor returns the node believed to be the current predecessor of a remote node.
	GetPredecessor(*chord.Node) (*chord.Node, error)
	// GetSuccessor returns the node believed to be the current successor of a remote node.
	GetSuccessor(*chord.Node) (*chord.Node, error)
	// SetPredecessor sets the predecessor of a remote node.
	SetPredecessor(*chord.Node, *chord.Node) error
	// SetSuccessor sets the successor of a remote node.
	SetSuccessor(*chord.Node, *chord.Node) error
	// FindSuccessor encuentra el nodo que sucede a la ID, partiendo desde un nodo remoto.
	FindSuccessor(*chord.Node, []byte) (*chord.Node, error)
	// Notify a remote node that it possibly have a new predecessor.
	Notify(*chord.Node, *chord.Node) error
	// Check Comprueba si un nodo remoto esta vivo.
	Check(*chord.Node) error

	// Get the value associated to a key on a remote node storage.
	Get(node *chord.Node, req *chord.GetRequest) (*chord.GetResponse, error)
	// Set a <key, value> pair on a remote node storage.
	Set(node *chord.Node, req *chord.SetRequest) error
	// Delete a <key, value> pair from a remote node storage.
	Delete(node *chord.Node, req *chord.DeleteRequest) error
	// Extend the storage dictionary of a remote node with a list of <key, values> pairs.
	Extend(node *chord.Node, req *chord.ExtendRequest) error
	// Partition return all <key, values> pairs in a given interval from the storage of a remote node.
	Partition(node *chord.Node, req *chord.PartitionRequest) (*chord.PartitionResponse, error)
	// Discard all <key, values> pairs in a given interval from the storage of a remote node.
	Discard(node *chord.Node, req *chord.DiscardRequest) error
}

type GRPCServices struct {
	*Configuration // Remote service configurations.

	connections    map[string]*RemoteNode // Dictionary of <address, open connection>.
	connectionsMtx sync.RWMutex           // Locks the dictionary for reading or writing.

	shutdown chan struct{} // Determine if the service is actually running.
}

// NewGRPCServices crea un nuevo objeto tipo GRPCServices .
func NewGRPCServices(config *Configuration) *GRPCServices {
	// Crea el GRPCServices.
	services := &GRPCServices{
		Configuration: config,
		connections:   nil,
		shutdown:      nil,
	}

	// Devuelve el GRPCServices.
	return services
}

// Empieza el servicio.
func (services *GRPCServices) Start() error {
	log.Info("Empezando el servicio de la capa de transporte...\n")

	// Si los servicios de la capa de tranporte estan en funcionamiento reporta un error .
	if IsOpen(services.shutdown) {
		message := "Error empezando el servicio:los servicios de la capa de tranporte estan en funcionamiento .\n"
		log.Error(message)
		return errors.New(message)
	}

	services.shutdown = make(chan struct{}) // Reportar que los servicios estan en ejecucion.

	services.connections = make(map[string]*RemoteNode) // Crear el diccionario de <address, open connection>.
	// Empezar hilos periodicos.
	go services.CloseOldConnections() // Comprueba y cierra viejas conexiones.

	log.Info("Servicios de la capa de tranporte en funcionamiento\n")
	return nil
}

// CloseOldConnections cierra las viejas conexiones abiertas.
func (services *GRPCServices) CloseOldConnections() {
	log.Trace("Cerrando viejas conexiones.\n")

	// Si el servicio esta apagado , cierra todas las conexiones  y regresa.
	if !IsOpen(services.shutdown) {
		services.connectionsMtx.Lock() // Bloquea el diccionario para que no se escriba en el,Se desbloquea despues.
		// Para  cerrar las conexiones con los nodos en el diccionario.
		for _, remoteNode := range services.connections {
			remoteNode.CloseConnection()
		}
		services.connections = nil // Borra el diccionario de conexiones.
		services.connectionsMtx.Unlock()
		return
	}

	services.connectionsMtx.RLock() // Bloquea el diccionario para lectura en el,Se desbloquea despues.
	if services.connections == nil {
		services.connectionsMtx.Unlock()
		log.Error("Error cerrando conexiones: la tabla de conexiones esta vacia.\n")
		return
	}
	services.connectionsMtx.RUnlock()

	services.connectionsMtx.Lock()

	for addr, remoteNode := range services.connections {
		if time.Since(remoteNode.lastActive) > services.MaxIdle {
			remoteNode.CloseConnection()
			delete(services.connections, addr) // Elimina los pares <address, open connection> en el diccionario.
		}
	}
	services.connectionsMtx.Unlock()
	log.Trace("Viejas conexiones cerradas.\n")
}

/*
Cierra el servicio, reportando que los servicios de la capa de transporte estan apagados
*/
func (services *GRPCServices) Stop() error {
	log.Info("Parando los servicios de la capa de transporte...\n")

	// Si los servicios de la capa de transporte no estan en funcionamiento reporta error.
	if !IsOpen(services.shutdown) {
		message := "Error parando los servicios: estos ya estan cerrados.\n"
		log.Error(message)
		return errors.New(message)
	}

	close(services.shutdown) // Reportar que los servicios se pararon.
	log.Info("Servicios de la capa de transporte parados.\n")
	return nil
}

// Connect establece una conexion con una direccion remota.
func (services *GRPCServices) Connect(addr string) (*RemoteNode, error) {
	log.Trace("Conectando con la direccion " + addr + ".\n")

	// Comprueba si el servicio esta apagado, en caso de que si se termina y se reporta.
	if !IsOpen(services.shutdown) {
		message := "Error creando la conexion: se deben empezar los servicios de la capa de transporte primero .\n"
		log.Error(message)
		return nil, errors.New(message)
	}

	services.connectionsMtx.RLock()
	if services.connections == nil {
		services.connectionsMtx.Unlock()
		message := "Error creando la conexion: la tabla de conexion esta vacia.\n"
		log.Error(message)
		return nil, errors.New(message)
	}
	remoteNode, ok := services.connections[addr]
	services.connectionsMtx.RUnlock()

	// Comprueba si existe una conexion que este en funcionamiento, de ser asi se devuelve dicha conexion .
	if ok {
		log.Trace("Conexion exitosa.\n")
		return remoteNode, nil
	}

	conn, err := grpc.Dial(addr, services.DialOpts...) // De otra forma se crea.
	if err != nil {
		message := "Error creando la conexion.\n"
		log.Error(message)
		return nil, errors.New(message + err.Error())
	}

	client := chord.NewChordClient(conn) // Crear el client chord asociado con la conexion.
	// Se construye el correspondiente nodo remoto.

	remoteNode = &RemoteNode{client,
		addr,
		conn,
		time.Now()}

	services.connectionsMtx.Lock()
	services.connections[addr] = remoteNode
	services.connectionsMtx.Unlock()

	log.Trace("Conexion exitosa.\n")
	return remoteNode, nil
}

// GetPredecessor returns the node believed to be the current predecessor of a remote node.
func (services *GRPCServices) GetPredecessor(node *chord.Node) (*chord.Node, error) {
	if node == nil {
		return nil, errors.New("No se puede establecer conexion con un nodo nulo.\n")
	}

	// Estableciendo conexion con el nodo remoto
	remoteNode, err := services.Connect(node.IP + ":" + node.Port)
	if err != nil {
		return nil, err
	}

	// Se obtiene el contexto de esta conexion y el tiempo de espera de la request .
	ctx, cancel := context.WithTimeout(context.Background(), services.Timeout)
	defer cancel()

	//Se devuelve el resultado de la llamada remota.
	res, err := remoteNode.GetPredecessor(ctx, &chord.GetPredecessorRequest{})
	return res.GetPredecessor(), err
}

// GetSuccessor returns the node believed to be the current successor of a remote node.
func (services *GRPCServices) GetSuccessor(node *chord.Node) (*chord.Node, error) {
	if node == nil {
		return nil, errors.New("Cannot establish connection with a null node.\n")
	}

	remoteNode, err := services.Connect(node.IP + ":" + node.Port) // Establish connection with the remote node.
	if err != nil {
		return nil, err
	}

	// Obtain the context of the connection and set the timeout of the request.
	ctx, cancel := context.WithTimeout(context.Background(), services.Timeout)
	defer cancel()

	// Return the result of the remote call.
	return remoteNode.GetSuccessor(ctx, emptyRequest)
}

// SetPredecessor sets the predecessor of a remote node.
func (services *GRPCServices) SetPredecessor(node, pred *chord.Node) error {
	if node == nil {
		return errors.New("Cannot establish connection with a null node.\n")
	}

	remoteNode, err := services.Connect(node.IP + ":" + node.Port) // Establish connection with the remote node.
	if err != nil {
		return err
	}

	// Obtain the context of the connection and set the timeout of the request.
	ctx, cancel := context.WithTimeout(context.Background(), services.Timeout)
	defer cancel()

	// Return the result of the remote call.
	_, err = remoteNode.SetPredecessor(ctx, pred)
	return err
}

// SetSuccessor sets the successor of a remote node.
func (services *GRPCServices) SetSuccessor(node, suc *chord.Node) error {
	if node == nil {
		return errors.New("Cannot establish connection with a null node.\n")
	}

	remoteNode, err := services.Connect(node.IP + ":" + node.Port) // Establish connection with the remote node.
	if err != nil {
		return err
	}

	// Obtain the context of the connection and set the timeout of the request.
	ctx, cancel := context.WithTimeout(context.Background(), services.Timeout)
	defer cancel()

	// Return the result of the remote call.
	_, err = remoteNode.SetSuccessor(ctx, suc)
	return err
}

// FindSuccessor encuentra el nodo que sucede a esta ID, empezando desde el nodo remoto.
func (services *GRPCServices) FindSuccessor(node *chord.Node, id []byte) (*chord.Node, error) {
	if node == nil {
		return nil, errors.New("No se puede establecer una conexion con un nodo nulo.\n")
	}

	remoteNode, err := services.Connect(node.IP + ":" + node.Port) // Establece la conexion con el nodo remoto.
	if err != nil {
		return nil, err
	}

	// Obtiene el contexto de la conexion y establece  el tiempo de espera de la peticion.
	ctx, cancel := context.WithTimeout(context.Background(), services.Timeout)
	defer cancel()

	// Devuelve el resultado de la llamada remota.
	req := chord.FindSuccesorRequest{ID: id}
	res, err := remoteNode.FindSuccessor(ctx, &req)
	if err != nil {
		return nil, err
	}

	return res.GetSuccesor(), nil
}

// Notify notifica a un nodo remoto que es posible que tenga un nuevo predecesor.
func (services *GRPCServices) Notify(node, pred *chord.Node) error {
	if node == nil {
		return errors.New("No se puede establecer conexion con un nodo vacio.\n")
	}

	// Estableciendo conexion con un nodo remoto
	remoteNode, err := services.Connect(node.IP + ":" + node.Port)
	if err != nil {
		return err
	}

	// Obteniendo el contexto de la conexion y el tiempo de espera de la.
	ctx, cancel := context.WithTimeout(context.Background(), services.Timeout)
	defer cancel()

	// Devuelve el resultado de la llamada remota.
	req := &chord.NotifyRequest{Notify: pred}
	_, err = remoteNode.Notify(ctx, req)
	return err
}

// Comprueba si un nodo remoto esta vivo.
func (services *GRPCServices) Check(node *chord.Node) error {
	if node == nil {
		return errors.New("No se puede establecer una conexion con un nodo nulo.\n")
	}

	// Establce conexion con el nodo remoto
	remoteNode, err := services.Connect(node.IP + ":" + node.Port)
	if err != nil {
		return err
	}

	// Obtiene el contexto de la conexion y establece el tiempo de duracion de esta.
	ctx, cancel := context.WithTimeout(context.Background(), services.Timeout)
	defer cancel()

	// Devuelve el resultado de la llamada remota.
	_, err = remoteNode.Check(ctx, &chord.CheckRequest{})
	return err
}

// Extend agrega una lista de pares  <key, values> en el diccionario de almacenamiento de un  odo remoto.
func (services *GRPCServices) Extend(node *chord.Node, req *chord.ExtendRequest) error {
	if node == nil {
		return errors.New("No se puede establecer conexion con un nodo nulo.\n")
	}

	// Estableciendo conexion con un nodo remoto
	remoteNode, err := services.Connect(node.IP + ":" + node.Port)
	if err != nil {
		return err
	}

	// Obteniendo el contexto y el tiempo de duracion de la request.
	ctx, cancel := context.WithTimeout(context.Background(), services.Timeout)
	defer cancel()

	// Devuelve el resultado de la llamada remota.
	_, err = remoteNode.Extend(ctx, req)
	return err
}

// Discard elimina todos los pares <key, values>  en el almacenamiento interno de un nodo remoto.
func (services *GRPCServices) Discard(node *chord.Node, req *chord.DiscardRequest) error {
	if node == nil {
		return errors.New("No se puede establecer una conexion con un nodo vacio.\n")
	}

	// Estableciendo conexion con un nodo remoto
	remoteNode, err := services.Connect(node.IP + ":" + node.Port)
	if err != nil {
		return err
	}

	// Se obtiene el contexto de la conexion y el tiempo de espera de la request.
	ctx, cancel := context.WithTimeout(context.Background(), services.Timeout)
	defer cancel()

	// Devuelve el resultado de la llamada remota.
	_, err = remoteNode.Discard(ctx, req)
	return err
}
