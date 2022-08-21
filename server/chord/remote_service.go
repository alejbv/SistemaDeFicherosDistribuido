package chord

import (
	"sync"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
)

// RemoteServices enables a node to interact with other nodes in the ring, as a client of its servers.
type RemoteServices interface {
	Start() error // Start the services.
	Stop() error  // Stop the services.

	// GetPredecessor returns the node believed to be the current predecessor of a remote node.
	GetPredecessor(*chord.Node) (*chord.Node, error)
	// GetSuccessor returns the node believed to be the current successor of a remote node.
	GetSuccessor(*chord.Node) (*chord.Node, error)
	// SetPredecessor sets the predecessor of a remote node.
	SetPredecessor(*chord.Node, *chord.Node) error
	// SetSuccessor sets the successor of a remote node.
	SetSuccessor(*chord.Node, *chord.Node) error
	// FindSuccessor finds the node that succeeds this ID, starting at a remote node.
	FindSuccessor(*chord.Node, []byte) (*chord.Node, error)
	// Notify a remote node that it possibly have a new predecessor.
	Notify(*chord.Node, *chord.Node) error
	// Check if a remote node is alive.
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

// NewGRPCServices creates a new GRPCServices object.
func NewGRPCServices(config *Configuration) *GRPCServices {
	// Create the GRPCServices object.
	services := &GRPCServices{
		Configuration: config,
		connections:   nil,
		shutdown:      nil,
	}

	// Return the GRPCServices object.
	return services
}
