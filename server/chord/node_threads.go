package chord

import (
	"time"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	log "github.com/sirupsen/logrus"
)

// Node server periodically threads.

/*
Estabiliza el nodo, Para esto el predecesor del sucesor es buscado.
Si el nodo obtenido no es este nodo, y es más cercano a este nodo que su actual sucesor,
se actualiza este nodo tomando este nodo recien descubierto como su nuevo sucesor.
Finalmente se notifica a el sucesor de este para que se puede actualizar el mismo
*/

func (node *Node) Stabilize() {
	log.Trace("Estabilizando el nodo.")

	// Bloquea el sucesor para leer de el, se desbloquea al terminar.
	node.sucLock.RLock()
	suc := node.successors.Beg()
	node.sucLock.RUnlock()

	// Si el sucesor es este nodo, no hay nada que estabilizar.
	if Equals(suc.ID, node.ID) {
		log.Trace("No es necesario estabilizar")
		return
	}

	candidate, err := node.RPC.GetPredecessor(suc) // En otro caso obten el predecesor de este nodo.
	if err != nil {
		log.Errorf("Error estabilizando el nodo. No se puede obtener el predecesor del sucesor en %s.\n%s", suc.IP, err.Error())
		return
	}

	/*
		Si el nodo candidato es mas cercano a este nodo que su actual sucesor, se actualiza el
		sucesor de este nodo con el candidato
	*/
	if Equals(node.ID, suc.ID) || Between(candidate.ID, node.ID, suc.ID) {
		log.Debug("Successor updated to node at " + candidate.IP + ".")
		// Bloquea el sucesor para escribir en el, se desbloquea el finalizar
		node.sucLock.Lock()
		node.successors.PushBeg(candidate) //Se actualiza el sucesor de este nodo con el obtenido.
		suc = candidate
		node.sucLock.Unlock()
	}

	// Notifica al sucesor de la existencia de su predecesor.
	err = node.RPC.Notify(suc, node.Node)
	if err != nil {
		log.Errorf("Error notifying successor at %s.\n%s", suc.IP, err.Error())
		return
	}

	log.Trace("Node stabilized.")
}

// PeriodicallyStabilize periodicamente estabiliza el nodo.
func (node *Node) PeriodicallyStabilize() {
	log.Debug("Empezado el hilo para estabilizar el nodo.")

	ticker := time.NewTicker(1 * time.Second) // Establece el tiempo de activacion entre rutinas.
	for {
		select {
		case <-node.shutdown: //Si el nodo esta caido, cierra el hilo.
			ticker.Stop()
			return
		case <-ticker.C: // Si ha transcurrido el tiempo, estabiliza el nodo.
			node.Stabilize()
		}
	}
}

/*
CheckPredecessor comprueba si el predecesor a fallado.
Para eso realiza una llamada remota al metodo Check en su predecesor.
En caso de fallo se asume que este está caido y se actualiza el nodo, absorviendo
las llaves de su predecesor. Estas llaves ya estan replicadas en el nodo, por lo que
las nuevas llaves son envidas al sucesor para mantener la replicacion
*/
func (node *Node) CheckPredecessor() {
	log.Trace("Comprobando el predecesor.")

	//Bloquea el predecesor para leer en el, al terminar se desbloquea.
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	// Si el predecesor de este nodo es distinto a este, comprueba si esta vivo.
	if !Equals(pred.ID, node.ID) {
		err := node.RPC.Check(pred)
		// En caso de error, se asume que el predecesor no esta activo.
		if err != nil {
			log.Errorf("Predecesor en %s a fallado.\n%s", pred.IP, err.Error())
			// Bloquea el predecesor para escribir en el, se desbloquea despues.
			node.predLock.Lock()
			node.predecessor = node.Node
			node.predLock.Unlock()
			// Si existia un predecesor viejo, absorce sus llaves.
			go node.AbsorbPredecessorKeys(pred)
		} else {
			log.Trace("Predecesor Activo.")
		}
	} else {
		log.Trace("No existe predecesor.")
	}
}

// PeriodicallyCheckPredecessor comprueba periodicamente si el predecesor a fallado.
func (node *Node) PeriodicallyCheckPredecessor() {
	log.Debug("Empezado el hilo para comprobar el predecesor.")

	ticker := time.NewTicker(500 * time.Millisecond) // Establece el tiempo entre la activacion de las rutinas.
	for {
		select {
		case <-node.shutdown: // Si el nodo esta caido, se cierra el hilo.
			ticker.Stop()
			return
		case <-ticker.C: // Si ha transcurrido el tiempo comprueba si el predecesor esta activo.
			node.CheckPredecessor()
		}
	}
}

