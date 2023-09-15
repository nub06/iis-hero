package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var iisStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all Internet services on target computer.",
	Long:  `Performing a "iisreset /start" command on a target computer`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		fmt.Println(color.HiGreenString("IIS starting..."))
		c.StartIIS()

	},
}

var iisStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all Internet services on target computer",
	Long:  `Performing a "iisreset /stop" command on a target computer`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		fmt.Println(color.HiGreenString("IIS stopping..."))
		c.StopIIS()
	},
}

var resetCmd = &cobra.Command{
	Use:     "reset",
	Short:   "Stop and then restart all Internet services on the target computer",
	Long:    `Performing a "iisreset" command on a target computer`,
	Args:    cobra.ExactArgs(0),
	Aliases: []string{"restart"},
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		fmt.Println(color.HiGreenString("Restarting IIS..."))
		c.ResetIIS()
	},
}

var iisBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup IIS configuration",
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		var folderName string
		if len(args) == 1 {
			folderName = args[0]
			c.BackupIISConfig(folderName)
		} else if len(args) == 0 {
			c.BackupIISConfig(folderName)
		} else {
			log.Fatal(color.HiRedString("e.g: iis-hero config backup\niis-hero config backup <foldername>"))
		}
	},
}

var iisConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage IIS configuration",
}

var iisConfigRestore = &cobra.Command{
	Use:   "restore",
	Short: "Restore IIS configuration from backup.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		isLatest, _ := cmd.Flags().GetBool("latest")
		folderName, _ := cmd.Flags().GetString("name")
		c := setRemoteComputerDetails()
		c.RestoreIISConfig(isLatest, folderName)

	},
}

var iisConfigList = &cobra.Command{
	Use:     "list",
	Short:   "List IIS configuration backups.",
	Args:    cobra.ExactArgs(0),
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		c.ConfigBackupLists()

	},
}

var iisConfigClear = &cobra.Command{
	Use:   "clear",
	Short: "Clear all IIS configuration data.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		fmt.Println(color.HiGreenString("Removing all IIS configuration..."))

		c.ConfigClearAll()

	},
}

var iisBackupRemove = &cobra.Command{
	Use:     "remove",
	Short:   "Remove IIS configuration backup file",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()

		isForce, _ := cmd.Flags().GetBool("all")
		var folder string

		if len(args) == 1 {
			folder = args[0]
			c.RemoveIISBackup(folder, isForce)
		} else if len(args) == 0 && isForce {
			c.RemoveIISBackup(folder, isForce)
		} else {
			log.Fatal(color.HiRedString("e.g: iis-hero config backup remove <folder>\n iis-hero config remove -a"))
		}

	},
}

func init() {

	iisConfigCmd.AddCommand(iisBackupCmd)
	iisConfigCmd.AddCommand(iisConfigRestore)
	iisConfigCmd.AddCommand(iisConfigList)
	iisConfigCmd.AddCommand(iisConfigClear)

	iisBackupCmd.AddCommand(iisBackupRemove)

	iisConfigRestore.Flags().Bool("latest", false, "Restores the most recently created backup file")
	iisConfigRestore.Flags().String("name", "", "Restores the backup file with the specified name")

	iisBackupRemove.Flags().BoolP("all", "a", false, "Flag for the remove all backup files")
}
