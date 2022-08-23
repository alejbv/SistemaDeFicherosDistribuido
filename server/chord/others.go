package chord

import (
	"bytes"
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

// Equals comprueba si 2 IDs son iguales
func Equals(ID1, ID2 []byte) bool {
	return bytes.Compare(ID1, ID2) == 0
}

func Keys[T any](dictionary map[string]T) []string {
	keys := make([]string, 0)

	for key := range dictionary {
		keys = append(keys, key)
	}

	return keys
}

func KeyBetween(key string, hash func() hash.Hash, L, R []byte) (bool, error) {
	ID, err := HashKey(key, hash) // Obtiene la ID correspondiente a la clave.
	if err != nil {
		return false, err
	}

	return Equals(ID, R) || Between(ID, L, R), nil
}

// Between comprueba si una ID esta dentro del intervalo (L,R),en el anillo chord.
func Between(ID, L, R []byte) bool {
	// Si L <= R, devuelve true si L < ID < R.
	if bytes.Compare(L, R) <= 0 {
		return bytes.Compare(L, ID) < 0 && bytes.Compare(ID, R) < 0
	}

	// Si L > R, es un segmento sobre el final del anillo.
	// Entonces, ID esta entre L y R si L < ID o ID < R.
	return bytes.Compare(L, ID) < 0 || bytes.Compare(ID, R) < 0
}
