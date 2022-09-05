/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/alejbv/SistemaDeFicherosDistribuido/client"
	"github.com/spf13/cobra"
)

// listFileCmd represents the listFile command
var listFileCmd = &cobra.Command{
	Use:   "listFile",
	Short: "Lista los archivos y sus etiqueta de",
	Long: `Este comando recibe una lista de etiquetas, donde identifica a todos
	los archivos del sistema que son representados por esta y devuelve de estos sus nombres y las etiquetas
	con las que fueron almcenados`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listFile llamado")

		chordClient, _ := client.StartClient()

		err := client.ClientListFiles(chordClient, args)
		if err != nil {

			fmt.Println("Hubo un error mientras se ejecutaba el comando listFile")
		}
	},
}

func init() {
	rootCmd.AddCommand(listFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
