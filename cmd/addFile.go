/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/alejbv/SistemaDeFicherosDistribuido/client"
	"github.com/spf13/cobra"
)

// addFileCmd represents the addFile command
var addFileCmd = &cobra.Command{
	Use:   "addFile",
	Short: "Agrega un archivo al sistema",
	Long: `Este comando recibe el path del archivo que se quiere agregar al sistema
	y una lista de las etiquetas que lo identifican`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addFile llamado")

		//Despues de comprobar que funciona se arregla esto
		chordClient, _ := client.StartClient()
		filePath := args[0]
		tags := args[1:]
		err := client.ClientAddFile(chordClient, filePath, tags)
		if err != nil {

			fmt.Println("Hubo un error mientras se ejecutaba el comando addFile")
		}

	},
}

func init() {
	rootCmd.AddCommand(addFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
