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

// addTagsCmd represents the addTags command
var addTagsCmd = &cobra.Command{
	Use:   "addTags",
	Short: "Añade etiquetas a los archivos",
	Long: `Este comando recibe 2 listas de etiquetas, donde con la primera identifica a todos
	los archivos del sistema que son representados por esta y le agregan las etiquetas de la segunda lista`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addTags llamado")
		chordClient, _ := client.StartClient()
		queryTags := strings.Split(args[0], "-")
		addTags := strings.Split(args[1], "-")
		err := client.ClientAddTags(chordClient, queryTags, addTags)
		if err != nil {

			fmt.Println("Hubo un error mientras se ejecutaba el comando addTags")
		}
	},
}

func init() {
	rootCmd.AddCommand(addTagsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTagsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTagsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
