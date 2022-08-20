package chord

import (
	"errors"
	"math/big"
	"net"
	"sync"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	"github.com/cloudflare/cfssl/log"
	"google.golang.org/grpc"
)

type Node struct {
	*chord.Node // Real node.

	predecessor *chord.Node  // Predecessor of this node in the ring.
	predLock    sync.RWMutex // Locks the predecessor for reading or writing.
	//successors  *Queue[chord.Node] // Queue of successors of this node in the ring.
	sucLock sync.RWMutex // Locks the queue of successors for reading or writing.

	fingerTable FingerTable  // FingerTable of this node.
	fingerLock  sync.RWMutex // Locks the finger table for reading or writing.

	//	RPC    RemoteServices // Transport layer of this node.
	config *Configuration // General configurations.

	//	dictionary Storage      // Storage dictionary of this node.
	dictLock sync.RWMutex // Locks the dictionary for reading or writing.

	server   *grpc.Server     // Node server.
	sock     *net.TCPListener // Node server listener socket.
	shutdown chan struct{}    // Determine if the node server is actually running.

	chord.UnimplementedChordServer
}

// NewNode creates and returns a new Node.
func NewNode(port string, configuration *Configuration, transport RemoteServices, storage Storage) (*Node, error) {
	// If configuration is null, report error.
	if configuration == nil {
		log.Error("Error creating node: configuration cannot be null.")
		return nil, errors.New("error creating node: configuration cannot be null")
	}

	// Creates the new node with the obtained ID and same address.
	innerNode := &chord.Node{ID: big.NewInt(0).Bytes(), IP: "0.0.0.0", Port: port}

	// Instantiates the node.
	node := &Node{Node: innerNode,
		predecessor: nil,
		successors:  nil,
		fingerTable: nil,
		RPC:         transport,
		config:      configuration,
		dictionary:  storage,
		server:      nil,
		shutdown:    nil}

	// Return the node.
	return node, nil
}

// DefaultNode creates and returns a new Node with default configurations.
func DefaultNode(port string) (*Node, error) {
	conf := DefaultConfig()                    // Creates a default configuration.
	transport := NewGRPCServices(conf)         // Creates a default RPC transport layer.
	dictionary := NewDiskDictionary(conf.Hash) // Creates a default dictionary.

	// Return the default node.
	return NewNode(port, conf, transport, dictionary)
}
