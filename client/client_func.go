package client

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/alejbv/SistemaDeFicherosDistribuido/chord"
)

func AddFile(client chord.ChordClient, filePath string, tags []string) error {

	chain := strings.Split(filePath, "/")
	tempName := strings.Split(chain[len(chain)-1], ".")
	fileName := tempName[0]
	fileExtension := tempName[1]

	log.Printf("Se va a subir el archivo %s con etiquetas\n %s", fileName, tags)

	fileInfo, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("Error mientras se cargaba el archivo:\n" + err.Error())
		return errors.New("Error mientras se cargaba el archivo:\n" + err.Error())
	}
	tagFile := &chord.TagFile{
		Name:      fileName,
		Extension: fileExtension,
		File:      fileInfo,
		Tags:      tags,
	}
	req := &chord.AddFileRequest{
		File:    tagFile,
		Replica: false,
	}
	// Se obtiene el contexto de la conexion y el tiempo de espera de la request.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = client.AddFile(ctx, req)
	if err != nil {
		log.Errorf("Error mientras se subia el archivo:\n" + err.Error())
		return errors.New("Error mientras se subia el archivo:\n" + err.Error())

	}
	return nil
}
func DeleteFile(client chord.ChordClient, tags []string) error {

	log.Printf("Se van a eliminar todos los archivos  con etiquetas\n %s", tags)
	// Se obtiene el contexto de la conexion y el tiempo de espera de la request.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.DeleteFileByQuery(ctx, &chord.DeleteFileByQueryRequest{Tag: tags})
	if err != nil {
		log.Errorf("Error mientras se eliminaban los archivos:\n" + err.Error())
		return errors.New("Error mientras se eliminaban los archivos:\n" + err.Error())

	}
	return nil
}
func ListFiles(client chord.ChordClient, tags []string) error {

	log.Printf("Se va a obtener la informacion de los archivos con etiquetas\n %s", tags)
	// Se obtiene el contexto de la conexion y el tiempo de espera de la request.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := client.ListByQuery(ctx, &chord.ListByQueryRequest{Tags: tags})
	if err != nil {
		log.Errorf("Error mientras se obtenia la informacion:\n" + err.Error())
		return errors.New("Error mientras se obtenia la informacion:\n" + err.Error())

	}
	info := res.Response
	for key, value := range info {
		var fileTags []string

		json.Unmarshal(value, &fileTags)
		log.Printf("Archivos \t\t Etiquetas")
		log.Printf("%s\t", key)
		space := "\t"
		for _, tag := range fileTags {
			log.Printf(space+"%s\n", tag)
			space = "\t\t"
		}

	}
	return nil

}
func AddTags(client chord.ChordClient, queryTags, addTags []string) error {

	log.Printf("Se van a agregar a todos los archivos  con etiquetas\n %s\n Las etiquetas:%s", queryTags, addTags)
	// Se obtiene el contexto de la conexion y el tiempo de espera de la request.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.AddTagsByQuery(ctx, &chord.AddTagsByQueryRequest{QueryTags: queryTags,
		AddTags: addTags})
	if err != nil {
		log.Errorf("Error mientras se agregaban las etiquetas:\n" + err.Error())
		return errors.New("Error mientras se agregaban las etiquetas:\n" + err.Error())

	}
	return nil
}

func DeleteTags(client chord.ChordClient, queryTags, deleteTags []string) error {

	log.Printf("Se van a eliminar de todos los archivos  con etiquetas\n %s\n Las etiquetas:%s", queryTags, deleteTags)
	// Se obtiene el contexto de la conexion y el tiempo de espera de la request.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.DeleteTagsByQuery(ctx, &chord.DeleteTagsByQueryRequest{QueryTags: queryTags,
		RemoveTags: deleteTags})
	if err != nil {
		log.Errorf("Error mientras se eliminaban las etiquetas:\n" + err.Error())
		return errors.New("Error mientras se eliminaban las etiquetas:\n" + err.Error())

	}
	return nil
}
