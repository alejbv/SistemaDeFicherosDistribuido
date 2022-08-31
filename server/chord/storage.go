package chord

import (
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storage interface {
	GetFile(string, string) ([]byte, []string, error)
	SetFile(*chord.TagFile) error
	DeleteFile(string, string) error
	PartitionFile([]byte, []byte) ([]*chord.TagFile, []*chord.TagFile, error)
	ExtendFiles([]*chord.TagFile) error
	DiscardFiles([]string) error

	SetTag(string, string, string, *chord.Node) error
	GetTag(string) ([]TagEncoding, error)
	DeleteTag(string) error
	PartitionTag([]byte, []byte) (map[string][]byte, map[string][]byte, error)
	DeleteFileFromTag(string, string, string) error
	ExtendTags(map[string][]byte) error
	DiscardTags([]string) error
	EditFileFromTag(string, *chord.TagEncoder) error

	GetPath() string
	ChangePath(string)

	Clear() error
}
type DiskDictionary struct {
	Path string
	//data map[string]void // Internal dictionary
	//lock map[string]string // Lock states
	Hash func() hash.Hash // Hash function to use
}

func NewDiskDictionary(hash func() hash.Hash, path string) *DiskDictionary {
	return &DiskDictionary{
		//data: make(map[string]void),
		//lock: make(map[string]string),
		Path: path,
		Hash: hash,
	}
}

func (dictionary *DiskDictionary) GetPath() string {
	return dictionary.Path
}
func (dictionary *DiskDictionary) ChangePath(newPath string) {

	dictionary.Path = newPath
}

func (dictionary *DiskDictionary) GetFile(fileName, fileExtension string) ([]byte, []string, error) {
	log.Debugf("Cargando archivo: %s\n", fileName)

	// El path del directorio destion
	fileDir := dictionary.Path + "files/" + fileName + "-" + fileExtension + "/"

	// El path del fichero destion
	filePath := fileDir + fileName + "." + fileExtension

	// El path del json de las etiquetass
	fileJson := fileDir + "tags.json"

	// La informacion del archivo
	value, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("Error cargando el archivo %s:\n%v\n", fileName, err)
		return nil, nil, status.Error(codes.Internal, "No se pudo cargar el archivo")
	}

	// La informacion del json
	tagsJson, err := ioutil.ReadFile(fileJson)
	if err != nil {
		log.Errorf("Error cargando el json %s:\n%v\n", fileName, err)
		return nil, nil, status.Error(codes.Internal, "No se pudo cargar el json")
	}

	// Donde se van a almacenar las etiquetas
	var tags []string
	json.Unmarshal(tagsJson, &tags)

	if err != nil {
		log.Errorf("Error decodificando el json %s:\n%v\n", fileName, err)
		return nil, nil, status.Errorf(codes.Internal, fmt.Sprintf("No se pudo decodificar el json: %s\n", err))
	}
	return value, tags, nil
}

func (dictionary *DiskDictionary) SetFile(file *chord.TagFile) error {
	//dictionary.data[file.Name] = void{}

	log.Debugf("Guardando el archivo: %s\n", file.Name)

	// Creando el directorio
	if isthere := FileIsThere(dictionary.Path + "files/" + file.Name + "-" + file.Extension); !isthere {
		err := os.Mkdir(dictionary.Path+"files/"+file.Name+"-"+file.Extension, 0666)
		if err != nil {
			log.Errorf("No se pudo crear el directorio %s:\n%v\n", file.Name, err)
			return status.Error(codes.Internal, "No se pudo crear el directorio")
		}
	}
	// Path para crear el archivo
	fileDir := dictionary.Path + "files/" + file.Name + "-" + file.Extension + "/"
	filePath := fileDir + file.Name + "." + file.Extension
	fileJson := fileDir + "tags.json"

	// Guarda la informacion en el fichero filepath, si no existe lo crea
	err := ioutil.WriteFile(filePath, file.File, 0666)
	if err != nil {
		log.Errorf("Error guardando  el archivo %s:\n%v\n", file.Name, err)
		return status.Error(codes.Internal, "No se pudo guardar el archivo")
	}
	outJSON, err := json.MarshalIndent(file.Tags, "", "\t")
	if err != nil {
		log.Error("Error creando el archivo donde se van a guardar las etiquetas")
		return status.Error(codes.Internal, "No se pudo guardar las etiquetas")
	}
	err = ioutil.WriteFile(fileJson, outJSON, 0644)
	if err != nil {
		return nil
	}
	return nil
}
func (dictionary *DiskDictionary) DeleteFile(fileName, fileExtension string) error {
	log.Debugf("Eliminando el archivo: %s\n", fileName)

	//Path del archivo a eliminar
	filePath := dictionary.Path + "files/" + fileName + "-" + fileExtension

	if isthere := FileIsThere(filePath); !isthere {
		return nil
	}
	err := os.RemoveAll(filePath)
	if err != nil {
		log.Errorf("Error vaciando el directorio:\n%v\n", err)
		return status.Error(codes.Internal, "Error vaciando el directorio")
	}
	err = os.Remove(filePath)

	if err != nil {
		log.Errorf("Error eliminando el archivo:\n%v\n", err)
		return status.Error(codes.Internal, "No se pudo eliminar el archivo")
	}
	return nil
}

