/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/alejbv/SistemaDeFicherosDistribuido/server/logging"
	"github.com/alejbv/SistemaDeFicherosDistribuido/server/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serverStartCmd represents the serverStart command
var serverStartCmd = &cobra.Command{
	Use:   "serverStart",
	Short: "Empieza a correr el sistema como un nodo del chord",
	Long: `Este comando establece a el sistema como un nodo del chord, realizando
	todos los procedimientos y pasos para inicializar el sistema`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serverStart llamado")
		logging.SettingLogger(log.InfoLevel, "server")
		rsaPrivateKeyPath := "pv.pem"
		rsaPublicteKeyPath := "pub.pem"
		network := "tcp"
		services.StartServer(network, rsaPrivateKeyPath, rsaPublicteKeyPath)
		log.Info("Servicio Iniciado...")
		for {

		}

	},
}

func init() {
	rootCmd.AddCommand(serverStartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverStartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverStartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
