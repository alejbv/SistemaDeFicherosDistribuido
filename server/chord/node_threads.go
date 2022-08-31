package chord

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	log "github.com/sirupsen/logrus"
)

// Los hilos de ejecucion periodicos del nodo servidor.

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
		log.Errorf("Error estabilizando el nodo: no se puede obtener el predecesor del sucesor en %s.\n%s", suc.IP, err.Error())
		return
	}
	/*
		Si el nodo candidato es mas cercano a este nodo que su actual sucesor, se actualiza el
		sucesor de este nodo con el candidato
	*/
	if Equals(node.ID, suc.ID) || Between(candidate.ID, node.ID, suc.ID) {
		log.Debug("Sucesor actualizado al nodo en " + candidate.IP + ".")
		// Bloquea el sucesor para escribir en el, se desbloquea el finalizar
		node.sucLock.Lock()
		node.successors.PushBeg(candidate) //Se actualiza el sucesor de este nodo con el obtenido.
		suc = candidate
		node.sucLock.Unlock()
	}

	// Notifica al sucesor de la existencia de su predecesor.
	err = node.RPC.Notify(suc, node.Node)
	if err != nil {
		log.Errorf("Error notificando al sucesor en %s.\n%s", suc.IP, err.Error())
		return
	}

	log.Trace("Nodo estabilizado.")
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
			// Si existia un predecesor viejo, absorve sus llaves.
			//go node.AbsorbPredecessorKeys(pred)
			/*
				Posible modificacion
			*/
			intersectionFiles, _, _ := node.dictionary.PartitionFile(pred.ID, node.ID)
			go func() {

				node.AbsorbPredecessorKeys(pred)
				log.Debug("Se les va a notificar a las etiquetas que hubo un cambio en la ubicacion de unos archivos.")
				// Se recorre cada archivo
				for _, files := range intersectionFiles {
					tags := files.Tags

					// Se recorre cada etiqueta
					for _, tag := range tags {
						// Se Genera una request para editar las etiquetas
						enc := &chord.TagEncoder{
							FileName:      files.Name,
							FileExtension: files.Extension,
							NodeID:        pred.ID,
							NodeIP:        pred.IP,
							NodePort:      pred.Port,
						}

						req := &chord.EditFileFromTagRequest{Tag: tag, Mod: enc}

						go node.RPC.EditFileFromTag(node.Node, req)
					}

				}

			}()
			// Hasta aqui se extiende la modificacion
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
			log.Trace("El sucesor esta vivo.")
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
	log.Debugf("Sucesor actualizado al nodo en %s\n.", suc.IP)

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

// FixFinger actualiza una parte en particular de la FingerTable, y devuelve el index de la otro posicion a actualizar .
func (node *Node) FixFinger(index int) int {
	log.Trace("Arreglando entrada.")
	defer log.Trace("Entrada arreglada.")

	// Obtiene el tamaño de la FingerTable
	m := node.config.HashSize
	// Obtiene  node.ID + 2^(next) mod(2^m).
	ID := FingerID(node.ID, index, m)
	// Obtiene el nodo que  sucede a  ID = node.ID + 2^(next) mod(2^m).
	suc, err := node.FindIDSuccessor(ID)
	// En caso de error buscando el sucesor, reporta el error y se salta esta parte.
	if err != nil || suc == nil {
		log.Errorf("Sucesor de la ID no encontrado: el arreglo de esta parte fue saltado.\n%s", err.Error())
		//Devuelve el siguiente indice a arreglar.
		return (index + 1) % m
	}

	log.Tracef("Parte correspondiente encontrada en  %s.", suc.IP)

	// Si el sucesor de ID es este nodo, entonces el anillo esta listo.
	// Limpia las posiciones restantes y devuelve el indice 0 para reiniciar el ciclo arreglado.
	if Equals(suc.ID, node.ID) {
		for i := index; i < m; i++ {
			// Bloquea la FingerTable para escribir en ella, la desbloquea al terminar
			node.fingerLock.Lock()
			// CLimpia las posiciones correspondientes en la FingerTable
			node.fingerTable[i] = nil
			node.fingerLock.Unlock()
		}
		return 0
	}

	// Bloquea la FingerTable para escribir en ella, la desbloquea al terminar
	node.fingerLock.Lock()
	// Actualiza la posicion correspondiente en la FingerTable
	node.fingerTable[index] = suc
	node.fingerLock.Unlock()

	// Devuelve el proximo indice a arreglar
	return (index + 1) % m
}