type TagEncoding struct {
	FileName      string `json:"file_name`
	FileExtension string `json:"file_extension"`
	NodeID        []byte `json:"node_id"`
	NodeIP        string `json:"node_ip"`
	NodePort      string `json:"node_port"`
}

func (dictionary *DiskDictionary) SetTag(tag, fileName, fileExtension string, node *chord.Node) error {

	log.Debugf("Guardando la informacion relacionada a la etiqueta: %s\n", tag)
	// Creando el directorio
	if isthere := FileIsThere(dictionary.Path + "tags/" + tag); !isthere {
		// Creando el archivo
		err := os.Mkdir(dictionary.Path+"tags/"+tag, 0666)
		if err != nil {
			log.Errorf("No se pudo crear el directorio %s:\n%v\n", tag, err)
			return status.Error(codes.Internal, "No se pudo crear el directorio")
		}

	}
	// Path del archivo
	filePath := dictionary.Path + "tags/" + tag + "/" + tag + "." + "json"

	// La lista de la informacion a almacenar
	var save []TagEncoding
	// Creo el objeto que voy a guardar
	en := TagEncoding{
		FileName:      fileName,
		FileExtension: fileExtension,
		NodeID:        node.ID,
		NodeIP:        node.IP,
		NodePort:      node.Port,
	}

	// Comprueba si hay creado un json, o sea si se almaceno informacion anteriormente
	// En caso de que eso ocurriera, se guarda en save los archivos previos
	if isThere := FileIsThere(filePath); isThere {
		jsonInfo, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Errorf("No se pudo abrir el JSON", err)
			return status.Error(codes.Internal, "No se pudo crear el archivo")
		}
		err = json.Unmarshal(jsonInfo, &save)
		if err != nil {
			log.Errorf("No se pudo abrir el JSON: %v\n", err)
			return status.Error(codes.Internal, "No se pudo crear el archivo")
		}
		// En otro caso se crea el json
	} else {
		// Debe leer el archivo para obtener la informacion
		fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Errorf("No se pudo crear el archivo %s:\n%v\n", tag+".json", err)
			return status.Error(codes.Internal, "No se pudo crear el archivo")
		}
		fd.Close()

	}
	// Agregado el nuevo elemento
	save = append(save, en)

	Cjson, err := json.MarshalIndent(&save, "", "  ")
	if err != nil {
		log.Errorf("No se pudo modificar el archivo correspondiente a la etiqueta %s\n%v\n", tag, err)
		return status.Error(codes.Internal, "No se pudo modificar el archivo")

	}

	return ioutil.WriteFile(filePath, Cjson, 0644)

}

