package chord

import (
	"crypto/sha1"
	"hash"
	"time"

	"google.golang.org/grpc"
)

// Configuracion de un node
type Configuration struct {
	Hash     func() hash.Hash //Funcion de hash a usar.
	HashSize int              // El tamaño del hash usado.

	ServerOpts []grpc.ServerOption // Opciones del Servidor.
	DialOpts   []grpc.DialOption   // Opciones de la Conexion.

	Timeout time.Duration // Timeout de la conexion.
	MaxIdle time.Duration // El máximo lifetime de una conexion.

	StabilizingNodes int // El número máximo de sucesores a registrar en un nodo para asegurar la estabilidad

}

// DefaultConfig crea y devuelve una configuracion por defecto
func DefaultConfig() *Configuration {
	Hash := sha1.New              // Usa sha1 como la funcion de hash.
	HashSize := Hash().Size() * 8 // Se almacena el tamaño de hash soportado.
	DialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(5 * time.Second),
		grpc.FailOnNonTempDialError(true),
		grpc.WithInsecure(),
	}

	// Construyendo la configuracion.
	config := &Configuration{
		Hash:             sha1.New,
		HashSize:         HashSize,
		DialOpts:         DialOpts,
		Timeout:          10 * time.Second,
		MaxIdle:          5 * time.Second,
		StabilizingNodes: 5,
	}

	return config
}
