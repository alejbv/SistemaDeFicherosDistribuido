package chord

import (
	"hash"
	"net"

	log "github.com/sirupsen/logrus"
)

func IsOpen[T any](channel <-chan T) bool {
	select {
	case <-channel:
		return false
	default:
		if channel == nil {
			return false
		}
		return true
	}
}

// GetOutboundIP obtiene la  IP de este nodo en la red.
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func HashKey(key string, hash func() hash.Hash) ([]byte, error) {
	log.Trace("Hashing key: " + key + ".\n")
	h := hash()
	if _, err := h.Write([]byte(key)); err != nil {
		log.Error("Error hashing key " + key + ".\n" + err.Error() + ".\n")
		return nil, err
	}
	value := h.Sum(nil)

	return value, nil
}