func (dictionary *DiskDictionary) GetTag(tag string) ([]TagEncoding, error) {
	log.Debugf("Recuperando la informacion relacionada a la etiqueta: %s\n", tag)
	// Creando el directorio
	if isthere := FileIsThere(dictionary.Path + "tags/" + tag); !isthere {

		log.Errorf("No existe el directorio: %s\n", tag)
		return nil, status.Error(codes.Internal, "No existe el directorio: "+tag)

	}
	// Path del archivo
	filePath := dictionary.Path + "tags/" + tag + "/" + tag + "." + "json"
	fileDir := dictionary.Path + "tags/" + tag
	// La lista de la informacion a almacenar
	var save []TagEncoding

	// Comprueba si hay creado un json, o sea si se almaceno informacion anteriormente
	// En caso de que eso ocurriera, se guarda en save los archivos previos
	if isThere := FileIsThere(filePath); isThere {
		jsonInfo, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Errorf("No se pudo abrir el JSON", err)
			return nil, status.Error(codes.Internal, "No se pudo abrir el JSON")
		}
		err = json.Unmarshal(jsonInfo, &save)
		if err != nil {
			log.Errorf("No se pudo abrir el JSON: %v\n", err)
			return nil, status.Error(codes.Internal, "No se pudo abrir el archivo")
		}
		if len(save) == 0 {
			log.Errorf("El JSON esta vacio")
			os.RemoveAll(fileDir)
			os.Remove(fileDir)
			return nil, status.Error(codes.Internal, "El JSON esta vacio")
		}

		// En caso de que la carpeta este vacia
	} else {
		log.Errorf("No existe un JSON dentro de la carpeta: %s\n ", tag)
		return nil, status.Error(codes.Internal, "No existe un JSON dentro de la carpeta:"+tag+"\n")

	}

	return save, nil
}
func (dictionary *DiskDictionary) DeleteFileFromTag(tag, fileName, fileExtension string) error {

	log.Debugf("Recuperando la informacion relacionada a la etiqueta: %s\n Para despues eliminar la relacionada al fichero %s", tag, fileName)
	// Creando el directorio
	if isthere := FileIsThere(dictionary.Path + "tags/" + tag); !isthere {

		log.Errorf("No existe el directorio: %s\n", tag)
		//return status.Error(codes.Internal, "No existe el directorio: "+tag)
		return nil
	}
	// Path del archivo
	filePath := dictionary.Path + "tags/" + tag + "/" + tag + "." + "json"
	fileDir := dictionary.Path + "tags/" + tag
	// La lista de la informacion a almacenar
	var save []TagEncoding

	// Comprueba si hay creado un json, o sea si se almaceno informacion anteriormente
	// En caso de que eso ocurriera, se guarda en save los archivos previos
	if isThere := FileIsThere(filePath); isThere {
		jsonInfo, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Errorf("No se pudo abrir el JSON", err)
			return status.Error(codes.Internal, "No se pudo abrir el JSON")
		}
		err = json.Unmarshal(jsonInfo, &save)
		if err != nil {
			log.Errorf("No se pudo abrir el JSON: %v\n", err)
			return status.Error(codes.Internal, "No se pudo abrir el archivo")
		}
		if len(save) == 0 {
			log.Errorf("El JSON esta vacio")
			os.RemoveAll(fileDir)
			os.Remove(fileDir)
			return status.Error(codes.Internal, "El JSON esta vacio")
		}
		var temp []TagEncoding
		for _, value := range save {
			// Si es el archivo adecuado se ignora
			if value.FileName == fileName && value.FileExtension == fileExtension {
				continue

				// En otro caso se almacena
			} else {
				temp = append(temp, value)
			}
		}

		Cjson, err := json.MarshalIndent(&save, "", "  ")
		if err != nil {
			log.Errorf("No se pudo modificar el archivo correspondiente a la etiqueta %s\n%v\n", tag, err)
			return status.Error(codes.Internal, "No se pudo modificar el archivo")
		}

		return ioutil.WriteFile(filePath, Cjson, 0644)
		// En caso de que la carpeta este vacia
	} else {
		log.Errorf("No existe un JSON dentro de la carpeta: %s\n ", tag)
		return status.Error(codes.Internal, "No existe un JSON dentro de la carpeta:"+tag+"\n")

	}

}
func (dictionary *DiskDictionary) DeleteTag(tag string) error {

	log.Debugf("Eliminando la etiqueta: %s\n  del almacenamiento local", tag)
	// Creando el directorio
	if isthere := FileIsThere(dictionary.Path + "tags/" + tag); !isthere {

		log.Errorf("No existe el directorio: %s\n", tag)
		//return status.Error(codes.Internal, "No existe el directorio: "+tag)
		return nil
	}
	// Path del archivo
	fileDir := dictionary.Path + "tags/" + tag
	err := os.RemoveAll(fileDir)
	if err != nil {
		log.Errorf("Error al vaciar el directorio %s\n", tag)
		return err
	}
	err = os.Remove(fileDir)
	if err != nil {
		log.Errorf("Error eliminando el directorio %s\n", tag)
		return err
	}

	return nil
}

