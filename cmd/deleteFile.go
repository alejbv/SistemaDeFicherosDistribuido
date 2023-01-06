/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/alejbv/SistemaDeFicherosDistribuido/client"
	"github.com/spf13/cobra"
)

// deleteFileCmd represents the deleteFile command
var deleteFileCmd = &cobra.Command{
	Use:   "deleteFile",
	Short: "Elimina un conjunto de ficheros del sistema",
	Long: `Este comando recibe una lista de etiquetas y luego procede a eliminar del sistema
	a todos los ficheros que tienen a dicha lista como un subconjunto de sus etiquetas`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deleteFile llamado")

		chordClient, conn, _ := client.StartClient()
		defer conn.Close()

		err := client.ClientDeleteFile(chordClient, args)
		if err != nil {

			fmt.Println("Hubo un error mientras se ejecutaba el comando deleteFile")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
