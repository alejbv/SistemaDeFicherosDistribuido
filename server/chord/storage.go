package chord

import (
	"github.com/alejbv/SistemaDeFicherosDistribuido/chord"
)

type Storage interface {
	// Recupera la informacion almacenada de un archivo y las etiquetas que lo identifican
	GetFile(string, string) ([]byte, []string, error)

	// Almacena localmente un nuevo archivo
	SetFile(*chord.TagFile) error

	// Elimina un archivo local
	DeleteFile(string, string) error

	// Particiona los archivos almacenados localmente en 2 dependiedo si se ubican o no dentro de un
	// Determinado rango
	PartitionFile([]byte, []byte) ([]*chord.TagFile, []*chord.TagFile, error)

	// Agrega una lista nueva de archivos al almacenamiento local
	ExtendFiles([]*chord.TagFile) error

	// Elimina un conjunto de archviso
	DiscardFiles([]string) error

	GetFileInfo(string, string) ([]byte, error)

	DeleteTagsFromFile(string, string, []string) ([]string, error)
	AddTagsToFile(string, string, []string) (bool, error)
	/*
	   				Metodos de las etiquetas
	   ==============================================================================================
	*/

	// Agrega a la informacion almacenada de una determinada etiqueta una referencia a un nuevo archivo
	SetTag(string, string, string, *chord.Node) error

	// Recupera la informacion que almacena una etiqueta
	GetTag(string) ([]TagEncoding, error)

	// Elimina la informacion de una etiqueta
	DeleteTag(string) error

	// Particiona las etiquetas almacenadas localmente en 2 dependiedo si se ubican o no dentro de rango
	PartitionTag([]byte, []byte) (map[string][]byte, map[string][]byte, error)

	// Elimina la informacion de un archivo de una etiqueta
	DeleteFileFromTag(string, string, string) error

	// Modifica la informacion de un archivo de una etiqueta
	EditFileFromTag(string, *chord.TagEncoder) error

	// Agrega la informacion de una lista de archivos a una etiqueta
	AddFileToTag(string, []*chord.TagEncoder) error

	// Agrega un nuevo conjunto de etiqueta al almacenamiento local
	ExtendTags(map[string][]byte) error

	// Elimina un conjunto de etiquetas
	DiscardTags([]string) error

	// Se obtiene el path local donde se trabaja
	GetPath() string

	// Permite cambiar el path
	ChangePath(string)

	// Crea las carpetes donde se almacenaran los archivos y las etiquetas
	Inicialize()

	// Elimina toda la informacion almacenada localmente
	Clear() error
}
