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
	"github.com/nub06/iis-hero/service"
	"github.com/nub06/iis-hero/util"
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
	Long: `iis-hero is a CLI management tool for IIS servers. It provides access and management capabilities to remote IIS servers. 
It supports a wide range of operations for application pools, web sites, applications, and virtual directories.
In addition to IIS, it also provides support for Windows Services and file operations. You can manage Windows services and perform operations such as file transfer between remote servers.
`,
}

var execCmd = &cobra.Command{
	Use:     "exec",
	Short:   "Execute custom powershell commands on target computer",
	Args:    cobra.ExactArgs(1),
	Example: `iis-hero exec "iisreset /status"`,

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
		service.RemoveCurrentConf()

		profile, _ := cmd.Flags().GetString("profile")

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

			fmt.Println(color.HiCyanString("iis-hero will be run on target computer(localhost) '%s'", remoteHost))
			fmt.Println(color.HiCyanString("Current profile information is resetting."))

			isRemote = false

			viper.Set("isRemote", isRemote)

		} else {
			isRemote = true
			viper.Set("isRemote", isRemote)
			cyan := color.New(color.FgHiCyan).SprintFunc()

			credInfo := fmt.Sprintf("%s: %s\n%s: %s\n%s: %s\n%s: %s",
				cyan("Domain"), cyan(remoteDomain),
				cyan("Host"), cyan(remoteHost),
				cyan("Username"), cyan(remoteUsername),
				cyan("Password"), cyan(remotePassword))

			fmt.Println(credInfo)

			if profile != "" {
				service.SaveConfig(profile)
				service.UseConfig(profile)

			} else {

				fmt.Println(cyan("Current Profile: Empty"))
				fmt.Println(color.HiGreenString("Profile information is not specified."))
				fmt.Println(color.HiGreenString("Current profile information is resetting."))

			}

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

		service.RemoveCurrentConf()
	},
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "This command allows you to manage your configuration profiles",
}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "This command allows you to save the credentials created with the 'iis-hero login' command as a configuration profile",
	Example: `iis-hero profile save dev
iis-hero profile save --name dev`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name, _ := cmd.Flags().GetString("name")

		if len(args) == 1 && name != "" {

			log.Fatal(color.HiRedString("You cannot use both an argument and the --name flag together. Please either provide an argument or use the --name flag."))
		}

		if name == "" {
			if len(args) != 0 {
				name = args[0]
			} else {

				log.Fatalf(color.HiRedString("You must specify a configuration profile name when saving. Use the following command to save a configuration profile: %s", color.HiCyanString("\niis-hero login save --name <profile name>\niis-hero login save <profile name>")))
			}
		}

		service.SaveConfig(name)
	},
}

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use command allows you to switch between configuration profiles.",
	Long:  `If you have previously saved a configuration profile with the 'iis-hero profile save' command, you can start using a profile you've saved before with the 'use' command`,
	Example: `iis-hero profile use dev
iis-hero profile use --name dev`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name, _ := cmd.Flags().GetString("name")

		if len(args) == 1 && name != "" {

			log.Fatal(color.HiRedString("You cannot use both an argument and the --name flag together. Please either provide an argument or use the --name flag."))
		}

		if name == "" {
			if len(args) != 0 {
				name = args[0]
			} else {

				log.Fatalf(color.HiRedString("You must specify a configuration profile name when trying to change the profile. Please use the following command to specify a configuration profile %s", color.HiCyanString("\niis-hero login use --name <profile name>\niis-hero login use <profile name>")))
			}
		}
		service.UseConfig(name)
	},
}

var confRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "This command allows you to save the credentials created with the 'iis-hero login' command as a configuration profile",
	Example: `iis-hero profile save dev
iis-hero profile save --name dev`,
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"rm"},
	Run: func(cmd *cobra.Command, args []string) {

		var confName string

		isAll, _ := cmd.Flags().GetBool("all")

		if len(args) != 1 || isAll {
			if isAll && len(args) == 1 {
				log.Fatal(color.HiRedString("You can not use Profile Name and --all flag together."))
			}
		} else {
			confName = args[0]
		}

		service.DeleteConfiguration(confName, isAll)

	},
}

var showCmd = &cobra.Command{
	Use:     "list",
	Short:   "This command lists saved configuration profiles.",
	Example: `iis-hero profile list`,
	Aliases: []string{"ls"},

	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		err, configList := service.ShowSavedConfigs()

		if err != nil {

			log.Fatal(color.HiRedString("Cannot find any saved configuration profile"))

		} else {

			util.MakeTable(configList)
		}
	},
}

var showCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "This command displays the currently used configuration profiles.",
	Long: `If you create new credentials with the 'iis-hero login' command,
your current configuration profile information will appear empty until you save this information with the 'iis-hero profile save' command`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		err, config := service.ShowCurrentConfig()

		if err != nil {
			log.Fatalf(color.HiRedString("Cannot find the current configuration profile. You may not have saved your configuration or you may not have specified a profile with the \n'iis-hero profile use <profile name>' command.\nor you can save your configuration using the following command:\niis-hero profile save --name <profile name>"))

		} else {
			util.MakeTable(config)

		}
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
	profileCmd.AddCommand(saveCmd)
	profileCmd.AddCommand(useCmd)
	profileCmd.AddCommand(showCmd)
	profileCmd.AddCommand(confRemoveCmd)
	profileCmd.AddCommand(showCurrentCmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.AddCommand(credCmd)
	rootCmd.AddCommand(winSvcCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(folderCmd)

	rootCmd.AddCommand(profileCmd)

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
	loginCmd.Flags().String("profile", "", "Identify the configuration profile name to save as profile")

	saveCmd.Flags().String("name", "", "Identify the configuration profile name")
	useCmd.Flags().String("name", "", "Identify the configuration profile name")

	confRemoveCmd.Flags().BoolP("all", "a", false, "Remove all saved Configuration Profiles")
}
