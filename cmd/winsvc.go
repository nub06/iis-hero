package cmd

import (
	"fmt"

	"github.com/nub06/iis-hero/util"
	"github.com/spf13/cobra"
)

var winSvcCmd = &cobra.Command{
	Use:     "winsvc",
	Aliases: []string{"ws"},
	Short:   "Manage Windows Services",
	//Long:  `Set the winsvc for the application`,
	Args: cobra.ExactArgs(0),
}

var startwinSvcCmd = &cobra.Command{
	Use:   "start <service name>",
	Short: "Start the Windows Service",
	Long: `Start the Windows Service
You can use <service name> or <display name>`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		serviceName := args[0]
		c := setRemoteComputerDetails()
		c.StartWindowsService(serviceName)

	},
}

var stopwinSvcCmd = &cobra.Command{
	Use:   "stop <service name>",
	Short: "Stop the Windows Service",
	Long: `Stop the Windows Service
You can use <service name> or <display name>`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		serviceName := args[0]
		c := setRemoteComputerDetails()
		c.StopWindowsService(serviceName)
	},
}

var restartwinSvcCmd = &cobra.Command{
	Use:   "restart <service name>",
	Short: "Restart the Windows Service",
	Long: `Restart the Windows Service
You can use <service name> or <display name>`,
	Aliases: []string{"reset"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		serviceName := args[0]
		c := setRemoteComputerDetails()
		c.RestartWindowsService(serviceName)
	},
}

var winSvcCmdState = &cobra.Command{
	Use:   "state <service name>",
	Short: "Display the status of a Windows Service",
	Long: `state the winsvc for the application
You can use <service name> or <display name>`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		serviceName := args[0]
		c := setRemoteComputerDetails()
		x, _ := (c.GetWindowsServiceState(serviceName))
		fmt.Println(x)
	},
}

var winSvcCmdList = &cobra.Command{
	Use:   "list <service name>",
	Short: "List Windows Service and properties",
	Long: `list the winsvc for the application
You can use <service name> or <display name>`,
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		serviceName := args[0]
		c := setRemoteComputerDetails()
		util.MakeTableFromStruct(c.GetWindowsServiceStats(serviceName))
	},
}

var winSvcCmdCreate = &cobra.Command{
	Use:   "create <service name>",
	Short: "Create a new Windows Service",
	//Long:  `create the winsvc for the application`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		serviceName := args[0]
		displayName, _ := cmd.Flags().GetString("displayname")
		description, _ := cmd.Flags().GetString("description")
		exePath, _ := cmd.Flags().GetString("exepath")
		startup, _ := cmd.Flags().GetString("startup")

		c := setRemoteComputerDetails()
		c.CreateWinsvc(serviceName, displayName, description, exePath, startup)
	},
}

var winSvcRemove = &cobra.Command{
	Use:   "remove <service name>",
	Short: "Remove a Windows Service with given name",
	Long: `Removes a Windows Service with given name
You can use <service name> or <display name>
`,
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		serviceName := args[0]

		c := setRemoteComputerDetails()

		c.DeleteWinSvc(serviceName)
	},
}

var winSvcCmdChange = &cobra.Command{
	Use:   "change <service name>",
	Short: "Change the properties of the existing Windows Service",
	Long: `Changes the properties of the existing Windows Service
You can use <service name> or <display name>`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		serviceName := args[0]
		displayName, _ := cmd.Flags().GetString("displayname")
		description, _ := cmd.Flags().GetString("description")
		exePath, _ := cmd.Flags().GetString("exepath")
		startmode, _ := cmd.Flags().GetString("startmode")

		c := setRemoteComputerDetails()
		c.ChangeWinSvc(serviceName, displayName, description, exePath, startmode)
	},
}

func init() {

	winSvcCmd.AddCommand(startwinSvcCmd)
	winSvcCmd.AddCommand(stopwinSvcCmd)
	winSvcCmd.AddCommand(restartwinSvcCmd)
	winSvcCmd.AddCommand(winSvcCmdState)
	winSvcCmd.AddCommand(winSvcCmdList)
	winSvcCmd.AddCommand(winSvcCmdChange)
	winSvcCmd.AddCommand(winSvcCmdCreate)
	winSvcCmd.AddCommand(winSvcRemove)

	winSvcCmdCreate.Flags().String("displayname", "<service name>", "Identify a DisplayName value for Windows Service")
	winSvcCmdCreate.Flags().String("description", "Description of <service name>", "Identify a Description value for Windows Service")
	winSvcCmdCreate.Flags().String("exepath", "", "Identify the Executable Path value of Windows Service")
	winSvcCmdCreate.Flags().StringP("startup", "", "Automatic", "Identify the StartupType value for Windows Service")

	winSvcCmdChange.Flags().String("displayname", "", "Set new DisplayName value for Windows Service")
	winSvcCmdChange.Flags().String("description", "", "Set new Description value for Windows Service")
	winSvcCmdChange.Flags().String("exepath", "", "Set new Executable Path value for Windows Service")
	winSvcCmdChange.Flags().String("startmode", "", "Set new StartupType value for Windows Service")

}