// PeriodicallyFixFinger Periodicamente arregla la FingerTable.
func (node *Node) PeriodicallyFixFinger() {
	log.Debug("Hilo para arreglar la FingerTable empezado.")

	// Posicion de la actual posicion de la FingerTable a arreglar
	next := 0
	// Establece el tiempo entre la activacion de las rutinas
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		// Si el servidor esta caido termina el hilo
		case <-node.shutdown:
			ticker.Stop()
			return
		case <-ticker.C:
			// Si es tiempo arregla la posicion correspondiente en la FingerTable
			next = node.FixFinger(next)
		}
	}
}

/*
FixSuccessor arregla una entrada de la cola de los predecesores.
Dado una entrada de la cola de sucesores, obtiene las referencias a los nodos remotos
que contiene y hace una llamada remota a GetSuccessor para obtener sus sucesores.
Si la llamada falla, asume que el nodo remoto esta caido y remueve dicha entrada de la cola.
En otro caso arregla la siguiente entrada actualizando su valor con el sucesor obtenido

*/
func (node *Node) FixSuccessor(entry *QueueNode[chord.Node]) *QueueNode[chord.Node] {
	log.Trace("Arreglando la cola de sucesores introducida.")

	// Si es vacia devuelve error.
	if entry == nil {
		log.Error("Error arreglando la cola de sucesores introducida: no puede ser vacia.")
		return nil
	}

	// Bloquea  la cola para leer de ella, la desbloquea al terminar
	node.sucLock.RLock()
	// Obtiene el sucesor obtenido en la cola
	value := entry.value
	// Obtiene el nodo previo al de la entrada
	prev := entry.prev
	// Obtiene el nodo siguiente en la cola
	next := entry.next
	inside := entry.inside
	fulfilled := node.successors.Fulfilled()
	node.sucLock.RUnlock()

	/*
		Si el nodo no esta dentro de la cola, devuelve el nodo siguiente.
		Si es el ultimo y la cola esta completa devuelve un null para arreglar el ciclo
	*/
	if !inside || next == nil && fulfilled {
		log.Trace("Cola de sucesores arreglada.")
		return next
	}

	// En otro caso se obtiene el sucesor de este nodo
	suc, err := node.RPC.GetSuccessor(value)
	// Si hay un error se asume que el nodo esta muerto
	if err != nil {
		/*
			SI este sucesor es el sucesor inmediato del nodo, no se reporta el error,
			para esperar por  CheckSuccessor a que lo detecte y removerlo de la cola
		*/
		if prev == nil {
			// En este caso , devuevle el siguiente nodo en esta cola..
			return next
		} else {
			// En otro caso reporta en error y remuve el nodo de la cola
			log.Errorf("Error consiguiendo el sucesor del sucesor en  %s."+
				"Por lo tanto se asume que esta muerto y se remueve de la cola.\n%s", value.IP, err.Error())

			// Bloquea la cola para escribir en ella, desbloqueala al finalizar
			node.sucLock.Lock()
			//Se remueve de la cola.
			node.successors.Remove(entry)
			// Se agrega el propio nodo para asegurar que la cola no este vacia.
			if node.successors.Empty() {
				node.successors.PushBack(node.Node)
			}
			node.sucLock.Unlock()

			// En este caso, se devuelve el nodo previo al de entrada, para arreglarlo despues.
			return prev
		}
	}

	// Bloquea la cola para leer de ella, se desbloquea el finalizar
	node.sucLock.RLock()
	// Obtiene el nodo siguiente de la cola
	next = entry.next
	// Comprueba si este nodo sigue dentro de la cola
	inside = entry.inside
	node.sucLock.RUnlock()

	// Si el sucesor obtenido no es este nodo, y no es el mismo que el de entrada.
	if !Equals(suc.ID, node.ID) && !Equals(suc.ID, value.ID) {
		// Si esta todavia dentro de la cola
		if inside {
			// Si es el ultimo nodo se agregan sus sucesores al final de la cola.
			if next == nil {
				// Bloquea la cola para escribir en ella, se desbloquea al final
				node.sucLock.Lock()
				// Agrega al final de la cola este sucesor
				node.successors.PushBack(suc)
				node.sucLock.Unlock()
			} else {
				// En otro caso arregla el siguiente nodo de la cola.
				// Bloquea la cola para escribir en ella, se desbloquea al final
				node.sucLock.Lock()
				// Establece el  sucesor como el siguiente nodo de la cola
				next.value = suc
				node.sucLock.Unlock()
			}
		} else {
			// En otro caso, se salta este nodo y sigue con el siguiente
			return next
		}
	} else if Equals(suc.ID, value.ID) {
		// Si el nodo es igual a su sucesor, salta este nodo y sigue con el proximo
		return next
	} else {
		/*
		 En el caso de que el sucesor obtenido sea este nodo, entoces el anillo a sido
		 invertido por lo que no hay mas sucesor que agregar a la cola. Por tanto se
		 devuelve null para reiniciar el arreglo del ciclo
		*/
		return nil
	}

	log.Trace("Cola de sucesor arreglada.")
	return next
}

