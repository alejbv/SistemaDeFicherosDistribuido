package chord

import (
	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
)

// Definicion de la FingerTable.
type FingerTable []*chord.Node

// NewFingerTable crea y devuelve una nueva FingerTable.
func NewFingerTable(size int) FingerTable {
	hand := make([]*chord.Node, size) // Se crea el nuevo arreglo.

	// Se inicializan todos en  null.
	for i := range hand {
		hand[i] = nil
	}

	return hand
}
