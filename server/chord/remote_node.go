package chord

import (
	"time"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// RemoteNode stores the properties of a connection with a remote chord node.
type RemoteNode struct {
	chord.ChordClient // Chord client connected with the remote node server.

	addr       string           // Address of the remote node.
	conn       *grpc.ClientConn // Grpc connection with the remote node server.
	lastActive time.Time        // Last time the connection was used.
}

// CloseConnection close the connection with a RemoteNode.
func (connection *RemoteNode) CloseConnection() {
	err := connection.conn.Close() // Close the connection with the remote node server.
	if err != nil {
		log.Error("Error closing connection with a remote node.\n" + err.Error() + "\n")
		return
	}
}
