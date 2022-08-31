package chord

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Node struct {
	*chord.Node // Nodo real.

	predecessor *chord.Node        // Predecesor de este nodo en el anillo.
	predLock    sync.RWMutex       // Bloquea el predecesor para lectura o escritura.
	successors  *Queue[chord.Node] // Cola de sucesores de este nodo en el anillo.
	sucLock     sync.RWMutex       // Bloquea la cola de sucesores para lectura o escritura.

	fingerTable FingerTable  // FingerTable de este nodo.
	fingerLock  sync.RWMutex // Bloquea la FingerTable para lectura o escritura.

	RPC    RemoteServices // Capa de transporte para este nodo(implementa la parte del cliente del chord).
	config *Configuration // Configuraciones generales.

	// Los metodos con el dictionario se deben modificar
	dictionary Storage      // Diccionario de almacenamiento de este nodo.
	dictLock   sync.RWMutex // Bloquea el diccionario para lectura o escritura.

	server   *grpc.Server     // Nodo servidor.
	sock     *net.TCPListener // Socket para escuchar conexiones del nodo servidor.
	shutdown chan struct{}    // Determina si el nodo esta actualemente en ejecucion.

	chord.UnimplementedChordServer
}

// NewNode crea y devuelve un nuevo nodo.
func NewNode(port string, configuration *Configuration, transport RemoteServices, storage Storage) (*Node, error) {
	// Su la configuracion es vacia devuelve error.
	if configuration == nil {
		log.Error("Error creando el nodo: la configuracion no puede ser vacia.")
		return nil, errors.New("error creando el nodo: la configuracion no puede ser vacia")
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
	conf := DefaultConfig() // Crea una configuracion por defecto.
	transport := NewGRPCServices(conf)
	path := "server/data/"                           // Crea un objeto RPC por defecto  para interactuar con la capa de transporte.
	dictionary := NewDiskDictionary(conf.Hash, path) // Crea un diccionario por defecto.

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

// GetSuccessor regresa el nodo que se cree que es el sucesor actual del nodo.
func (node *Node) GetSuccessor(ctx context.Context, req *chord.GetSuccessorRequest) (*chord.GetSuccessorResponse, error) {
	log.Trace("Obteniendo el sucesor del nodo.")

	// Bloquea el sucesor para leer de el, se desbloquea el terminar
	node.sucLock.RLock()
	suc := node.successors.Beg()
	node.sucLock.RUnlock()

	// Devuelve el sucesor de este nodo
	res := &chord.GetSuccessorResponse{Successor: suc}
	return res, nil
}

// SetPredecessor establece el predecesor de este nodo.
func (node *Node) SetPredecessor(ctx context.Context, req *chord.SetPredecessorRequest) (*chord.SetPredecessorResponse, error) {
	log.Trace("Fijando el sucesor del nodo.")

	candidate := req.GetPredecessor()
	log.Tracef("Estableciendo el nodo predecesor a %s.", candidate.IP)

	// Si el nuevo predecesor no es este nodo, se actualiza el nuevo predecesor
	if !Equals(candidate.ID, node.ID) {
		// Bloquea el predecesor para lectura y escritura, se desbloquea el finalizar
		node.predLock.Lock()
		old := node.predecessor
		node.predecessor = candidate
		node.predLock.Unlock()
		// Si existia un anterior predecesor se absorven sus claves
		go node.AbsorbPredecessorKeys(old)
	} else {
		log.Trace("El candidato a predecesor es el mismo nodo. No es necesario actualizar.")
	}

	return &chord.SetPredecessorResponse{}, nil
}

// SetSuccessor establece el sucesor de este nodo.
func (node *Node) SetSuccessor(ctx context.Context, req *chord.SetSuccessorRequest) (*chord.SetSuccessorResponse, error) {

	candidate := req.GetSuccessor()
	log.Tracef("Estableciendo el nodo sucesor a %s.", candidate.IP)

	// Si el nuevo sucesor es distinto al actual nodo, se actualiza
	if !Equals(candidate.ID, node.ID) {
		//Bloquea el sucesor para escribir en el, se desbloquea al terminar
		node.sucLock.Lock()
		node.successors.PushBeg(candidate)
		node.sucLock.Unlock()
		// Actualiza este nuevo sucesor con las llaves de este nodo
		go node.UpdateSuccessorKeys()
	} else {
		log.Trace("Candidato a sucesor es el mismo nodo. No hay necesidad de actualizar.")
	}

	return &chord.SetSuccessorResponse{}, nil
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
		log.Debugf("Predecesor actualizado al nodo en %s.", candidate.IP)

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

// Posible Metodo a modificar
// Partition devuelve todos los pares <key, values>  de este almacenamiento local
func (node *Node) Partition(ctx context.Context, req *chord.PartitionRequest) (*chord.PartitionResponse, error) {
	log.Trace("Obteniendo todos las pares <key, values> en el almacenamiento local.")

	//Bloquea el predecesor para leer de el, lo desbloquea al terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	//Bloquea el predecesor para leer de el, lo desbloquea al terminar
	node.dictLock.RLock()
	// Obtiene los pares <key, value> del almacenamiento.
	inFiles, outFiles, err1 := node.dictionary.PartitionFile(pred.ID, node.ID)
	inTags, outTags, err2 := node.dictionary.PartitionTag(pred.ID, node.ID)
	node.dictLock.RUnlock()
	if err1 != nil {
		log.Error("Error obteniendo los archivos del almacenamiento local.")
		return &chord.PartitionResponse{}, errors.New("error obteniendo los archivos del almacenamiento local.\n" + err1.Error())
	}
	if err2 != nil {
		log.Error("Error obteniendo las etiquetas del almacenamiento local.")
		return &chord.PartitionResponse{}, errors.New("error obteniendo las etiquetas del almacenamiento local.\n" + err2.Error())
	}

	/*
		Devuelve el diccionario correspondiente al almacenamiento local de este nodo, y el correspondiente
		al almacenamiento local de replicacion
	*/
	return &chord.PartitionResponse{InFiles: inFiles, OutFiles: outFiles, InTags: inTags, OutTags: outTags}, nil
}

// Posible Metodo a modificar
// Extend agrega al diccionario local un nuevo conjunto de pares<key, values> .
func (node *Node) Extend(ctx context.Context, req *chord.ExtendRequest) (*chord.ExtendResponse, error) {
	log.Debug("Agregando nuevos elementos al almacenamiento local.")

	// Bloquea el diccionario para escribir en el , al terminar se desbloquea
	node.dictLock.Lock()
	err1 := node.dictionary.ExtendFiles(req.Files) // Agrega los pares <key, value> al almacenamiento.
	err2 := node.dictionary.ExtendTags(req.Tags)
	node.dictLock.Unlock()
	if err1 != nil {
		log.Error("Error agregando los archivos al almacenamiento local.")
		return &chord.ExtendResponse{}, errors.New("error agregando los archivos al almacenamiento local\n" + err1.Error())
	}
	if err2 != nil {
		log.Error("Error agregando las etiquetas al almacenamiento local.")
		return &chord.ExtendResponse{}, errors.New("error agregando las etiquetas al almacenamiento local\n" + err1.Error())
	}
	return &chord.ExtendResponse{}, nil
}

// Posible Metodo a modificar
// Discard a list of keys from local storage dictionary.
func (node *Node) Discard(ctx context.Context, req *chord.DiscardRequest) (*chord.DiscardResponse, error) {
	log.Debug("Descartando llaves desde el diccionario de almacenamiento local.")

	//Bloquea el diccionario para escribir en el, se desbloquea el final
	node.dictLock.Lock()

	// Elimina los archivos del almacenamiento.
	err1 := node.dictionary.DiscardFiles(req.Files)
	err2 := node.dictionary.DiscardTags(req.Tags)
	node.dictLock.Unlock()

	if err1 != nil {
		log.Error("Error descartando las llaves del diccionario de almacenamiento.")
		return &chord.DiscardResponse{}, errors.New("error descartando las llaves del diccionario de almacenamiento\n" + err1.Error())
	}
	// Elimina las etiquetas del almacenamiento.
	if err2 != nil {
		log.Error("Error descartando las llaves del diccionario de almacenamiento.")
		return &chord.DiscardResponse{}, errors.New("error descartando las llaves del diccionario de almacenamiento\n" + err2.Error())
	}
	return &chord.DiscardResponse{}, nil
}

/*
Metodos propios de la aplicacion
*/

// AddFile almacena  un fichero en el almacenamiento local .
func (node *Node) AddFile(ctx context.Context, req *chord.AddFileRequest) (*chord.AddFileResponse, error) {

	file := req.GetFile()

	log.Infof("Almacenar fichero %s.", file.GetName())

	// Si la request es una replica se resuelve de forma local.
	if req.Replica {
		log.Debug("Resolviendo la request de forma local (replicacion).")

		// Bloquea el diccionario para escribir en el, se desbloquea el terminar
		node.dictLock.Lock()
		// Almacena el fichero en el almacenamiento
		//err := node.dictionary.SetWithLock(file, address)
		err := node.dictionary.SetFile(file)
		node.dictLock.Unlock()

		if err != nil {
			log.Error("Error almacenado el fichero.")
			return &chord.AddFileResponse{}, errors.New("error almacenado la llave.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error almacenado la llave: ya esta bloqueada.\n" + err.Error())
			return &chord.AddFileResponse{}, err
		}

		log.Info("Se almaceno de forma exitosa.")
		//new_node := chord.Node{ID: node.ID, IP: node.IP, Port: node.Port}
		return &chord.AddFileResponse{Destine: node.Node}, nil
	}

	// Por defecto,toma este nodo para almacenar el archivo en el almacenamiento local.
	keyNode := node.Node
	// Bloquea el predecesor para leer de el, se desbloquea el terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	/*
		Si la ID de la llave no esta entre la ID de este nodo y su predecesor
		entonces la request no es necesariamente local
	*/
	if between, err := KeyBetween(file.Name, node.config.Hash, pred.ID, node.ID); !between && err == nil {
		log.Debug("Buscando por el nodo correspondiente.")
		// Localiza el nodo que corresponde a esta llave

		keyNode, err = node.LocateKey(file.Name)
		if err != nil {
			log.Error("Error estableciendo la llave.")
			return &chord.AddFileResponse{}, errors.New("Error estableciendo la llave.\n" + err.Error())
		}
	} else if err != nil {
		log.Error("Error estableciendo la llave.")
		return &chord.AddFileResponse{}, errors.New("Error estableciendo la llave.\n" + err.Error())
	}

	// Si la llave corresponde a este nodo, directamente almacena el fichero de forma local.
	if Equals(keyNode.ID, node.ID) {
		log.Debug("Resolviendo la request de forma local.")

		// Bloquea el diccionario para escribir en el, se desbloquea el terminar
		node.dictLock.Lock()
		//Almacena el fichero de forma local .
		//err := node.dictionary.SetWithLock(req.Key, req.Value, address)
		err := node.dictionary.SetFile(file)
		node.dictLock.Unlock()
		if err != nil {
			log.Error("Error almacenando el fichero.")
			return &chord.AddFileResponse{}, errors.New("error almacenando la llave.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error almacenando la llave: ya está bloqueada.\n" + err.Error())
			return &chord.AddFileResponse{}, err
		}

		log.Info("Resolucion exitosa.")

		// Bloquea el sucesor para leer de el, se desbloquea el terminar
		node.sucLock.RLock()
		suc := node.successors.Beg()
		node.sucLock.RUnlock()

		/*
			Antes de terminar se recorren las etiquetas y por cada una se almacena la informacion
			de donde se almaceno el fichero actual
		*/
		tags := file.Tags
		for _, tag := range tags {
			// Encuentro el nodo donde se debe almacenar la etiqueta
			temp_node, err := node.LocateKey(tag)
			// Si hubo algun error encontrando dicho nodo
			if err != nil {
				log.Error("Error localizando el nodo correspondiente a la etiqueta %s\n.", tag)
				return &chord.AddFileResponse{}, errors.New(fmt.Sprintln("Error localizando el nodo correspondiente a la etiqueta %s\n.", tag))
			}
			tagRequest := &chord.AddTagRequest{
				Tag:          tag,
				FileName:     file.Name,
				ExtensioName: file.Extension,
				TargetNode:   keyNode,
				Replica:      false,
			}
			// Si hubo un error almacenando alguna etiqueta
			err = node.RPC.AddTag(temp_node, tagRequest)
			if err != nil {
				log.Error("Error almacenando la etiqueta %s\n.", tag)
				return &chord.AddFileResponse{}, errors.New(fmt.Sprintf("Error almacenando la etiqueta %s\n.", tag))
			}
		}

		//Si el sucesor no es este nodo, replica la request a el.
		if !Equals(suc.ID, node.ID) {
			go func() {
				req.Replica = true
				log.Debugf("Replicando la request set a %s.", suc.IP)
				err := node.RPC.AddFile(suc, req)
				if err != nil {
					log.Errorf("Error replicando la request a %s.\n%s", suc.IP, err.Error())
				}
			}()
		}

		return &chord.AddFileResponse{Destine: node.Node}, nil
	} else {
		log.Infof("Redirigiendo la request a %s.", keyNode.IP)
		// En otro caso, devuelve el resultado de la llamada remota en el nodo correspondiente.
		return &chord.AddFileResponse{Destine: keyNode}, node.RPC.AddFile(keyNode, req)
	}

}

// AddFile almacena  una etiqueta y la informacion relevante en el almacenamiento local .
func (node *Node) AddTag(ctx context.Context, req *chord.AddTagRequest) (*chord.AddTagResponse, error) {
	log.Infof("Establece: llave=%s.", req.Tag)

	// Si la request es una replica se resuelve de forma local.
	if req.Replica {
		log.Debug("Resolviendo la request de forma local (replicacion).")

		// Bloquea el diccionario para escribir en el, se desbloquea el terminar
		node.dictLock.Lock()
		// Almacena la etiqueta y la informacion relevante en el almacenamiento local.
		err := node.dictionary.SetTag(req.Tag, req.FileName, req.ExtensioName, req.TargetNode)
		node.dictLock.Unlock()
		if err != nil {
			log.Errorf("Error almacenando la replica de la etiquetas: en el nod %v.", node)
			return &chord.AddTagResponse{}, errors.New("error almacenado la replica de etiqueta .\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error almacenado la llave: ya esta bloqueada.\n" + err.Error())
			return &chord.AddTagResponse{}, err
		}

		log.Info("Se almaceno de forma exitosa.")
		return &chord.AddTagResponse{}, nil
	}

	// Por defecto,toma este nodo para almacenar la nueva etiqueta  en el almacenamiento local.
	keyNode := node.Node
	// Bloquea el predecesor para leer de el, se desbloquea el terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	/*
		Si la ID de la llave no esta entre la ID de este nodo y su predecesor
		entonces la request no es necesariamente local
	*/
	if between, err := KeyBetween(req.Tag, node.config.Hash, pred.ID, node.ID); !between && err == nil {
		log.Debug("Buscando por el nodo correspondiente.")
		// LLocaliza el nodo que corresponde a esta llave
		keyNode, err = node.LocateKey(req.Tag)
		if err != nil {
			log.Error("Error estableciendo la llave.")
			return &chord.AddTagResponse{}, errors.New("error estableciendo la llave.\n" + err.Error())
		}
	} else if err != nil {
		log.Error("Error estableciendo la llave.")
		return &chord.AddTagResponse{}, errors.New("error estableciendo la llave.\n" + err.Error())
	}

	// Si la llave corresponde a este nodo , directamente almacena la llave de forma local.
	if Equals(keyNode.ID, node.ID) {
		log.Debug("Resolviendo la request de forma local.")

		// Bloquea el diccionario para escribir en el, se desbloquea el terminar
		node.dictLock.Lock()
		//Almacena la etiqueta .
		err := node.dictionary.SetTag(req.Tag, req.FileName, req.ExtensioName, req.TargetNode)
		node.dictLock.Unlock()
		if err != nil {
			log.Error("Error almacenando la etiqueta.")
			return &chord.AddTagResponse{}, errors.New("error almacenando la llave.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error almacenando la llave: ya está bloqueada.\n" + err.Error())
			return &chord.AddTagResponse{}, err
		}

		log.Info("Resolucion exitosa.")

		// Bloquea el sucesor para leer de el, se desbloquea el terminar
		node.sucLock.RLock()
		suc := node.successors.Beg()
		node.sucLock.RUnlock()

		//Si el sucesor no es este nodo, replica la request a el.
		if !Equals(suc.ID, node.ID) {
			go func() {
				req.Replica = true
				log.Debugf("Replicando la request set a %s.", suc.IP)
				err := node.RPC.AddTag(suc, req)
				if err != nil {
					log.Errorf("Error error replicando la request a %s.\n%s", suc.IP, err.Error())
				}
			}()
		}

		return &chord.AddTagResponse{}, nil
	} else {
		log.Infof("Redirigiendo la request a %s.", keyNode.IP)
		// En otro caso, devuelve el resultado de la llamada remota en el nodo correspondiente.
		return &chord.AddTagResponse{}, node.RPC.AddTag(keyNode, req)
	}

}

// Elimina un fichero del almacenamiento.
func (node *Node) DeleteFile(ctx context.Context, req *chord.DeleteFileRequest) (*chord.DeleteFileResponse, error) {
	log.Infof("Elimina: el archivo=%s.", req.FileName)
	// Si la request es una replica se resuelve local
	if req.Replica {
		log.Debug("Resolviendo la request Delete de forma local (replicacion).")

		// Bloquea el diccionario para escribir en el, lo desbloquea al terminar.
		node.dictLock.Lock()
		// Elimina el par <key, value> del almacenamiento.
		err := node.dictionary.DeleteFile(req.FileName, req.FileExtension)
		node.dictLock.Unlock()
		if err != nil && err != os.ErrPermission {
			log.Error("Error eliminando el fichero.")
			return &chord.DeleteFileResponse{}, errors.New("error eliminando el fichero.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error eliminando el fichero: ya esta bloqueado.\n" + err.Error())
			return &chord.DeleteFileResponse{}, err
		}

		log.Info("Eliminacion exitosa.")
		return &chord.DeleteFileResponse{}, nil
	}

	// Por defecto, se toma este nodo para eliminar el par <key, value> del almacenamiento local.
	keyNode := node.Node
	// Bloquea el predecesor para poder leer de el, se desbloquea al terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	/*
		Si el ID correspondiente al nombre del fichero no esta entre la ID de este nodo y la de su predecesor
		entonces la request no es necesariamente local
	*/
	if between, err := KeyBetween(req.FileName, node.config.Hash, pred.ID, node.ID); !between && err == nil {
		log.Debug("Buscando por el nodo correspondiente.")
		// Localiza el nodo que almacena el fichero
		keyNode, err = node.LocateKey(req.FileName)
		if err != nil {
			log.Error("Error eliminando el fichero.")
			return &chord.DeleteFileResponse{}, errors.New("error eliminando el fichero.\n" + err.Error())
		}
	} else if err != nil {
		log.Error("Error eliminando el fichero.")
		return &chord.DeleteFileResponse{}, errors.New("error eliminando el fichero.\n" + err.Error())
	}

	// Si el fichero corresponde a este nodo, se elimina directamente del almacenamiento.
	if Equals(keyNode.ID, node.ID) {
		log.Debug("Resolviendo la request delete de forma local.")

		// Bloquea el diccionario para escribir en el y se desbloquea al terminar la funcion.
		node.dictLock.Lock()
		// Elimina el fichero .
		err := node.dictionary.DeleteFile(req.FileName, req.FileExtension)
		node.dictLock.Unlock()
		if err != nil && err != os.ErrPermission {
			log.Error("Error eliminando el fichero.")
			return &chord.DeleteFileResponse{}, errors.New("error eliminando el fichero.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error eliminando el fichero: ya esta bloqueado.\n" + err.Error())
			return &chord.DeleteFileResponse{}, err
		}

		log.Info("Eliminacion exitosa.")

		// Bloquea el sucesor para leer de el, se desbloquea al terminar.
		node.sucLock.RLock()
		suc := node.successors.Beg()
		node.sucLock.RUnlock()

		// Si el sucesor no es este nodo, se replica la request para el
		if Equals(suc.ID, node.ID) {
			go func() {
				req.Replica = true
				log.Debugf("Replicando la request delete para %s.", suc.IP)
				err := node.RPC.DeleteFile(suc, req)
				if err != nil {
					log.Errorf("Error replicando la request delete para %s.\n%s", suc.IP, err.Error())
				}
			}()
		}
		// En otro caso se devuelve.
		return &chord.DeleteFileResponse{}, nil
	} else {
		log.Infof("Redirigiendo la request delete para %s.", keyNode.IP)
		// En otro caso, se devuelve el resultado de la llamada remota en el nodo correspondiente
		return &chord.DeleteFileResponse{}, node.RPC.DeleteFile(keyNode, req)

	}

}

func (node *Node) DeleteFileByQuery(ctx context.Context, req *chord.DeleteFileByQueryRequest) (*chord.DeleteFileByQueryResponse, error) {

	// Objeto poder llevar un registro de los ficheros, sabiendo en cuantas querys está presente
	querys := make(map[string][]string)
	target := make(map[string]*chord.Node)

	/*
		Lo primero que se necesita es poder obtener toda la informacion de cada uno de las etiquetas.
		Para eso debe haber una parte que por cada etiqueta pida todos los ficheros que tiene
	*/

	tags := req.Tag
	for _, tag := range tags {

		// Por defecto, se toma este nodo para eliminar el par <key, value> del almacenamiento local.
		keyNode := node.Node
		// Bloquea el predecesor para poder leer de el, se desbloquea al terminar
		node.predLock.RLock()
		pred := node.predecessor
		node.predLock.RUnlock()
		// Encuentra en donde esta ubicada dicha etiqueta
		if between, err := KeyBetween(tag, node.config.Hash, pred.ID, node.ID); !between && err == nil {
			log.Debug("Buscando por el nodo correspondiente.")
			// Localiza el nodo que almacena el fichero
			keyNode, err = node.LocateKey(tag)
			if err != nil {
				log.Error("Error encontrando el nodo.")
				return &chord.DeleteFileByQueryResponse{}, errors.New("error encontrando el nodo.\n" + err.Error())
			}
		} else if err != nil {
			log.Error("Error encontrando el nodo.")
			return &chord.DeleteFileByQueryResponse{}, errors.New("error encontrando el nodo.\n" + err.Error())
		}

		// Si pasa por aqui es q ya tiene nodo al que buscar
		// Si el fichero corresponde a este nodo, se elimina directamente del almacenamiento.
		if Equals(keyNode.ID, node.ID) {
			log.Debug("Resolviendo la request de forma local.")

			// Bloquea el diccionario para escribir en el y se desbloquea al terminar la funcion.
			node.dictLock.Lock()
			// obteniendo la informacion de la etiqueta .
			values, err := node.dictionary.GetTag(tag)
			node.dictLock.Unlock()
			if err != nil && err != os.ErrPermission {
				log.Error("Error recuperando la informacion de las etiquetas.")
				return &chord.DeleteFileByQueryResponse{}, errors.New("error recuperando la informacion de las etiquetas.\n" + err.Error())
			} else if err == os.ErrPermission {
				log.Error("Error recuperando la informacion de las etiquetas: esta bloqueada.\n" + err.Error())
				return &chord.DeleteFileByQueryResponse{}, err
			}

			log.Info("Recuperacion exitosa.")
			// Aqui vendria lo que se hace una vez con la informacion de las tags
			for _, value := range values {
				file := value.FileName + "." + value.FileExtension
				tempNode := &chord.Node{ID: value.NodeID, IP: value.NodeIP, Port: value.NodePort}

				if list, ok := querys[file]; ok {
					list = append(list, tag)
					querys[file] = list
				} else {
					querys[file] = []string{tag}
					target[file] = tempNode
				}
			}

		} else {
			log.Infof("Redirigiendo la request delete para %s.", keyNode.IP)
			// En otro caso, se devuelve el resultado de la llamada remota en el nodo correspondiente

			getTagRequest := &chord.GetTagRequest{Tag: tag}
			res, err := node.RPC.GetTag(keyNode, getTagRequest)
			if err != nil {

				log.Error("Error al recibr los TagEncoders")
				return &chord.DeleteFileByQueryResponse{}, errors.New("Error al recibr los TagEncoders\n" + err.Error())
			}
			// Se tiene la respuesta y se esta trabajando con ella
			for _, value := range res {
				file := value.FileName + "." + value.FileExtension
				tempNode := &chord.Node{ID: value.NodeID, IP: value.NodeIP, Port: value.NodePort}

				if list, ok := querys[file]; ok {
					list = append(list, tag)
					querys[file] = list
				} else {
					querys[file] = []string{tag}
					target[file] = tempNode
				}
			}
			/*
				Hasta aqui esta la implementacion general, que se deberá usar en los otros metodos
			*/
		}

		/*

			A partir de aqui se debe tener listos tanto el diccionario de querys como el de target

		*/
		// Primero se eliminan los archivos
		for key, value := range querys {
			// Esto representa la interseccion de todas las querys, o sea se procesa si
			// las cumple todas
			if len(value) == len(tags) {
				// Separa la clave en el nombre del fichero y su extension
				idef := strings.Split(key, ".")
				//Nombre del fichero
				name := idef[0]
				// Nombre de la extension
				extension := idef[1]

				req := &chord.DeleteFileRequest{
					FileName:      name,
					FileExtension: extension,
					Replica:       false,
				}
				// Llama a eliminar el fichero
				err := node.RPC.DeleteFile(target[key], req)
				if err != nil {
					return &chord.DeleteFileByQueryResponse{}, err
				}
				for _, tag := range value {
					req := &chord.DeleteFileFromTagRequest{
						Tag:           tag,
						FileName:      name,
						FileExtension: extension,
						Replica:       false,
					}
					// Eliminar la informacion del fichero en las etiquetas
					err := node.RPC.DeleteFileFromTag(target[key], req)

					if err != nil {
						return &chord.DeleteFileByQueryResponse{}, err
					}
				}
			}
		}
	}
	return &chord.DeleteFileByQueryResponse{}, nil
}
func (node *Node) GetTag(req *chord.GetTagRequest, stream chord.Chord_GetTagServer) error {

	log.Infof("Obtener la informacion relacionada a la etiqueta: %s.", req.Tag)

	// Por defecto, toma este nodo para obtener el valor de esta llave del almacenamiento local.
	keyNode := node.Node
	// Bloquea el predecesor para leer de el y al terminar lo desbloquea.
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	/*
	 Si el identificador de la llave no esta entre la ID del nodo y la ID de su predecesor
	 entonces la llave requerida no esta local necesariamente
	*/
	if between, err := KeyBetween(req.Tag, node.config.Hash, pred.ID, node.ID); !between && err == nil {
		log.Debug("Buscando por el nodo correspondiente.")
		// Localiza el nodo que almacena la llave
		keyNode, err = node.LocateKey(req.Tag)
		if err != nil {
			log.Errorf("Error localizando el nodo correspondiente a la etiqueta: %s\n.", req.Tag)
			return errors.New("error localizando el nodo correspondiente a la etiqueta: " + req.Tag + "\n" + err.Error() + "\n")
		}
	} else if err != nil {
		log.Errorf("Error localizando el nodo correspondiente a la etiqueta: %s\n.", req.Tag)
		return errors.New("error localizando el nodo correspondiente a la etiqueta: " + req.Tag + "\n" + err.Error() + "\n")
	}

	/*
		Si el nodo que almacena la llave es este nodo, entonces consigue el valor asociado desde el almacenamiento del nodo
	*/
	if Equals(keyNode.ID, node.ID) {
		log.Debug("Resolviendo la request de forma local.")

		//Bloquea el diccionario para leer de el, se desbloqua al terminar.
		node.dictLock.RLock()
		//Consigue el valor asociado a esta llave desde el almacenamiento
		encoding, err := node.dictionary.GetTag(req.Tag)
		node.dictLock.RUnlock()
		if err != nil && err != os.ErrPermission {
			log.Error("Error consiguiendo la informacion.\n" + err.Error())
			return os.ErrNotExist
		} else if err == os.ErrPermission {
			log.Error("Error consiguiendo la llave: ya esta bloqueado.\n" + err.Error())
			return err
		}

		log.Info("Se recupero de forma exitosa.")

		for _, info := range encoding {

			encoder :=
				&chord.TagEncoder{
					FileName:      info.FileName,
					FileExtension: info.FileExtension,
					NodeID:        info.NodeID,
					NodeIP:        info.NodeIP,
					NodePort:      info.NodePort,
				}
			res := &chord.GetTagResponse{Encoder: encoder}
			stream.Send(res)
		}
		return nil
	} else {
		log.Infof("Redirigiendo la request a %s.", keyNode.IP)
	}
	// En otro caso, devuelve  el resultado de la llamada remota al nodo correspondiente.
	res, err := node.RPC.GetTag(keyNode, req)
	if err != nil {
		log.Error("Error consiguiendo la informacion.\n" + err.Error())
		return err
	}
	for _, info := range res {

		encoder :=
			&chord.TagEncoder{
				FileName:      info.FileName,
				FileExtension: info.FileExtension,
				NodeID:        info.NodeID,
				NodeIP:        info.NodeIP,
				NodePort:      info.NodePort,
			}
		res := &chord.GetTagResponse{Encoder: encoder}
		stream.Send(res)
	}
	return nil
}

func (node *Node) DeleteFileFromTag(ctx context.Context, req *chord.DeleteFileFromTagRequest) (*chord.DeleteFileFromTagResponse, error) {
	log.Infof("Elimina la informacion del fichero=%s en la etiqueta %s\n.", req.FileName, req.Tag)
	// Si la request es una replica se resuelve local
	if req.Replica {
		log.Debug("Resolviendo la request Delete de forma local (replicacion).")

		// Bloquea el diccionario para escribir en el, lo desbloquea al terminar.
		node.dictLock.Lock()
		// Elimina el par <key, value> del almacenamiento.
		err := node.dictionary.DeleteFileFromTag(req.Tag, req.FileName, req.FileExtension)
		node.dictLock.Unlock()
		if err != nil && err != os.ErrPermission {
			log.Error("Error eliminando la informacion.")
			return &chord.DeleteFileFromTagResponse{}, errors.New("error eliminando la informacion.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error eliminando la informacion: ya esta bloqueado.\n" + err.Error())
			return &chord.DeleteFileFromTagResponse{}, err
		}

		log.Info("Eliminacion exitosa.")
		return &chord.DeleteFileFromTagResponse{}, nil
	}

	// Por defecto, se toma este nodo para eliminar el par <key, value> del almacenamiento local.
	keyNode := node.Node
	// Bloquea el predecesor para poder leer de el, se desbloquea al terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	/*
		Si el ID correspondiente al nombre del fichero no esta entre la ID de este nodo y la de su predecesor
		entonces la request no es necesariamente local
	*/
	if between, err := KeyBetween(req.Tag, node.config.Hash, pred.ID, node.ID); !between && err == nil {
		log.Debug("Buscando por el nodo correspondiente.")
		// Localiza el nodo que almacena el fichero
		keyNode, err = node.LocateKey(req.FileName)
		if err != nil {
			log.Error("Error buscando el nodo correspondiente a la etiqueta: %s.", req.Tag)
			return &chord.DeleteFileFromTagResponse{}, errors.New("error buscando el nodo correspondiente a la etiqueta: " + req.Tag + "\n" + err.Error())
		}
	} else if err != nil {
		log.Error("Error comparando los ID.")
		return &chord.DeleteFileFromTagResponse{}, errors.New("error comparando los ID.\n" + err.Error())
	}

	// Si el fichero corresponde a este nodo, se elimina directamente del almacenamiento.
	if Equals(keyNode.ID, node.ID) {
		log.Debug("Resolviendo la request DeleteTag de forma local.")

		// Bloquea el diccionario para escribir en el y se desbloquea al terminar la funcion.
		node.dictLock.Lock()
		// Elimina la informacion .
		err := node.dictionary.DeleteFileFromTag(req.Tag, req.FileName, req.FileExtension)
		node.dictLock.Unlock()
		if err != nil && err != os.ErrPermission {
			log.Error("Error eliminando el fichero.")
			return &chord.DeleteFileFromTagResponse{}, errors.New("error eliminando el fichero.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error eliminando el fichero: ya esta bloqueado.\n" + err.Error())
			return &chord.DeleteFileFromTagResponse{}, err
		}

		log.Info("Eliminacion exitosa.")

		// Bloquea el sucesor para leer de el, se desbloquea al terminar.
		node.sucLock.RLock()
		suc := node.successors.Beg()
		node.sucLock.RUnlock()

		// Si el sucesor no es este nodo, se replica la request para el
		if Equals(suc.ID, node.ID) {
			go func() {
				req.Replica = true
				log.Debugf("Replicando la request delete para %s.", suc.IP)
				err := node.RPC.DeleteFileFromTag(suc, req)
				if err != nil {
					log.Errorf("Error replicando la request delete para %s.\n%s", suc.IP, err.Error())
				}
			}()
		}
		// En otro caso se devuelve.
		return &chord.DeleteFileFromTagResponse{}, nil
	} else {
		log.Infof("Redirigiendo la request delete para %s.", keyNode.IP)
		// En otro caso, se devuelve el resultado de la llamada remota en el nodo correspondiente
		return &chord.DeleteFileFromTagResponse{}, node.RPC.DeleteFileFromTag(keyNode, req)

	}

}
func (node *Node) DeleteTag(ctx context.Context, req *chord.DeleteTagRequest) (*chord.DeleteTagResponse, error) {
	log.Infof("Elimina la informacion de la etiqueta %s en el almacenamiento local\n.", req.Tag)
	// Si la request es una replica se resuelve local
	if req.Replica {
		log.Debug("Resolviendo la request Delete de forma local (replicacion).")

		// Bloquea el diccionario para escribir en el, lo desbloquea al terminar.
		node.dictLock.Lock()
		// Elimina el par <key, value> del almacenamiento.
		err := node.dictionary.DeleteTag(req.Tag)
		node.dictLock.Unlock()
		if err != nil && err != os.ErrPermission {
			log.Error("Error eliminando la informacion.")
			return &chord.DeleteTagResponse{}, errors.New("error eliminando la informacion.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error eliminando la informacion: ya esta bloqueado.\n" + err.Error())
			return &chord.DeleteTagResponse{}, err
		}

		log.Info("Eliminacion exitosa.")
		return &chord.DeleteTagResponse{}, nil
	}

	// Por defecto, se toma este nodo para eliminar el par <key, value> del almacenamiento local.
	keyNode := node.Node
	// Bloquea el predecesor para poder leer de el, se desbloquea al terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	/*
		Si el ID correspondiente al nombre del fichero no esta entre la ID de este nodo y la de su predecesor
		entonces la request no es necesariamente local
	*/
	if between, err := KeyBetween(req.Tag, node.config.Hash, pred.ID, node.ID); !between && err == nil {
		log.Debug("Buscando por el nodo correspondiente.")
		// Localiza el nodo que almacena el fichero
		keyNode, err = node.LocateKey(req.Tag)
		if err != nil {
			log.Error("Error buscando el nodo correspondiente a la etiqueta: %s.", req.Tag)
			return &chord.DeleteTagResponse{}, errors.New("error buscando el nodo correspondiente a la etiqueta: " + req.Tag + "\n" + err.Error())
		}
	} else if err != nil {
		log.Error("Error comparando los ID.")
		return &chord.DeleteTagResponse{}, errors.New("error comparando los ID.\n" + err.Error())
	}

	// Si el fichero corresponde a este nodo, se elimina directamente del almacenamiento.
	if Equals(keyNode.ID, node.ID) {
		log.Debug("Resolviendo la request DeleteTag de forma local.")

		// Bloquea el diccionario para escribir en el y se desbloquea al terminar la funcion.
		node.dictLock.Lock()
		// Elimina la informacion .
		err := node.dictionary.DeleteTag(req.Tag)
		node.dictLock.Unlock()
		if err != nil && err != os.ErrPermission {
			log.Error("Error eliminando la etiqueta.")
			return &chord.DeleteTagResponse{}, errors.New("error eliminando la etiqueta.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error eliminando la etiqueta: ya esta bloqueado.\n" + err.Error())
			return &chord.DeleteTagResponse{}, err
		}

		log.Info("Eliminacion exitosa.")

		// Bloquea el sucesor para leer de el, se desbloquea al terminar.
		node.sucLock.RLock()
		suc := node.successors.Beg()
		node.sucLock.RUnlock()

		// Si el sucesor no es este nodo, se replica la request para el
		if Equals(suc.ID, node.ID) {
			go func() {
				req.Replica = true
				log.Debugf("Replicando la request delete para %s.", suc.IP)
				err := node.RPC.DeleteTag(suc, req)
				if err != nil {
					log.Errorf("Error replicando la request delete para %s.\n%s", suc.IP, err.Error())
				}
			}()
		}
		// En otro caso se devuelve.
		return &chord.DeleteTagResponse{}, nil
	} else {
		log.Infof("Redirigiendo la request delete para %s.", keyNode.IP)
		// En otro caso, se devuelve el resultado de la llamada remota en el nodo correspondiente
		return &chord.DeleteTagResponse{}, node.RPC.DeleteTag(keyNode, req)

	}

}
func (node *Node) EditFileFromTag(ctx context.Context, req *chord.EditFileFromTagRequest) (*chord.EditFileFromTagResponse, error) {
	log.Infof("Se esta modificando la etiqueta =%s.", req.Tag)

	// Por defecto, se toma este nodo .
	keyNode := node.Node
	// Bloquea el predecesor para poder leer de el, se desbloquea al terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	/*
		Si el ID correspondiente al nombre del fichero no esta entre la ID de este nodo y la de su predecesor
		entonces la request no es necesariamente local
	*/
	if between, err := KeyBetween(req.Tag, node.config.Hash, pred.ID, node.ID); !between && err == nil {
		log.Debug("Buscando por el nodo correspondiente.")
		// Localiza el nodo que almacena el fichero
		keyNode, err = node.LocateKey(req.Tag)
		if err != nil {
			log.Error("Error editando la etiqueta.")
			return &chord.EditFileFromTagResponse{}, errors.New("error editando la etiqueta.\n" + err.Error())
		}
	} else if err != nil {
		log.Error("Error editando la etiqueta.")
		return &chord.EditFileFromTagResponse{}, errors.New("error editando la etiqueta.\n" + err.Error())
	}

	// Si el fichero corresponde a este nodo, se Edita directamente del almacenamiento.
	if Equals(keyNode.ID, node.ID) {
		log.Debug("Resolviendo la request Edit de forma local.")

		// Bloquea el diccionario para escribir en el y se desbloquea al terminar la funcion.
		node.dictLock.Lock()
		// Edita el fichero .
		err := node.dictionary.EditFileFromTag(req.Tag, req.Mod)
		node.dictLock.Unlock()
		if err != nil && err != os.ErrPermission {
			log.Error("Error editando la etiqueta.")
			return &chord.EditFileFromTagResponse{}, errors.New("error editando la etiqueta.\n" + err.Error())
		} else if err == os.ErrPermission {
			log.Error("Error editando la etiqueta: ya esta bloqueado.\n" + err.Error())
			return &chord.EditFileFromTagResponse{}, err
		}

		log.Info("Se edito exitosamente.")

		return &chord.EditFileFromTagResponse{}, nil
	} else {
		log.Infof("Redirigiendo la request EditFileFromTag para %s.", keyNode.IP)
		// En otro caso, se devuelve el resultado de la llamada remota en el nodo correspondiente
		return &chord.EditFileFromTagResponse{}, node.RPC.EditFileFromTag(keyNode, req)

	}

}