func (dictionary *DiskDictionary) PartitionTag(L, R []byte) (map[string][]byte, map[string][]byte, error) {
	in := make(map[string][]byte)
	out := make(map[string][]byte)
	all := false

	var fileDir = dictionary.Path + "tags"
	if Equals(L, R) {
		all = true
	}

	// Lista todos los directorios dentro de tags.
	// O sea lista todas las etiquetas almacenadas en el almacenamiento local
	keys, err := ioutil.ReadDir(fileDir)
	if err != nil {
		log.Error("No se pudo listar los directorios de las etiquetas")
		return nil, nil, err
	}

	for _, dir := range keys {
		// Me muevo por todos los directorios
		if dir.IsDir() {

			// El nombre de uno de los directorios, o sea la informacion de una de las etiquetas
			key := dir.Name()

			// Recupero la informacion que estaba almacenando
			//value, _ := ioutil.ReadFile(fileDir + key + "/" + key)
			value, _ := dictionary.GetTag(key)
			new_value, _ := json.MarshalIndent(&value, "", " ")
			/*
				var new_value []*chord.TagEncoder

				for _, encoding := range value {
					new_encoding := &chord.TagEncoder{
						FileName:      encoding.FileName,
						FileExtension: encoding.FileExtension,
						NodeID:        encoding.NodeID,
						NodeIP:        encoding.NodeIP,
						NodePort:      encoding.NodePort,
					}
					new_value = append(new_value, new_encoding)
				}
			*/
			// Si las etiquetas van almacenadas aqui
			if between, err := KeyBetween(key, dictionary.Hash, L, R); (all || between) && err == nil {
				in[key] = new_value

				// Si no van almacenadas aqui
			} else if err == nil {
				out[key] = new_value

				// Si hubo algun error
			} else {
				return nil, nil, err
			}

		}
	}
	return in, out, nil
}
func (dictionary *DiskDictionary) PartitionFile(L, R []byte) ([]*chord.TagFile, []*chord.TagFile, error) {

	in := make([]*chord.TagFile, 0)
	out := make([]*chord.TagFile, 0)
	all := false

	var fileDir = dictionary.Path + "files"
	if Equals(L, R) {
		all = true
	}

	// Lista todos los directorios dentro de files.
	// O sea lista todas los ficheros almacenados en el almacenamiento local
	keys, err := ioutil.ReadDir(fileDir)
	if err != nil {
		log.Error("No se pudo listar los directorios de las etiquetas")
		return nil, nil, err
	}

	for _, dir := range keys {
		// Me muevo por todos los directorios
		if dir.IsDir() {

			// El nombre de uno de los directorios, o sea la informacion de uno de los archivos
			key := dir.Name()

			// Cuando se guardo el archivo se almaceno en una carpeta identificada como fileName-fileExtension
			// chain es el resultado de separar fileName y fileExtension
			chain := strings.Split(key, "-")

			// Se hace un Join porque puede ser que en el nombre hubiese algun -, por lo que se tiene en cuenta
			// Dicho caso
			fileName := strings.Join(chain[:len(chain)-1], "-")

			// La extension va a ser el ultimo element de la formada de hacer el Split
			fileExtension := chain[len(chain)-1]

			// Recupero la informacion que estaba almacenando
			//value, _ := ioutil.ReadFile(fileDir + key + "/" + key)
			info, tags, _ := dictionary.GetFile(fileName, fileExtension)

			file := &chord.TagFile{
				Name:      fileName,
				Extension: fileExtension,
				File:      info,
				Tags:      tags,
			}

			// Si los archivos van almacenadoss aqui
			if between, err := KeyBetween(fileName, dictionary.Hash, L, R); (all || between) && err == nil {

				in = append(in, file)
				// Si no van almacenados aqui
			} else if err == nil {

				out = append(out, file)
				// Si hubo algun error
			} else {
				return nil, nil, err
			}

		}
	}
	return in, out, nil
}

