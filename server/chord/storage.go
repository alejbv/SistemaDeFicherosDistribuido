package chord

import (
	"encoding/json"
	"errors"
	"hash"
	"io/ioutil"
	"os"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/chord/chord"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storage interface {
	Get(string) ([]byte, error)
	//Set(string, []byte) error
	Set(*chord.TagFile) error
	SetTag(string, string, string, *chord.Node) error
	Delete(string) error

	/*
		GetWithLock(string, string) ([]byte, error)
		SetWithLock(*chord.TagFile, string) error
		DeleteWithLock(string, string) error
		Lock(string, string) error
		Unlock(string, string) error

	*/
	Partition([]byte, []byte) (map[string][]byte, map[string][]byte, error)
	Extend(map[string][]byte) error
	Discard([]string) error
	Clear() error
}

type Dictionary struct {
	data map[string][]byte // Internal dictionary
	lock map[string]string // Lock states
	Hash func() hash.Hash  // Hash function to use
}

func NewDictionary(hash func() hash.Hash) *Dictionary {
	return &Dictionary{
		data: make(map[string][]byte),
		lock: make(map[string]string),
		Hash: hash,
	}
}

func (dictionary *Dictionary) Get(key string) ([]byte, error) {
	value, ok := dictionary.data[key]

	if !ok {
		return nil, errors.New(" Key not found")
	}

	return value, nil
}

func (dictionary *Dictionary) Set(key string, value []byte) error {
	dictionary.data[key] = value
	return nil
}

func (dictionary *Dictionary) Delete(key string) error {
	delete(dictionary.data, key)
	return nil
}

func (dictionary *Dictionary) GetWithLock(key, permission string) ([]byte, error) {
	lock, ok := dictionary.lock[key]

	if !ok || permission == lock {
		return dictionary.Get(key)
	}

	return nil, os.ErrPermission
}

func (dictionary *Dictionary) SetWithLock(key string, value []byte, permission string) error {
	lock, ok := dictionary.lock[key]

	if !ok || permission == lock {
		return dictionary.Set(key, value)
	}

	return os.ErrPermission
}

func (dictionary *Dictionary) DeleteWithLock(key string, permission string) error {
	lock, ok := dictionary.lock[key]

	if !ok || permission == lock {
		return dictionary.Delete(key)
	}

	return os.ErrPermission
}

func (dictionary *Dictionary) Partition(L, R []byte) (map[string][]byte, map[string][]byte, error) {
	in := make(map[string][]byte)
	out := make(map[string][]byte)

	if Equals(L, R) {
		return dictionary.data, out, nil
	}

	for key, value := range dictionary.data {
		if between, err := KeyBetween(key, dictionary.Hash, L, R); between && err == nil {
			in[key] = value
		} else if err == nil {
			out[key] = value
		} else {
			return nil, nil, err
		}
	}

	return in, out, nil
}

func (dictionary *Dictionary) Extend(data map[string][]byte) error {
	for key, value := range data {
		dictionary.data[key] = value
	}
	return nil
}

func (dictionary *Dictionary) Discard(data []string) error {
	for _, key := range data {
		delete(dictionary.data, key)
	}

	return nil
}

func (dictionary *Dictionary) Clear() error {
	dictionary.data = make(map[string][]byte)
	return nil
}

func (dictionary *Dictionary) Lock(key, permission string) error {
	lock, ok := dictionary.lock[key]

	if ok && permission != lock {
		return os.ErrPermission
	}

	dictionary.lock[key] = permission

	return nil
}

func (dictionary *Dictionary) Unlock(key, permission string) error {
	lock, ok := dictionary.lock[key]

	if ok && permission != lock {
		return os.ErrPermission
	}

	delete(dictionary.lock, key)
	return nil
}

// Hard storage.

type void struct{}

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

func (dictionary *DiskDictionary) Get(key string) ([]byte, error) {
	log.Debugf("Cargando archivo: %s\n", key)

	value, err := ioutil.ReadFile(key)
	if err != nil {
		log.Errorf("Error cargando el archivo %s:\n%v\n", key, err)
		return nil, status.Error(codes.Internal, "No se pudo cargar el archivo")
	}

	return value, nil
}