/*
CheckSuccessor comprueba si es sucesor a fallado.
Para esto realiza una llamada remota a check del sucesor. Si la llamada falla, se asume
que el sucesor a fallado y es removido de la cola de sucesores y es necesario reemplazarlo
Es neceario transferir las llaves de este nodo a su nuevo sucesor para mantener la replicacion, ya que
este nuevo sucesor solamente tiene sus llaves y las correspondientes al viejo sucesor.
*/

func (node *Node) CheckSuccessor() {
	log.Trace("Comprobando sucesor.")

	// Bloquea el sucesor para leer de el, al terminar se desbloquea.
	node.sucLock.RLock()
	suc := node.successors.Beg()
	node.sucLock.RUnlock()

	//Bloquea el predecesor para leer de el, al terminar se desbloquea.
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	// Si el sucesor es distinto a este nodo, comprueba si esta vivo
	if !Equals(suc.ID, node.ID) {
		err := node.RPC.Check(suc)
		// De no estarlo se reemplaza
		if err != nil {
			// Bloquea el predecesor para escribir en el, al terminar se desbloquea.
			node.sucLock.Lock()
			node.successors.PopBeg() // Remueve el actual sucesor.
			// Se agrega el final este nodo para asegurar que la cola no este vacia.
			if node.successors.Empty() {
				node.successors.PushBack(node.Node)
			}
			suc = node.successors.Beg() //Toma el siguiente sucesor en la cola.
			node.sucLock.Unlock()
			log.Errorf("Sucesor en %s ha fallado.\n%s", suc.IP, err.Error())
		} else {
			// If successor is alive, return.
			log.Trace("Successor alive.")
			return
		}
	}

	//Si no hay sucesores pero si un predecesor, toma el sucesor como predecesor.
	if Equals(suc.ID, node.ID) {
		if !Equals(pred.ID, node.ID) {
			//Bloquea el sucesor para escribir en el, al terminar se desbloquea.
			node.sucLock.Lock()
			node.successors.PushBeg(pred)
			suc = node.successors.Beg() // Toma el siguiente sucesor.
			node.sucLock.Unlock()
		} else {
			// Si no hay predecesor tambien, no hay nada que hacer.
			log.Trace("No hay sucesor.")
			return
		}
	}

	// En otro caso reporta que hay un nuevo sucesor.
	log.Debugf("Sucesor actualizado al nod en %s\n.", suc.IP)

	// Actualiza el nuevo sucesor replicando las llaves de este nodo.
	go node.UpdateSuccessorKeys()
}

// PeriodicallyCheckSuccessor comprueba periodicamente si el sucesor a fallado.
func (node *Node) PeriodicallyCheckSuccessor() {
	log.Debug("El hilo para comprobar el sucesor a empezado.")

	ticker := time.NewTicker(1 * time.Second) // Establece el tiempo de comprobacion.
	for {
		select {
		case <-node.shutdown: // Si el servidor esta caido, parar el hilo.
			ticker.Stop()
			return
		case <-ticker.C: // Si es tiempo comprobar si el sucesor esta vivo.
			node.CheckSuccessor()
		}
	}
}

// FixFinger update a particular finger on the finger table, and return the index of the next finger to update.
func (node *Node) FixFinger(index int) int {
	log.Trace("Fixing finger entry.")
	defer log.Trace("Finger entry fixed.")

	m := node.config.HashSize            // Obtain the finger table size.
	ID := FingerID(node.ID, index, m)    // Obtain node.ID + 2^(next) mod(2^m).
	suc, err := node.FindIDSuccessor(ID) // Obtain the node that succeeds ID = node.ID + 2^(next) mod(2^m).
	// In case of error finding the successor, report the error and skip this finger.
	if err != nil || suc == nil {
		log.Errorf("Successor of ID not found.This finger fix was skipped.\n%s", err.Error())
		// Return the next index to fix.
		return (index + 1) % m
	}

	log.Tracef("Correspondent finger found at %s.", suc.IP)

	// If the successor of this ID is this node, then the ring has already been turned around.
	// Clean the remaining positions and return index 0 to restart the fixing cycle.
	if Equals(suc.ID, node.ID) {
		for i := index; i < m; i++ {
			node.fingerLock.Lock()    // Lock finger table to write on it, and unlock it after.
			node.fingerTable[i] = nil // Clean the correspondent position on the finger table.
			node.fingerLock.Unlock()
		}
		return 0
	}

	node.fingerLock.Lock()        // Lock finger table to write on it, and unlock it after.
	node.fingerTable[index] = suc // Update the correspondent position on the finger table.
	node.fingerLock.Unlock()

	// Return the next index to fix.
	return (index + 1) % m
}