func (dictionary *DiskDictionary) ExtendTags(data map[string][]byte) error {
	for key, value := range data {
		log.Debugf("Guardando la informacion relacionada a la etiqueta: %s\n", key)
		// Creando el directorio
		if isthere := FileIsThere(dictionary.Path + "tags/" + key); !isthere {
			// Creando el archivo
			err := os.Mkdir(dictionary.Path+"tags/"+key, 0666)
			if err != nil {
				log.Errorf("No se pudo crear el directorio %s:\n%v\n", key, err)
				return status.Error(codes.Internal, "No se pudo crear el directorio")
			}

		}
		// Path del archivo
		filePath := dictionary.Path + "tags/" + key + "/" + key + "." + "json"
		// Comprueba si hay creado un json, o sea si se almaceno informacion anteriormente
		// En caso de que eso ocurriera, se guarda en save los archivos previos
		return ioutil.WriteFile(filePath, value, 0666)

	}
	return nil
}

func (dictionary *DiskDictionary) ExtendFiles(data []*chord.TagFile) error {
	for _, value := range data {
		err := dictionary.SetFile(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dictionary *DiskDictionary) DiscardTags(data []string) error {
	fileDir := dictionary.Path + "tags/"
	for _, key := range data {
		err := dictionary.DeleteTag(fileDir + key)
		if err != nil {
			return err
		}
	}

	return nil
}
func (dictionary *DiskDictionary) DiscardFiles(data []string) error {

	var fileName string
	var fileExtension string
	for _, key := range data {

		chain := strings.Split(key, "-")
		fileName = strings.Join(chain[:len(chain)-1], "-")
		fileExtension = chain[len(chain)-1]

		err := dictionary.DeleteFile(fileName, fileExtension)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dictionary *DiskDictionary) Clear() error {
	// Se elimina todo del directorio de archivos
	filesPath := dictionary.Path + "files"
	err := os.RemoveAll(filesPath)
	if err != nil {
		return err
	}
	// Se elimina todo del directorio de etiquetas
	tagsPath := dictionary.Path + "tags"

	err = os.RemoveAll(tagsPath)
	if err != nil {
		return err
	}
	return nil
}

func (dictionary *DiskDictionary) EditFileFromTag(tag string, mod *chord.TagEncoder) error {

	log.Debugf("Recuperando la informacion relacionada a la etiqueta: %s\n Para despues modificar la informacion relacionada al archivo", tag, mod.FileName+"-"+mod.FileExtension)
	// Creando el directorio
	if isthere := FileIsThere(dictionary.Path + "tags/" + tag); !isthere {

		log.Errorf("No existe el directorio: %s\n", tag)
		//return status.Error(codes.Internal, "No existe el directorio: "+tag)
		return nil
	}
	// Path del archivo
	filePath := dictionary.Path + "tags/" + tag + "/" + tag + "." + "json"
	fileDir := dictionary.Path + "tags/" + tag
	// La lista de la informacion a almacenar
	var save []TagEncoding

	// Comprueba si hay creado un json, o sea si se almaceno informacion anteriormente
	// En caso de que eso ocurriera, se guarda en save los archivos previos
	if isThere := FileIsThere(filePath); isThere {
		jsonInfo, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Errorf("No se pudo abrir el JSON", err)
			return status.Error(codes.Internal, "No se pudo abrir el JSON")
		}
		err = json.Unmarshal(jsonInfo, &save)
		if err != nil {
			log.Errorf("No se pudo abrir el JSON: %v\n", err)
			return status.Error(codes.Internal, "No se pudo abrir el archivo")
		}
		if len(save) == 0 {
			log.Errorf("El JSON esta vacio")
			os.RemoveAll(fileDir)
			os.Remove(fileDir)
			return nil
		}
		var temp []TagEncoding
		for _, value := range save {
			// Si es el archivo adecuado se modifica y se almacena
			if value.FileName == mod.FileName && value.FileExtension == mod.FileExtension {
				value.NodeID = mod.NodeID
				value.NodeIP = mod.NodeIP
				value.NodePort = mod.NodePort
				temp = append(temp, value)
				// En otro caso se almacena directamente
			} else {
				temp = append(temp, value)
			}
		}

		Cjson, err := json.MarshalIndent(&save, "", "  ")
		if err != nil {
			log.Errorf("No se pudo modificar el archivo correspondiente a la etiqueta %s\n%v\n", tag, err)
			return status.Error(codes.Internal, "No se pudo modificar el archivo")
		}

		return ioutil.WriteFile(filePath, Cjson, 0644)
		// En caso de que la carpeta este vacia
	} else {
		log.Errorf("No existe un JSON dentro de la carpeta: %s\n ", tag)
		return status.Error(codes.Internal, "No existe un JSON dentro de la carpeta:"+tag+"\n")

	}

}