// PeriodicallyFixSuccessor Periodicamente arregla la cola de sucesor.
func (node *Node) PeriodicallyFixSuccessor() {
	log.Debug("Hilo para arreglar la cola de sucesor arreglado.")

	// Establece el tiempo entre la activacion de las rutinas
	ticker := time.NewTicker(500 * time.Millisecond)
	//Se va a iterar por la cola de sucesores
	var entry *QueueNode[chord.Node] = nil
	for {
		select {
		case <-node.shutdown: // SI el nodo servidor esta caido se termina el hilo
			ticker.Stop()
			return
			//Si es tiempo arregla una entrada de la cola
		case <-ticker.C:
			// Bloquea el sucesor para leer de el, se desbloquea al carg
			node.sucLock.RLock()
			// Se obtiene el sucesor de este nodo
			suc := node.successors.Beg()
			node.sucLock.RUnlock()

			// Si el sucesor es distinto del nodo, entonces la cola posee al menos un nodo
			// Por tanto debe ser arreglado
			if !Equals(suc.ID, node.ID) {
				// Si el nodo actual es vacio entonces el ciclod debe ser arreglado,
				// empezando por el primer nodo de la query
				if entry == nil {
					// Bloquea el sucesor para leer de el, se desbloquea al terminar
					node.sucLock.RLock()
					entry = node.successors.first
					node.sucLock.RUnlock()
				}
				// Arregla la entrada correspondiente
				entry = node.FixSuccessor(entry)
			} else {
				// En otro caso se resetea
				entry = nil
			}
		}
	}
}