// PeriodicallyFixFinger periodically fix finger table.
func (node *Node) PeriodicallyFixFinger() {
	log.Debug("Fix finger thread started.")

	next := 0                                        // Index of the actual finger entry to fix.
	ticker := time.NewTicker(100 * time.Millisecond) // Set the time between routine activations.
	for {
		select {
		case <-node.shutdown: // If node server is shutdown, stop the thread.
			ticker.Stop()
			return
		case <-ticker.C: // If it's time, fix the correspondent finger table entry.
			next = node.FixFinger(next)
		}
	}
}

// FixSuccessor fix an entry of the queue of successors.
// Given an entry of the successor queue, gets the reference to a remote node it contains and
// makes a remote call to GetSuccessor to get its successor.
// If the call fails, assume the remote node is dead, and remove this entry from the queue of successors.
// In this case, return the previous entry of this one, to fix this entry later.
// Otherwise, fix the next entry, updating its value with the obtained successor,
// and return the next entry of the queue.
func (node *Node) FixSuccessor(entry *QueueNode[chord.Node]) *QueueNode[chord.Node] {
	log.Trace("Fixing successor queue entry.")

	// If the queue node is null, report error.
	if entry == nil {
		log.Error("Error fixing successor queue entry: queue node argument cannot be null.")
		return nil
	}

	node.sucLock.RLock()                     // Lock the queue to read it, and unlock it after.
	value := entry.value                     // Obtain the successor contained in this queue node.
	prev := entry.prev                       // Obtain the previous node of this queue node.
	next := entry.next                       // Obtain the next node of this queue node.
	inside := entry.inside                   // Check if this queue node still being inside the queue.
	fulfilled := node.successors.Fulfilled() // Check if the queue is fulfilled.
	node.sucLock.RUnlock()

	// If the queue node is not inside the queue, return the next node.
	// If this queue node is the last one, and the queue is fulfilled, return null to restart the fixing cycle.
	if !inside || next == nil && fulfilled {
		log.Trace("Successor queue entry fixed.")
		return next
	}

	suc, err := node.RPC.GetSuccessor(value) // Otherwise, get the successor of this successor.
	// If there is an error, then assume this successor is dead.
	if err != nil {
		// If this successor is the immediate successor of this node, don't report the error,
		// to wait for CheckSuccessor to detect it and pop this node from the queue
		// (it's necessary for the correct transfer of keys).
		if prev == nil {
			// In this case, return the next node of this queue node, to skip this one and fix the remaining.
			return next
		} else {
			// Otherwise, report the error and remove the node from the queue of successors.
			log.Errorf("Error getting successor of successor at %s."+
				"Therefore is assumed dead and removed from the queue of successors.\n%s", value.IP, err.Error())

			node.sucLock.Lock()           // Lock the queue to write on it, and unlock it after.
			node.successors.Remove(entry) // Remove it from the queue of successors.
			// Push back this node, to ensure the queue is not empty.
			// Push back this node, to ensure the queue is not empty.
			if node.successors.Empty() {
				node.successors.PushBack(node.Node)
			}
			node.sucLock.Unlock()

			// In this case, return the previous node of this queue node, to fix this entry later.
			return prev
		}
	}

	node.sucLock.RLock()  // Lock the queue to read it, and unlock it after.
	next = entry.next     // Obtain the next node of this queue node.
	inside = entry.inside // Check if this queue node still being inside the queue.
	node.sucLock.RUnlock()

	// If the obtained successor is not this node, and is not the same node of this entry.
	if !Equals(suc.ID, node.ID) && !Equals(suc.ID, value.ID) {
		// If this queue node still on the queue.
		if inside {
			// If this queue node is the last one, push its successor at the end of queue.
			if next == nil {
				node.sucLock.Lock()           // Lock the queue to write on it, and unlock it after.
				node.successors.PushBack(suc) // Push this successor in the queue.
				node.sucLock.Unlock()
			} else {
				// Otherwise, fix next node of this queue node.
				node.sucLock.Lock() // Lock the queue to write on it, and unlock it after.
				next.value = suc    // Set this successor as value of the next node of this queue node.
				node.sucLock.Unlock()
			}
		} else {
			// Otherwise, skip this node and continue with the next one.
			return next
		}
	} else if Equals(suc.ID, value.ID) {
		// If the node is equal than its successor, skip this node and continue with the next one.
		return next
	} else {
		// Otherwise, if the obtained successor is this node, then the ring has already been turned around,
		// so there are no more successors to add to the queue.
		// Therefore, return null to restart the fixing cycle.
		return nil
	}

	log.Trace("Successor queue entry fixed.")
	return next
}

