/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var poolName string

var remoteHost string
var remoteUsername string
var remotePassword string
var remoteDomain string
var isRemote bool

var rootCmd = &cobra.Command{
	Use: "iis-hero",
	Long: `
	iis-hero is a CLI management tool for IIS servers. It provides access and management capabilities to remote IIS servers. 
It supports a wide range of operations for application pools, web sites, applications, and virtual directories.
In addition to IIS, it also provides support for Windows services and file operations. You can manage Windows services and perform operations such as file transfer between remote servers.
`,
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute custom powershell commands on target computer",
	Args:  cobra.ExactArgs(1),
	Long:  `A longer description`,

	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		fmt.Println(c.RunCommandPlain(args[0]))
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Identifies the credentials of the target computer",
	Long:  `Identifies the credentials of the target computer on which the application will run`,
	Example: `iis-hero login -c <ComputerName>
iis-hero -c localhost	
iis-hero  -c <ComputerName> -d <Domain> -u <UserName> -p <Password>`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		clearViperInfo()

		if remoteHost == "local" || remoteHost == "localhost" {

			remoteHost, _ = os.Hostname()

		}
		viper.Set("remoteHost", encrypt(remoteHost))
		viper.Set("remoteUsername", encrypt(remoteUsername))
		viper.Set("remotePassword", encrypt(remotePassword))
		viper.Set("remoteDomain", encrypt(remoteDomain))

		if remoteHost == "" {
			log.Fatal(color.HiRedString("'--computer, -c' flag cannot be empty"))
		}

		viper.WriteConfig()

		if remoteHost == "" || remoteUsername == "" || remoteDomain == "" || remotePassword == "" {

			//The isRemote flag is being set to false because one of the required credentials is empty.

			isRemote = false

			viper.Set("isRemote", isRemote)

		} else {
			isRemote = true
			viper.Set("isRemote", isRemote)
			fmt.Println(remoteDomain, remoteHost, remoteUsername, remotePassword)
		}
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error saving configuration file:", err)
			return
		}

	},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the current login credentials",
	Long:  `This command clears the credential information that is currently being used`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		clearViperInfo()

		if err := viper.WriteConfig(); err != nil {
			fmt.Println(color.HiRedString("Error saving configuration file:"), err)
			return
		}
		fmt.Println(color.HiGreenString("Target computer credentials have been successfully cleared"))
	},
}

func Execute() error {

	configFilePath := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "config.yaml")
	cfgPath := filepath.Join(os.Getenv("APPDATA"), "iis-hero")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(configFilePath), 0755); err != nil {
			panic(err)
		}
		if _, err := os.Create(configFilePath); err != nil {
			panic(err)
		}
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return rootCmd.Execute()
}

func init() {

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(loginCmd)
	loginCmd.AddCommand(clearCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.AddCommand(credCmd)
	rootCmd.AddCommand(winSvcCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(folderCmd)

	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(iisStopCmd)
	rootCmd.AddCommand(iisStartCmd)

	rootCmd.AddCommand(poolCmd)
	rootCmd.AddCommand(siteCmd)
	rootCmd.AddCommand(iisConfigCmd)

	rootCmd.AddCommand(virtualDir)
	virtualDir.AddCommand(virtualDirList)
	virtualDir.AddCommand(virtualDirCreate)
	virtualDir.AddCommand(virtualDirRemove)
	virtualDir.AddCommand(virtualDirFind)
	virtualDir.AddCommand(virtualDirChange)

	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appListCmd)

	loginCmd.PersistentFlags().StringVarP(&remoteDomain, "domain", "d", "", "Identify the domain name")
	loginCmd.PersistentFlags().StringVarP(&remoteHost, "computer", "c", "", "Identify the target computer hostname")
	loginCmd.PersistentFlags().StringVarP(&remoteUsername, "username", "u", "", "Identify the user name")
	loginCmd.PersistentFlags().StringVarP(&remotePassword, "password", "p", "", "Identify the password")

}