// Posible Metodo a modificar
/*
FixStorage arregla la localizacon de  una llave en particular dentro del
diccionario de almacenamiento. Para eso localiza el nodo  correspondiente. Si
el nodo es este nodo o su predecesor, entonces no hay problema.
En otro caso esta mal localizada y por tanto debe ser reubicada y eliminada
*/
// Este es un metodo a tener en cuenta
func (node *Node) FixTagStorage(key string) {

	//Bloquea el predecesor para leer de el, lo desbloquea al terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	// Localiza el nodo al que le corresponde la clave
	keyNode, err := node.LocateKey(key)
	if err != nil {
		log.Errorf("Error arreglando las etiquetas en el almacenamiento local .\n%s", err.Error())
		return
	}

	// Si el nodo obtenido es distinto a el nodo actual o su predecesor, entonces esta mal ubicado
	if !Equals(keyNode.ID, node.ID) && !Equals(keyNode.ID, pred.ID) {
		// Bloquea el diccionario para leer de el, lo desbloquea al terminar
		node.dictLock.RLock()
		// Consigue el valor asociado a esta llave
		value, err := node.dictionary.GetTag(key)
		node.dictLock.RUnlock()
		if err != nil {
			/*
				En este caso no se reporta el error, porque la clave ya no esta aqui.Por lo
				que no es necesario reubicarlo.
			*/
			return
		}

		encod, _ := json.MarshalIndent(value, "", " ")
		var mapa map[string][]byte
		mapa[key] = encod

		// Bloquea el diccionario para escribir en el, lo desbloquea al terminar
		node.dictLock.Lock()
		// Elimina la llave del almacenamiento local
		err = node.dictionary.DeleteTag(key)
		node.dictLock.Unlock()
		if err != nil {
			log.Errorf("Error eliminando la etiqueta %s en el almacenamiento local.\n%s", key, err.Error())
			return
		}

		// Establece este par <key, value> en el nodo correspondiente
		err = node.RPC.Extend(keyNode, &chord.ExtendRequest{Tags: mapa, Files: nil})
		if err != nil {
			log.Errorf("Error reubicando la llave  %s a %s.\n%s", key, keyNode.IP, err.Error())
			/*
				En caso de error, reinserta la llave en este nodo, para prevenir
				la perdidad de informacion.
			*/
			// Bloquea el diccionario para escribir en el, desbloquea al terminar
			node.dictLock.Lock()
			// Reinserta la clave en el almacenamiento local
			err = node.dictionary.ExtendTags(mapa)
			node.dictLock.Unlock()
			if err != nil {
				log.Errorf("Error reinsertando la llave %s en el almacenamiento local .\n%s", keyNode.IP, err.Error())
				return
			}
			return
		}
	}
}

/*
FixStorage arregla la localizacon de  una llave en particular dentro del
diccionario de almacenamiento. Para eso localiza el nodo  correspondiente. Si
el nodo es este nodo o su predecesor, entonces no hay problema.
En otro caso esta mal localizada y por tanto debe ser reubicada y eliminada
*/
// Este es un metodo a tener en cuenta
func (node *Node) FixFileStorage(key string) {
	//Bloquea el predecesor para leer de el, lo desbloquea al terminar
	chain := strings.Split(key, "-")
	fileName := strings.Join(chain[:len(chain)-1], "-")
	fileExtension := chain[len(chain)-1]

	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	// Localiza el nodo al que le corresponde la clave
	keyNode, err := node.LocateKey(fileName)
	if err != nil {
		log.Errorf("Error arreglando el almacenamiento local: No se pudo localizar el nodo adecuado .\n%s", err.Error())
		return
	}

	// Si el nodo obtenido es distinto a el nodo actual o su predecesor, entonces esta mal ubicado
	if !Equals(keyNode.ID, node.ID) && !Equals(keyNode.ID, pred.ID) {
		// Bloquea el diccionario para leer de el, lo desbloquea al terminar
		node.dictLock.RLock()
		// Consigue el valor asociado a esta llave
		value, tags, err := node.dictionary.GetFile(fileName, fileExtension)
		node.dictLock.RUnlock()
		if err != nil {
			/*
				En este caso no se reporta el error, porque la clave ya no esta aqui.Por lo
				que no es necesario reubicarlo.
			*/
			return
		}

		// Bloquea el diccionario para escribir en el, lo desbloquea al terminar
		node.dictLock.Lock()
		// Elimina la llave del almacenamiento local
		err = node.dictionary.DeleteFile(fileName, fileExtension)
		node.dictLock.Unlock()
		if err != nil {
			log.Errorf("Error eliminando el fichero %s en el almacenamiento local.\n%s", fileName, err.Error())
			return
		}
		tagFile := &chord.TagFile{
			Name:      fileName,
			Extension: fileExtension,
			File:      value,
			Tags:      tags,
		}

		// Establece este par <key, value> en el nodo correspondiente
		err = node.RPC.Extend(keyNode, &chord.ExtendRequest{Tags: nil, Files: []*chord.TagFile{tagFile}})
		if err != nil {
			log.Errorf("Error reubicando la llave  %s a %s.\n%s", key, keyNode.IP, err.Error())
			/*
				En caso de error, reinserta la llave en este nodo, para prevenir
				la perdidad de informacion.
			*/
			// Bloquea el diccionario para escribir en el, desbloquea al terminar
			node.dictLock.Lock()
			// Reinserta la clave en el almacenamiento local
			err = node.dictionary.SetFile(tagFile)
			node.dictLock.Unlock()
			if err != nil {
				log.Errorf("Error reinsertando la llave %s en el almacenamiento local .\n%s", keyNode.IP, err.Error())
				return
			}
			return
		}

		/*
			Como esto se llama para arreglar la replicacion, pero se deja comentado por si acaso
			// Una vez que el archivo se movio adecuadamente paso a modificar el cambio que hubo
			for _, tag := range tags {
				// Se crea un nuevo objeto TagEncoder que va a tener la informacion a modificar
				encoder := &chord.TagEncoder{
					FileName:      fileName,
					FileExtension: fileExtension,
					NodeID:        keyNode.ID,
					NodeIP:        keyNode.IP,
					NodePort:      keyNode.Port,
				}
				// Se genera la request
				req := &chord.EditFileFromTagRequest{Tag: tag, Mod: encoder}
				// Por cada una de las etiquetas se modificia
				go node.RPC.EditFileFromTag(keyNode, req)
			}
		*/
	}
}

