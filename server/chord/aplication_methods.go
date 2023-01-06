package chord

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alejbv/SistemaDeFicherosDistribuido/chord"
	log "github.com/sirupsen/logrus"
)

// AddFile almacena  un fichero en el almacenamiento local .
func (node *Node) AddFile(ctx context.Context, req *chord.AddFileRequest) (*chord.AddFileResponse, error) {

	file := req.File

	log.Infof("Almacenar fichero %s.", file.GetName())

	// Si la request es una replica se resuelve de forma local.
	if req.Replica {
		log.Info("Resolviendo la request de forma local (replicacion).")

		// Bloquea el diccionario para escribir en el, se desbloquea el terminar
		node.dictLock.Lock()
		// Almacena el fichero en el almacenamiento
		//err := node.dictionary.SetWithLock(file, address)
		err := node.dictionary.SetFile(file)
		node.dictLock.Unlock()

		if err != nil {
			log.Error("Error almacenado el fichero.")
			return &chord.AddFileResponse{}, errors.New("error almacenado la llave.\n" + err.Error())

		} else if err == os.ErrPermission {
			log.Error("Error almacenado la llave: ya esta bloqueada.\n" + err.Error())
			return &chord.AddFileResponse{}, err
		}

		log.Info("Se almaceno de forma exitosa.")
		//new_node := chord.Node{ID: node.ID, IP: node.IP, Port: node.Port}
		return &chord.AddFileResponse{Destine: node.Node}, nil
	}

	// Por defecto,toma este nodo para almacenar el archivo en el almacenamiento local.
	keyNode := node.Node
	// Bloquea el predecesor para leer de el, se desbloquea el terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	/*
		Si la ID de la llave no esta entre la ID de este nodo y su predecesor
		entonces la request no es necesariamente local
	*/
	if between, err := KeyBetween(file.Name, node.config.Hash, pred.ID, node.ID); !between && err == nil {
		log.Info("Buscando por el nodo correspondiente.")
		// Localiza el nodo que corresponde a esta llave

		keyNode, err = node.LocateKey(file.Name)
		if err != nil {
			log.Error("Error estableciendo la llave.")
			return &chord.AddFileResponse{}, errors.New("Error estableciendo la llave.\n" + err.Error())
		}

	} else if err != nil {
		log.Error("Error estableciendo la llave.")
		return &chord.AddFileResponse{}, errors.New("Error estableciendo la llave.\n" + err.Error())
	}

	// Si la llave corresponde a este nodo, directamente almacena el fichero de forma local.
	if Equals(keyNode.ID, node.ID) {
		log.Info("Resolviendo la request de forma local.")

		// Bloquea el diccionario para escribir en el, se desbloquea el terminar
		node.dictLock.Lock()
		//Almacena el fichero de forma local .
		err := node.dictionary.SetFile(file)
		node.dictLock.Unlock()
		if err != nil {
			log.Error("Error almacenando el fichero.")
			return &chord.AddFileResponse{}, errors.New("error almacenando el fichero.\n" + err.Error())

		} else if err == os.ErrPermission {
			log.Error("Error almacenando la llave: ya está bloqueada.\n" + err.Error())
			return &chord.AddFileResponse{}, err
		}

		log.Info("Resolucion exitosa.")

		// Bloquea el sucesor para leer de el, se desbloquea el terminar
		node.sucLock.RLock()
		suc := node.successors.Beg()
		node.sucLock.RUnlock()

		/*
			Antes de terminar se recorren las etiquetas y por cada una se almacena la informacion
			de donde se almaceno el fichero actual
		*/
		tags := file.Tags
		for _, tag := range tags {
			// Encuentro el nodo donde se debe almacenar la etiqueta
			temp_node, err := node.LocateKey(tag)
			// Si hubo algun error encontrando dicho nodo
			if err != nil {
				log.Error("Error localizando el nodo correspondiente a la etiqueta %s\n.", tag)
				return &chord.AddFileResponse{}, fmt.Errorf(fmt.Sprintf("Error localizando el nodo correspondiente a la etiqueta %s\n.", tag))
			}
			tagRequest := &chord.AddTagRequest{
				Tag:          tag,
				FileName:     file.Name,
				ExtensioName: file.Extension,
				TargetNode:   keyNode,
				Replica:      false,
			}
			// Si hubo un error almacenando alguna etiqueta
			err = node.RPC.AddTag(temp_node, tagRequest)
			if err != nil {
				log.Error("Error almacenando la etiqueta %s\n.", tag)
				return &chord.AddFileResponse{}, fmt.Errorf(fmt.Sprintf("Error almacenando la etiqueta %s\n.", tag))
			}
		}

		//Si el sucesor no es este nodo, replica la request a el.
		if !Equals(suc.ID, node.ID) {
			go func() {
				req.Replica = true
				log.Infof("Replicando la request set a %s.", suc.IP)
				err := node.RPC.AddFile(suc, req)
				if err != nil {
					log.Errorf("Error replicando la request a %s.\n%s", suc.IP, err.Error())
				}
			}()
		}

		return &chord.AddFileResponse{Destine: node.Node}, nil
	} else {
		log.Infof("Redirigiendo la request a %s.", keyNode.IP)
		// En otro caso, devuelve el resultado de la llamada remota en el nodo correspondiente.
		return &chord.AddFileResponse{Destine: keyNode}, node.RPC.AddFile(keyNode, req)
	}
}
func (node *Node) DeleteFileByQuery(ctx context.Context, req *chord.DeleteFileByQueryRequest) (*chord.DeleteFileByQueryResponse, error) {

	// Objetos para poder llevar un registro de los ficheros, sabiendo en cuantas querys está presente
	querys, target, err := node.QuerySolver(req.Tag)

	if err != nil {
		log.Errorf("Error mientras se realizaban las de las etiquetas:")
		return &chord.DeleteFileByQueryResponse{}, err
	}

	// Se bloquea el sucesor para poder leer de el, se desbloquea al terminar
	node.sucLock.RLock()
	suc := node.successors.Beg()
	node.sucLock.RUnlock()

	// Se bloquea el predecesor para poder leer de el, se desbloquea al terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()
	for key, value := range querys {
		// Esto representa la interseccion de todas las querys, o sea se procesa si las cumple todas
		if len(value) == len(req.Tag) {
			// Separa la clave en el nombre del fichero y su extension
			idef := strings.Split(key, ".")
			//Nombre del fichero
			name := idef[0]
			// Nombre de la extension
			extension := idef[1]

			Newreq := &chord.DeleteFileRequest{
				FileName:      name,
				FileExtension: extension,
				Replica:       false,
			}
			// Llama a eliminar el fichero localmente
			if Equals(target[key].ID, node.ID) {

				//Bloquea el diccionario para escribir en el, se desbloquea el final
				node.dictLock.Lock()
				err := node.dictionary.DeleteFile(name, extension)
				node.dictLock.Unlock()

				if err != nil {
					return &chord.DeleteFileByQueryResponse{}, err
				}
				if !Equals(suc.ID, node.ID) {
					go func() {
						Newreq.Replica = true
						log.Infof("Replicando la request a %s.", suc.IP)
						err := node.RPC.DeleteFile(suc, Newreq)
						if err != nil {
							log.Errorf("Error replicando la request a %s.\n%s", suc.IP, err.Error())
						}
					}()
				}

			} else {
				// Se elimina el fichero remotamente
				err := node.RPC.DeleteFile(target[key], Newreq)
				if err != nil {
					return &chord.DeleteFileByQueryResponse{}, err
				}

			}
			for _, tag := range value {
				keyNode := target[key]
				DeleteReq := &chord.DeleteFileFromTagRequest{
					Tag:           tag,
					FileName:      name,
					FileExtension: extension,
					Replica:       false,
				}

				if between, err := KeyBetween(tag, node.config.Hash, pred.ID, node.ID); between && err == nil {
					log.Debug("Buscando por el nodo correspondiente.")
					//La informacion esta guardada local
					keyNode = node.Node

				} else if err != nil {
					log.Error("Error encontrando el nodo.")
					return &chord.DeleteFileByQueryResponse{}, errors.New("error encontrando el nodo.\n" + err.Error())
				}
				if Equals(keyNode.ID, node.ID) {
					// Se elimina la informacion local

					//Bloquea el diccionario para escribir en el, se desbloquea el final
					node.dictLock.Lock()
					err := node.dictionary.DeleteFileFromTag(tag, name, extension)
					node.dictLock.Unlock()

					if err != nil {
						log.Errorf("Error eliminando  la informacion del archivo %s de la etiqueta %s\n", name, err.Error())
					}
					if !Equals(suc.ID, node.ID) {
						go func() {
							DeleteReq.Replica = true
							log.Debugf("Replicando la request a %s.", suc.IP)
							err := node.RPC.DeleteFileFromTag(suc, DeleteReq)
							if err != nil {
								log.Errorf("Error replicando la request a %s.\n%s", suc.IP, err.Error())
							}
						}()

					}
				} else {
					// Si se tiene que hacer una request remota
					go func() {
						node.RPC.DeleteFileFromTag(target[key], DeleteReq)
					}()
				}

			}
		}
	}
	return &chord.DeleteFileByQueryResponse{}, nil
}
func (node *Node) ListByQuery(ctx context.Context, req *chord.ListByQueryRequest) (*chord.ListByQueryResponse, error) {

	// Objeto poder llevar un registro de los ficheros, sabiendo en cuantas querys está presente
	querys, target, err := node.QuerySolver(req.Tags)

	if err != nil {
		log.Errorf("Error mientras se realizaban las de las etiquetas:")
		return &chord.ListByQueryResponse{}, err
	}
	response := make(map[string][]byte)

	log.Info("Terminado QuerySolver: Procesar info")
	for key, value := range querys {
		// Esto representa la interseccion de todas las querys, o sea se procesa si
		// las cumple todas
		if len(value) == len(req.Tags) {
			// Separa la clave en el nombre del fichero y su extension
			idef := strings.Split(key, ".")
			//Nombre del fichero
			name := idef[0]
			// Nombre de la extension
			extension := idef[1]

			req := &chord.GetFileInfoRequest{
				FileName:      name,
				FileExtension: extension,
			}
			// Obtener el archivo
			if Equals(target[key].ID, node.ID) {

				//Bloquea el diccionario para escribir en el, se desbloquea el final
				log.Infof("La informacion del archivo %s es local", name)
				node.dictLock.Lock()
				tempValue, err := node.dictionary.GetFileInfo(name, extension)
				node.dictLock.Unlock()

				if err != nil {
					log.Errorf("Hubo un error recuperando la informacion del archivo %s: %s", name, err.Error())
					return &chord.ListByQueryResponse{}, err
				}
				response[key] = tempValue

			} else {
				log.Infof("La informacion del archivo %s no es local: Buscando", name)
				resp, err := node.RPC.GetFile(target[key], req)
				if err != nil {
					log.Errorf("Hubo un error recuperando la informacion del archivo %s: %s", name, err.Error())
					return &chord.ListByQueryResponse{}, err
				}
				response[key] = resp.Info
			}

		}
	}
	return &chord.ListByQueryResponse{Response: response}, nil
}
func (node *Node) AddTagsByQuery(ctx context.Context, req *chord.AddTagsByQueryRequest) (*chord.AddTagsByQueryResponse, error) {

	// Objeto poder llevar un registro de los ficheros, sabiendo en cuantas querys está presente

	querys, target, err := node.QuerySolver(req.QueryTags)

	if err != nil {
		log.Errorf("Error mientras se realizaban las de las etiquetas:")
		return &chord.AddTagsByQueryResponse{}, err
	}

	// Se bloquea el sucesor para poder leer de el, se desbloquea al terminar
	node.sucLock.RLock()
	suc := node.successors.Beg()
	node.sucLock.RUnlock()

	// Se bloquea el predecesor para poder leer de el, se desbloquea al terminar
	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	for key, value := range querys {
		// Esto representa la interseccion de todas las querys, o sea se procesa si
		// las cumple todas
		if len(value) == len(req.QueryTags) {
			// Separa la clave en el nombre del fichero y su extension
			idef := strings.Split(key, ".")
			//Nombre del fichero
			name := idef[0]
			// Nombre de la extension
			extension := idef[1]
			RemoteReq := &chord.AddTagsToFileRequest{
				FileName:      name,
				FileExtension: extension,
				ListTags:      req.AddTags,
				Replica:       false,
			}
			if Equals(target[key].ID, node.ID) {

				//Bloquea el diccionario para escribir en el, se desbloquea el final
				node.dictLock.Lock()
				_, err := node.dictionary.AddTagsToFile(name, extension, req.AddTags)
				node.dictLock.Unlock()

				if err != nil {
					log.Errorf("Error mientras se agregaban etiquetas a los archivos locales:")
					return &chord.AddTagsByQueryResponse{}, err
				}

				// Si debe replicarle la request al sucesor
				if !Equals(suc.ID, node.ID) {
					go func() {
						RemoteReq.Replica = true
						log.Debugf("Replicando la request  a %s.", suc.IP)
						_, err := node.RPC.AddTagsToFile(suc, RemoteReq)
						if err != nil {
							log.Errorf("Error replicando la request a %s.\n%s", suc.IP, err.Error())
						}
					}()
				}

			} else {
				// Se llama remoto para agregar las nuevas etiquetas al fichero
				_, err := node.RPC.AddTagsToFile(target[key], RemoteReq)
				if err != nil {
					return &chord.AddTagsByQueryResponse{}, err
				}
			}
			for _, tag := range req.QueryTags {
				keyNode := target[key]
				AddRequest := &chord.AddTagRequest{
					Tag:          tag,
					FileName:     name,
					ExtensioName: extension,
					TargetNode:   target[key],
					Replica:      false,
				}

				if between, err := KeyBetween(tag, node.config.Hash, pred.ID, node.ID); between && err == nil {
					log.Debug("Buscando por el nodo correspondiente.")
					//La informacion esta guardada local
					keyNode = node.Node

				} else if err != nil {
					log.Error("Error encontrando el nodo.")
					return &chord.AddTagsByQueryResponse{}, errors.New("error encontrando el nodo.\n" + err.Error())
				}

				if Equals(keyNode.ID, node.ID) {
					// Se elimina la informacion local

					//Bloquea el diccionario para escribir en el, se desbloquea el final
					node.dictLock.Lock()
					err := node.dictionary.SetTag(tag, name, extension, target[key])
					node.dictLock.Unlock()

					if err != nil {
						log.Errorf("Error eliminando  la informacion del archivo %s de la etiqueta %s\n", name, err.Error())
					}
					if !Equals(suc.ID, node.ID) {
						go func() {
							AddRequest.Replica = true
							log.Debugf("Replicando la request a %s.", suc.IP)
							err := node.RPC.AddTag(suc, AddRequest)
							if err != nil {
								log.Errorf("Error replicando la request a %s.\n%s", suc.IP, err.Error())
							}
						}()

					}
				} else {
					// Si se tiene que hacer una request remota
					go func() {
						node.RPC.AddTag(target[key], AddRequest)
					}()
				}
			}
		}
	}
	return &chord.AddTagsByQueryResponse{}, nil
}
func (node *Node) DeleteTagsByQuery(ctx context.Context, req *chord.DeleteTagsByQueryRequest) (*chord.DeleteTagsByQueryResponse, error) {
	// Objeto poder llevar un registro de los ficheros, sabiendo en cuantas querys está presente

	querys, target, err := node.QuerySolver(req.QueryTags)
	if err != nil {
		log.Errorf("Error mientras se realizaban las de las etiquetas:")
		return &chord.DeleteTagsByQueryResponse{}, err
	}

	// Se bloquea el sucesor para poder leer de el, se desbloquea al terminar

	node.sucLock.RLock()
	suc := node.successors.Beg()
	node.sucLock.RUnlock()

	node.predLock.RLock()
	pred := node.predecessor
	node.predLock.RUnlock()

	for key, value := range querys {
		// Esto representa la interseccion de todas las querys, o sea se procesa si
		// las cumple todas
		if len(value) == len(req.QueryTags) {
			// Separa la clave en el nombre del fichero y su extension
			idef := strings.Split(key, ".")
			//Nombre del fichero
			name := idef[0]
			// Nombre de la extension
			extension := idef[1]

			RemoteReq := &chord.DeleteTagsFromFileRequest{
				FileName:      name,
				FileExtension: extension,
				ListTags:      req.RemoveTags,
				Replica:       false,
			}
			if Equals(target[key].ID, node.ID) {

				//Bloquea el diccionario para escribir en el, se desbloquea el final
				node.dictLock.Lock()
				_, err := node.dictionary.DeleteTagsFromFile(name, extension, req.RemoveTags)
				node.dictLock.Unlock()

				if err != nil {
					log.Errorf("Error mientras se eliminaban la informacion de los archivos de las etiquetas locales:")
					return &chord.DeleteTagsByQueryResponse{}, err
				}

				// Si debe replicarle la request al sucesor
				if !Equals(suc.ID, node.ID) {
					go func() {
						RemoteReq.Replica = true
						log.Debugf("Replicando la request  a %s.", suc.IP)
						_, err := node.RPC.DeleteTagsFromFile(suc, RemoteReq)
						if err != nil {
							log.Errorf("Error replicando la request a %s.\n%s", suc.IP, err.Error())
						}
					}()
				}

			} else {
				// Llama a eliminar el fichero
				_, err := node.RPC.DeleteTagsFromFile(target[key], RemoteReq)
				if err != nil {
					return &chord.DeleteTagsByQueryResponse{}, err
				}
			}

			for _, tag := range req.QueryTags {
				keyNode := target[key]
				DeleteRquest := &chord.DeleteFileFromTagRequest{
					Tag:           tag,
					FileName:      name,
					FileExtension: extension,
					Replica:       false,
				}

				if between, err := KeyBetween(tag, node.config.Hash, pred.ID, node.ID); between && err == nil {
					log.Debug("Buscando por el nodo correspondiente.")
					//La informacion esta guardada local
					keyNode = node.Node

				} else if err != nil {
					log.Error("Error encontrando el nodo.")
					return &chord.DeleteTagsByQueryResponse{}, errors.New("error encontrando el nodo.\n" + err.Error())
				}

				if Equals(keyNode.ID, node.ID) {
					// Se elimina la informacion local

					//Bloquea el diccionario para escribir en el, se desbloquea el final
					node.dictLock.Lock()
					err := node.dictionary.DeleteFileFromTag(tag, name, extension)
					node.dictLock.Unlock()

					if err != nil {
						log.Errorf("Error eliminando  la informacion del archivo %s de la etiqueta %s\n", name, err.Error())
					}
					if !Equals(suc.ID, node.ID) {
						go func() {
							DeleteRquest.Replica = true
							log.Debugf("Replicando la request a %s.", suc.IP)
							err := node.RPC.DeleteFileFromTag(suc, DeleteRquest)
							if err != nil {
								log.Errorf("Error replicando la request a %s.\n%s", suc.IP, err.Error())
							}
						}()

					}
				} else {
					// Si se tiene que hacer una request remota
					go func() {
						node.RPC.DeleteFileFromTag(target[key], DeleteRquest)
					}()
				}

			}
		}

	}
	return &chord.DeleteTagsByQueryResponse{}, nil
}
func (node *Node) QuerySolver(Tags []string) (map[string][]string, map[string]*chord.Node, error) {
	// Objeto poder llevar un registro de los ficheros, sabiendo en cuantas querys está presente
	querys := make(map[string][]string)
	target := make(map[string]*chord.Node)
	/*
		Lo primero que se necesita es poder obtener toda la informacion de cada uno de las etiquetas.
		Para eso debe haber una parte que por cada etiqueta pida todos los ficheros que tiene
	*/
	for _, tag := range Tags {

		// Por defecto, se toma este nodo.
		keyNode := node.Node
		// Bloquea el predecesor para poder leer de el, se desbloquea al terminar
		node.predLock.RLock()
		pred := node.predecessor
		node.predLock.RUnlock()
		// Encuentra en donde esta ubicada dicha etiqueta
		if between, err := KeyBetween(tag, node.config.Hash, pred.ID, node.ID); !between && err == nil {
			log.Info("Buscando por el nodo correspondiente.")
			// Localiza el nodo que almacena el fichero
			keyNode, err = node.LocateKey(tag)
			if err != nil {
				log.Error("Error encontrando el nodo remoto.")
				return nil, nil, errors.New("error encontrando el nodo remoto.\n" + err.Error())
			}

		} else if err != nil {
			log.Error("Error comprobando si es el nodo correspondiente.")
			return nil, nil, errors.New("error comprobando si es el nodo correspondiente.\n" + err.Error())
		}

		// Si pasa por aqui es q ya tiene nodo al que buscar
		// Si el fichero corresponde a este nodo, se busca directamente del almacenamiento.
		if Equals(keyNode.ID, node.ID) {
			log.Info("Resolviendo la request de forma local.")

			// Bloquea el diccionario para escribir en el y se desbloquea al terminar la funcion.
			node.dictLock.Lock()
			// obteniendo la informacion de la etiqueta .

			values, err := node.dictionary.GetTag(tag)
			node.dictLock.Unlock()
			if err != nil && err != os.ErrPermission {
				log.Error("Error recuperando la informacion de las etiquetas.")
				return nil, nil, errors.New("error recuperando la informacion de las etiquetas.\n" + err.Error())

			} else if err == os.ErrPermission {
				log.Error("Error recuperando la informacion de las etiquetas: esta bloqueada.\n" + err.Error())
				return nil, nil, err
			}

			log.Info("Recuperacion exitosa.")
			// Aqui vendria lo que se hace una vez con la informacion de las tags
			for _, value := range values {
				file := value.FileName + "." + value.FileExtension
				tempNode := &chord.Node{ID: value.NodeID, IP: value.NodeIP, Port: value.NodePort}

				if list, ok := querys[file]; ok {
					list = append(list, tag)
					querys[file] = list
				} else {
					querys[file] = []string{tag}
					target[file] = tempNode
				}
			}

		} else {
			log.Infof("Redirigiendo la request  para %s.", keyNode.IP)
			// En otro caso, se devuelve el resultado de la llamada remota en el nodo correspondiente

			getTagRequest := &chord.GetTagRequest{Tag: tag}
			res, err := node.RPC.GetTag(keyNode, getTagRequest)
			if err != nil {

				log.Error("Error al recibir los TagEncoders")
				return nil, nil, errors.New("Error al recibr los TagEncoders\n" + err.Error())
			}
			// Se tiene la respuesta y se esta trabajando con ella
			for _, value := range res {
				file := value.FileName + "." + value.FileExtension
				tempNode := &chord.Node{ID: value.NodeID, IP: value.NodeIP, Port: value.NodePort}

				if list, ok := querys[file]; ok {
					list = append(list, tag)
					querys[file] = list
				} else {
					querys[file] = []string{tag}
					target[file] = tempNode
				}
			}
			/*
				Hasta aqui esta la implementacion general, que se deberá usar en los otros metodos
			*/
		}

	}
	log.Infof("Se obtuvo toda la informacion de las etiquetas %s", Tags)
	return querys, target, nil
}
