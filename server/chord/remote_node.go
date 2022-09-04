package chord

import (
	"time"

	"github.com/alejbv/SistemaDeFicherosDistribuido/chord"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// RemoteNode almacena las propiedades de una conexion con un nodo chord remoto.
type RemoteNode struct {
	chord.ChordClient // Cliente Chord conectado con un servidor remoto.

	addr       string           // Direccion del nodo remoto.
	conn       *grpc.ClientConn // conexion Grpc con el servidor remoto.
	lastActive time.Time        // Ultima vez que la conexion fue usada.
}

// CloseConnection cierra  la conexion con un RemoteNode.
func (connection *RemoteNode) CloseConnection() {
	err := connection.conn.Close() // Cierra la conexion con el nodo remoto.
	if err != nil {
		log.Error("Error cerrando la conexion con el nodo remoto.\n" + err.Error() + "\n")
		return
	}
}