// PeriodicallyFixSuccessor periodically fix entries of the queue of successors.
func (node *Node) PeriodicallyFixSuccessor() {
	log.Debug("Fix successor thread started.")

	ticker := time.NewTicker(500 * time.Millisecond) // Set the time between routine activations.
	var entry *QueueNode[chord.Node] = nil           // Queue node entry for iterations.
	for {
		select {
		case <-node.shutdown: // If node server is shutdown, stop the thread.
			ticker.Stop()
			return
		case <-ticker.C: // If it's time, fix an entry of the queue.
			// Lock the successor to read it, and unlock it after.
			node.sucLock.RLock()
			suc := node.successors.Beg() // Obtain this node successor.
			node.sucLock.RUnlock()

			// If successor is not this node, then the queue of successors contains at least one successor.
			// Therefore, it needs to be fixed.
			if !Equals(suc.ID, node.ID) {
				// If actual queue node entry is null, restart the fixing cycle,
				// starting at the first queue node.
				if entry == nil {
					// Lock the successor to read it, and unlock it after.
					node.sucLock.RLock()
					entry = node.successors.first
					node.sucLock.RUnlock()
				}
				entry = node.FixSuccessor(entry) // Fix the corresponding entry.
			} else {
				entry = nil // Otherwise, reset the queue node entry.
			}
		}
	}
}

// FixStorage fix a particular key location on storage dictionary.
// To do this, it locates the node that corresponds to this key.
// If the node is this node, or its predecessor, then the key is correctly stored.
// Otherwise, the key is in a bad location, and therefore it's relocated and deleted from this node storage.
func (node *Node) FixStorage(key string) {
	// Lock the predecessor to read it, and unlock it after.
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	keyNode, err := node.LocateKey(key) // Locate the node that corresponds to the key.
	if err != nil {
		log.Errorf("Error fixing local storage dictionary.\n%s", err.Error())
		return
	}

	// If the obtained node is not this node, neither this node predecessor, then it is bad located.
	if !Equals(keyNode.ID, node.ID) && !Equals(keyNode.ID, pred.ID) {
		// Lock the storage dictionary to read it, and unlock it after.
		node.dictLock.RLock()
		value, err := node.dictionary.Get(key) // Get the value associated to the key.
		node.dictLock.RUnlock()
		if err != nil {
			// Don't report the error and return, because if the key is no longer in the dictionary
			// then it's simply no longer necessary to relocate it.
			return
		}

		// Lock the storage dictionary to write on it, and unlock it after.
		node.dictLock.Lock()
		err = node.dictionary.Delete(key) // Delete the key from local storage.
		node.dictLock.Unlock()
		if err != nil {
			log.Errorf("Error deleting key %s in local storage.\n%s", keyNode.IP, err.Error())
			return
		}

		// Set this <key, value> pair on the corresponding node.
		err = node.RPC.Set(keyNode, &chord.SetRequest{Key: key, Value: value})
		if err != nil {
			log.Errorf("Error relocating key %s to %s.\n%s", key, keyNode.IP, err.Error())
			// In case of error, reinsert the key on this node storage, to prevent the loss of information.
			// Lock the storage dictionary to write on it, and unlock it after.
			node.dictLock.Lock()
			err = node.dictionary.Set(key, value) // Reinsert the key on local storage.
			node.dictLock.Unlock()
			if err != nil {
				log.Errorf("Error reinserting key %s in local storage.\n%s", keyNode.IP, err.Error())
				return
			}
			return
		}
	}
}

// PeriodicallyFixStorage periodically fix storage dictionary.
func (node *Node) PeriodicallyFixStorage() {
	log.Debug("Fix storage thread started.")

	next := 0                                        // Index of the actual storage key to fix.
	keys := make([]string, 0)                        // Keys to fix.
	ticker := time.NewTicker(500 * time.Millisecond) // Set the time between routine activations.
	for {
		select {
		case <-node.shutdown: // If node server is shutdown, stop the thread.
			ticker.Stop()
			return
		case <-ticker.C: // If it's time, fix the correspondent storage entry.
			// If there is no more keys to fix.
			if next == len(keys) {
				// Lock the predecessor to read it, and unlock it after.
				node.predLock.RLock()
				pred := node.predecessor
				node.predLock.RUnlock()

				// Lock the storage dictionary to read it, and unlock it after
				node.dictLock.RLock()
				_, out, err := node.dictionary.Partition(pred.ID, node.ID) // Obtain the dictionary of replicated keys.
				node.dictLock.RUnlock()
				if err != nil {
					log.Errorf("Error fixing local storage dictionary. "+
						"Cannot obtain replicated keys on this node.\n%s", err.Error())
					continue
				}

				keys = Keys(out) // Obtain the replicated keys.
				next = 0         // Reset the index of the key to fix.
			}

			// If there are remaining keys to fix.
			if next < len(keys) {
				node.FixStorage(keys[next]) // Fix the corresponding key.
				next++                      // Update the index of the key to fix.
			}
		}
	}
}
