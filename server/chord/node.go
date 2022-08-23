package chord

import (
	"context"
	"errors"
	"math/big"
	"net"
	"os"
	"sync"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Node struct {
	*chord.Node // Real node.

	predecessor *chord.Node        // Predecessor of this node in the ring.
	predLock    sync.RWMutex       // Locks the predecessor for reading or writing.
	successors  *Queue[chord.Node] // Queue of successors of this node in the ring.
	sucLock     sync.RWMutex       // Locks the queue of successors for reading or writing.

	fingerTable FingerTable  // FingerTable of this node.
	fingerLock  sync.RWMutex // Locks the finger table for reading or writing.

	RPC    RemoteServices // Transport layer of this node.
	config *Configuration // General configurations.

	dictionary Storage      // Storage dictionary of this node.
	dictLock   sync.RWMutex // Locks the dictionary for reading or writing.

	server   *grpc.Server     // Node server.
	sock     *net.TCPListener // Node server listener socket.
	shutdown chan struct{}    // Determine if the node server is actually running.

	chord.UnimplementedChordServer
}

// NewNode crea y devuelve un nuevo nodo.
func NewNode(port string, configuration *Configuration, transport RemoteServices, storage Storage) (*Node, error) {
	// Su la configuracion es vacia devuelve error.
	if configuration == nil {
		log.Error("Error creando el nodo: la configuracion no puede ser vacia.")
		return nil, errors.New("error creando el nodo: la configuracion no puede ser vacia.")
	}

	// Crea un nuevo nodo con la ID obtenida y la misma dirección.
	innerNode := &chord.Node{ID: big.NewInt(0).Bytes(), IP: "0.0.0.0", Port: port}

	// Crea la instancia del nodo.
	node := &Node{Node: innerNode,
		predecessor: nil,
		successors:  nil,
		fingerTable: nil,
		RPC:         transport,
		config:      configuration,
		dictionary:  storage,
		server:      nil,
		shutdown:    nil}

	// Devuelve el nodo.
	return node, nil
}

// DefaultNode crea y devuelve un nuevo nodo con una configuracion por defecto.
func DefaultNode(port string) (*Node, error) {
	conf := DefaultConfig()                    // Crea una configuracion por defecto.
	transport := NewGRPCServices(conf)         // Crea un objeto RPC por defecto  para interactuar con la capa de transporte.
	dictionary := NewDiskDictionary(conf.Hash) // Crea un diccionario por defecto.

	// Devuelve el nodo creado.
	return NewNode(port, conf, transport, dictionary)
}

// GetPredecessor devuelve el nodo que se cree que es el actual predecesor.
func (node *Node) GetPredecessor(ctx context.Context, req *chord.GetPredecessorRequest) (*chord.GetPredecessorResponse, error) {
	log.Trace("Obteniendo el predecesor del nodo.")

	// Bloquea el predecesor para leerlo, se desbloquea el terminar.
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	// Se crea el mensaje respuesta que contiene al nodo predecesor
	res := &chord.GetPredecessorResponse{
		Predecessor: pred,
	}
	// Devuelve el predecesor de este nodo.
	return res, nil
}

// GetSuccessor returns the node believed to be the current successor.
func (node *Node) GetSuccessor(ctx context.Context, req *chord.EmptyRequest) (*chord.Node, error) {
	log.Trace("Getting node successor.")

	// Lock the successor to read it, and unlock it after.
	node.sucLock.RLock()
	suc := node.successors.Beg()
	node.sucLock.RUnlock()

	// Return the successor of this node.
	return suc, nil
}

// SetPredecessor sets the predecessor of this node.
func (node *Node) SetPredecessor(ctx context.Context, candidate *chord.Node) (*chord.EmptyResponse, error) {
	log.Tracef("Setting node predecessor to %s.", candidate.IP)

	// If the new predecessor is not this node, update this node predecessor.
	if !Equals(candidate.ID, node.ID) {
		// Lock the predecessor to read and write on it, and unlock it after.
		node.predLock.Lock()
		old := node.predecessor
		node.predecessor = candidate
		node.predLock.Unlock()
		// If there was an old predecessor, absorb its keys.
		go node.AbsorbPredecessorKeys(old)
	} else {
		log.Trace("Candidate predecessor is this same node. Update refused.")
	}

	return emptyResponse, nil
}

// SetSuccessor sets the successor of this node.
func (node *Node) SetSuccessor(ctx context.Context, candidate *chord.Node) (*chord.EmptyResponse, error) {
	log.Tracef("Setting node successor to %s.", candidate.IP)

	// If the new successor is not this node, update this node successor.
	if !Equals(candidate.ID, node.ID) {
		// Lock the successor to write on it, and unlock it after.
		node.sucLock.Lock()
		node.successors.PushBeg(candidate)
		node.sucLock.Unlock()
		// Update this new successor with this node keys.
		go node.UpdateSuccessorKeys()
	} else {
		log.Trace("Candidate successor is this same node. Update refused.")
	}

	return emptyResponse, nil
}

// FindSuccessor busca el nodo sucesor de esta ID.
func (node *Node) FindSuccessor(ctx context.Context, req *chord.FindSuccesorRequest) (*chord.FindSuccesorResponse, error) {
	// Busca el sucesor de esta ID.
	node_id := req.GetID()
	new_node, err := node.FindIDSuccessor(node_id)
	if err != nil {
		return nil, err
	}
	succesor := chord.Node{ID: new_node.ID, IP: new_node.IP, Port: new_node.Port}
	res := &chord.FindSuccesorResponse{Succesor: &succesor}
	return res, nil
}

// Notify notifica a este nodo que es posible que tenga un nuevo predecesor.
func (node *Node) Notify(ctx context.Context, req *chord.NotifyRequest) (*chord.NotifyResponse, error) {
	log.Trace("Comprobando la notificacion de predecesor.")

	candidate := req.GetNotify()

	// Bloquea el predecesor para leer de el, lo desbloquea al terminar.
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	/*

		Si este nodo no tiene predecesor o el candidato a predecesor esta más cerca a este
		nodo que su actual predecesor, se actualiza el predecesor con el candidato
	*/
	if Equals(pred.ID, node.ID) || Between(candidate.ID, pred.ID, node.ID) {
		log.Debugf("Predecesor actualizad al nodo en %s.", candidate.IP)

		// Bloquea el predecesor para escribir en el, lo desbloquea al terminar.
		node.predLock.Lock()
		node.predecessor = candidate
		node.predLock.Unlock()

		//Actualiza el nuevo predecesor con la correspondiente clave.
		go node.UpdatePredecessorKeys(pred)
	}

	return &chord.NotifyResponse{}, nil
}

//Comprueba si el nodo esta vivo.
func (node *Node) Check(ctx context.Context, req *chord.CheckRequest) (*chord.CheckResponse, error) {
	return &chord.CheckResponse{}, nil
}

// Get the value associated to a key.
func (node *Node) Get(ctx context.Context, req *chord.GetRequest) (*chord.GetResponse, error) {
	log.Infof("Get: key=%s.", req.Key)
	address := req.IP // Obtain the requesting address.

	// If block is needed.
	if req.Lock {
		node.dictLock.RLock()                         // Lock the dictionary to read it, and unlock it after.
		err := node.dictionary.Lock(req.Key, address) // Lock this key on storage.
		node.dictLock.RUnlock()
		if err != nil {
			log.Errorf("Error locking key: already locked.\n%s", err)
			return nil, err
		}
	} else {
		node.dictLock.RLock()                           // Lock the dictionary to read it, and unlock it after.
		err := node.dictionary.Unlock(req.Key, address) // Unlock this key on storage.
		node.dictLock.RUnlock()
		if err != nil {
			log.Errorf("Error locking key: already locked.\n%s", err)
			return nil, err
		}
	}

	keyNode := node.Node  // By default, take this node to get the value of this key from the local storage.
	node.predLock.RLock() // Lock the predecessor to read it, and unlock it after.
	pred := node.predecessor
	node.predLock.RUnlock()

	// If the key ID is not between the predecessor ID and this node ID,
	// then the requested key is not necessarily local.
	if between, err := KeyBetween(req.Key, node.config.Hash, pred.ID, node.ID); !between && err == nil {
		log.Debug("Searching for the corresponding node.")
		keyNode, err = node.LocateKey(req.Key) // Locate the node that stores the key.
		if err != nil {
			log.Error("Error getting key.")
			return &chord.GetResponse{}, errors.New("error getting key\n" + err.Error())
		}
	} else if err != nil {
		log.Error("Error getting key.")
		return &chord.GetResponse{}, errors.New("error getting key\n" + err.Error())
	}

	// If the node that stores the key is this node, directly get the associated value from this node storage.
	if Equals(keyNode.ID, node.ID) {
		log.Debug("Resolving get request locally.")

		// Lock the dictionary to read it, and unlock it after.
		node.dictLock.RLock()
		value, err := node.dictionary.GetWithLock(req.Key, address) // Get the value associated to this key from storage.
		node.dictLock.RUnlock()
		if err != nil && err != os.ErrPermission {
			log.Error("Error getting key.\n" + err.Error())
			return &chord.GetResponse{}, os.ErrNotExist
		} else if err == os.ErrPermission {
			log.Error("Error getting key: already locked.\n" + err.Error())
			return &chord.GetResponse{}, err
		}

		log.Info("Successful get.")
		// Return the value associated to this key.
		return &chord.GetResponse{Value: value}, nil
	} else {
		log.Infof("Redirecting get request to %s.", keyNode.IP)
	}
	// Otherwise, return the result of the remote call on the correspondent node.
	return node.RPC.Get(keyNode, req)
}

// Extend agrega al diccionario local un nuevo conjunto de pares<key, values> .
func (node *Node) Extend(ctx context.Context, req *chord.ExtendRequest) (*chord.ExtendResponse, error) {
	log.Debug("Agregando nuevos elementos al diccionario de almacenamiento local.")

	//Si no hay llaves que agregar regresa.
	if req.Dictionary == nil || len(req.Dictionary) == 0 {
		return &chord.ExtendResponse{}, nil
	}

	// Bloquea el diccionario para escribir en el , al terminar se desbloquea
	node.dictLock.Lock()
	err := node.dictionary.Extend(req.Dictionary) // Agrega los pares <key, value> al almacenamiento.
	node.dictLock.Unlock()
	if err != nil {
		log.Error("Error agregando los elementos al diccionario de almacenamiento.")
		return &chord.ExtendResponse{}, errors.New("error agregando los elementos al diccionario de almacenamiento\n" + err.Error())
	}
	return &chord.ExtendResponse{}, err
}

// Discard a list of keys from local storage dictionary.
func (node *Node) Discard(ctx context.Context, req *chord.DiscardRequest) (*chord.DiscardResponse, error) {
	log.Debug("Descartando llaves desde el diccionario de almacenamiento local.")

	//Bloquea el diccionario para escribir en el, se desbloquea el final
	node.dictLock.Lock()
	err := node.dictionary.Discard(req.GetKeys()) // Elimina las llaves del almacenamiento.
	node.dictLock.Unlock()
	if err != nil {
		log.Error("Error descartando las llaves del diccionario de almacenamiento.")
		return &chord.DiscardResponse{}, errors.New("error descartando las llaves del diccionario de almacenamiento\n" + err.Error())
	}
	return &chord.DiscardResponse{}, err
}
