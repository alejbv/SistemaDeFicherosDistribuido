package chord

import (
	"net"
	"sync"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
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