// Posible Metodo a modificar
// PeriodicallyFixStorage periodicamente arregla el almcenamiento local.
func (node *Node) PeriodicallyFixStorage() {
	log.Debug("Hilo para arreglar el almacenamiento empezado.")
	// Indice de la llave actual para arreglar.
	nextFile := 0
	nextTag := 0
	// Llaves a arreglar
	var keysTags []string
	var keysFiles []string
	// Establece el tiempo entre reactivacion de las rutinas
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		// Si el nodo servidor esta apagado, entonces termina el hilo
		case <-node.shutdown:
			ticker.Stop()
			return
			// Si es tiempo arregla el almacenamiento actual
		case <-ticker.C:
			// Si no hay mas claves para arreglar
			if nextFile == len(keysFiles) && nextTag == len(keysTags) {
				// Bloquea el predecesor para leer de el, desbloquealo al terminar
				node.predLock.RLock()
				pred := node.predecessor
				node.predLock.RUnlock()

				// Bloquea el diccionario para leer de el, desbloquealo al terminar
				node.dictLock.RLock()
				// Obtener el diccionario de las llaves replicadas
				_, outTags, err := node.dictionary.PartitionTag(pred.ID, node.ID)
				node.dictLock.RUnlock()
				if err != nil {
					log.Errorf("Error arreglando el diccionario de almacenamiento local: "+
						"No se puede obtener las etiquetas replicadas de este nodo.\n%s", err.Error())
					continue
				}
				_, outFiles, err := node.dictionary.PartitionFile(pred.ID, node.ID)
				node.dictLock.RUnlock()
				if err != nil {
					log.Errorf("Error arreglando el diccionario de almacenamiento local: "+
						"No se puede obtener los archivos replicados de este nodo.\n%s", err.Error())
					continue
				}

				// Obtener las etiquetas  replicadas
				keysTags = Keys(outTags)

				// Obtener los archivos replicados

				for _, file := range outFiles {
					keysFiles = append(keysFiles, file.Name+"-"+file.Extension)
				}
				// Restablece el indice de la llave a arreglar.

				nextFile = 0
				nextTag = 0
			}

			// Si quedan etiquetas por arreglar
			if nextTag < len(keysTags) {
				// Arregla la etiqueta correspondiente
				node.FixTagStorage(keysTags[nextTag])
				// Actualiza el indicie de la etiqueta a arreglar
				nextTag++
			}

			if nextFile < len(keysFiles) {
				// Arregla el archvio correspondiente
				node.FixFileStorage(keysFiles[nextFile])
				// Actualiza el indicie del archivo a arreglar
				nextFile++
			}
		}
	}
}
