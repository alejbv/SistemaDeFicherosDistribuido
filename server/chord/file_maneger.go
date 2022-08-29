package chord

type FileManeger struct {
}

/*
Funcion para añadir archivos segun sus etiquetas
*/
type File struct {
	Info []byte
}

func (file *FileManeger) AddFiles(tags []string, Files []File) error {

	return nil
}

/*
Funcion para eliminar archivos segun sus etiquetas
*/

func (file *FileManeger) DeleteFiles(tags []string) error {

	return nil
}

/*
Funcion para devolver una lista de archivos segun sus etiquetas
*/

func (file *FileManeger) GetFiles(tags []string) ([]File, error) {

	return nil, nil
}

/*
Funcion para devolver una listas de nombres de archivos  segun sus etiquetas
*/

func (file *FileManeger) ListFiles(tags []string) ([]string, error) {

	return nil, nil
}
