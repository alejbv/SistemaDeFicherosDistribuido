/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/alejbv/SistemaDeFicherosDistribuido/client"
	"github.com/spf13/cobra"
)

// deleteTagsCmd represents the deleteTags command
var deleteTagsCmd = &cobra.Command{
	Use:   "deleteTags",
	Short: "Elimina las etiquetas de los archivos",
	Long: `Este comando recibe 2 listas de etiquetas, donde con la primera se identificador a todos
	los archivos del sistema que fueron  etiquetados por esta y le elimina  las etiquetas de la segunda lista`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deleteTags llamado")

		chordClient, conn, _ := client.StartClient()
		defer conn.Close()
		queryTags := strings.Split(args[0], "-")
		deleteTags := strings.Split(args[1], "-")
		err := client.ClientDeleteTags(chordClient, queryTags, deleteTags)
		if err != nil {

			fmt.Println("Hubo un error mientras se ejecutaba el comando deleteTags")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteTagsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteTagsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteTagsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
