package cmd

import (
	"fmt"

	"github.com/nub06/iis-hero/util"
	"github.com/spf13/cobra"
)

var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Manage folders on target computer",
	Args:  cobra.ExactArgs(0),
}

var folderListCmd = &cobra.Command{
	Use:     "list <folderpath>",
	Short:   "List folders on given path",
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		folderName := args[0]
		c := setRemoteComputerDetails()
		model := c.FolderList(folderName)
		fmt.Println(model)
		util.MakeTable(c.FolderList(folderName))

	},
}

var folderBackupCmd = &cobra.Command{
	Use:   "backup <source folder path>",
	Short: "Create a backup of the specified folder on the target computer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		folderName := args[0]
		path, _ := cmd.Flags().GetString("dest")
		c := setRemoteComputerDetails()
		c.CreateBackupDir(folderName, path)

	},
}

var folderPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Perform file copying from target computer to your local computer",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()

		rem, _ := cmd.Flags().GetString("target")
		local, _ := cmd.Flags().GetString("local")

		c.CopyFromTarget(rem, local)

	},
}

var folderPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Perform file copying from your local computer to the target computer",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()

		rem, _ := cmd.Flags().GetString("target")
		local, _ := cmd.Flags().GetString("local")

		c.CopyToTarget(rem, local)

	},
}

func init() {

	folderCmd.AddCommand(folderListCmd)
	folderCmd.AddCommand(folderBackupCmd)
	folderCmd.AddCommand(folderPullCmd)
	folderCmd.AddCommand(folderPushCmd)

	folderBackupCmd.Flags().StringP("dest", "d", "D:\\Backups", "Specify backup folder destination path")

	folderPullCmd.Flags().String("target", "", "Specify file/folder path on target computer")
	folderPullCmd.Flags().String("local", "", " Specify destination folder path on local computer")

	folderPushCmd.Flags().String("local", "", "Specify file/folder path on local computer")
	folderPushCmd.Flags().String("target", "", "Specify destination folder path on target computer")
}