func (dictionary *DiskDictionary) Set(file *chord.TagFile) error {
	//dictionary.data[file.Name] = void{}

	log.Debugf("Guardando el archivo: %s\n", file.Name)

	// Creando el directorio
	if isthere := FileIsThere(dictionary.Path + "files/" + file.Name); !isthere {
		err := os.Mkdir(dictionary.Path+"files/"+file.Name, 0666)
		if err != nil {
			log.Errorf("No se pudo crear el directorio %s:\n%v\n", file.Name, err)
			return status.Error(codes.Internal, "No se pudo crear el directorio")
		}
	}
	// Path para crear el archivo
	filePath := dictionary.Path + "files/" + file.Name + "/" + file.Name + "." + file.Extension

	// Creando el archivo
	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Errorf("No se pudo crear el archivo %s:\n%v\n", file.Name+"."+file.Extension, err)
		return status.Error(codes.Internal, "No se pudo crear el archivo")
	}
	defer fd.Close()

	err = ioutil.WriteFile(filePath, file.File, 0644)
	if err != nil {
		log.Errorf("Error guardando  el archivo %s:\n%v\n", file.Name, err)
		return status.Error(codes.Internal, "No se pudo guardar el archivo")
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
		err := os.Mkdir(dictionary.Path+"tags/"+tag, 0666)
		if err != nil {
			log.Errorf("No se pudo crear el directorio %s:\n%v\n", tag, err)
			return status.Error(codes.Internal, "No se pudo crear el directorio")
		}

	}
	// Path del archivo
	filePath := dictionary.Path + "tags/" + tag + "/" + tag + "." + "json"
	// Creando el archivo

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

func (dictionary *DiskDictionary) Delete(key string) error {
	err := os.Remove(key)
	delete(dictionary.data, key)

	if err != nil {
		log.Errorf("Error eliminando el archivo:\n%v\n", err)
		return status.Error(codes.Internal, "No se pudo eliminar el archivo")
	}
	return nil
}

/*
func (dictionary *DiskDictionary) GetWithLock(key, permission string) ([]byte, error) {
	lock, ok := dictionary.lock[key]

	if !ok || permission == lock {
		return dictionary.Get(key)
	}

	return nil, os.ErrPermission
}

func (dictionary *DiskDictionary) SetWithLock(file *chord.TagFile, permission string) error {

	lock, ok := dictionary.lock[file.GetName()]

	if !ok || permission == lock {
		return dictionary.Set(file)
	}

	return os.ErrPermission
}

func (dictionary *DiskDictionary) DeleteWithLock(key string, permission string) error {
	lock, ok := dictionary.lock[key]

	if !ok || permission == lock {
		return dictionary.Delete(key)
	}

	return os.ErrPermission
}

func (dictionary *DiskDictionary) Lock(key, permission string) error {
	lock, ok := dictionary.lock[key]

	if ok && permission != lock {
		return os.ErrPermission
	}

	dictionary.lock[key] = permission

	return nil
}

func (dictionary *DiskDictionary) Unlock(key, permission string) error {
	lock, ok := dictionary.lock[key]

	if ok && permission != lock {
		return os.ErrPermission
	}

	delete(dictionary.lock, key)
	return nil
}


*/
func (dictionary *DiskDictionary) Partition(L, R []byte) (map[string][]byte, map[string][]byte, error) {
	in := make(map[string][]byte)
	out := make(map[string][]byte)
	all := false

	if Equals(L, R) {
		all = true
	}

	for key := range dictionary.data {
		value, err := dictionary.Get(key)
		if err != nil {
			return nil, nil, err
		}

		if between, err := KeyBetween(key, dictionary.Hash, L, R); (all || between) && err == nil {
			in[key] = value
		} else if err == nil {
			out[key] = value
		} else {
			return nil, nil, err
		}
	}

	return in, out, nil
}

func (dictionary *DiskDictionary) Extend(data map[string][]byte) error {
	for key, value := range data {
		err := dictionary.Set(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dictionary *DiskDictionary) Discard(data []string) error {
	for _, key := range data {
		err := dictionary.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dictionary *DiskDictionary) Clear() error {
	keys := Keys(dictionary.data)

	for _, key := range keys {
		err := dictionary.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
