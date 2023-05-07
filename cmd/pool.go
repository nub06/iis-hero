package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/util"
	"github.com/spf13/cobra"
)

var poolCmd = &cobra.Command{
	Use:   "pool",
	Short: "Manage Application Pools on IIS",
	Args:  cobra.ExactArgs(0),
}

var poolStartCmd = &cobra.Command{
	Use:   "start <poolname>",
	Short: "Start the Application Pool on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		poolName = args[0]

		c := setRemoteComputerDetails()
		c.StartWebAppPool(poolName)
	},
}

var poolStopCmd = &cobra.Command{
	Use:   "stop <poolname>",
	Short: "Stop the Application Pool on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		poolName = args[0]
		c := setRemoteComputerDetails()
		c.StopWebAppPool(poolName)
	},
}

var poolRestartCmd = &cobra.Command{
	Use:     "restart <poolname>",
	Short:   "Restart the Application Pool on IIS",
	Aliases: []string{"reset"},
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		all, _ := cmd.Flags().GetBool("all")

		if len(args) != 1 || all {
			if all && len(args) == 1 {
				log.Fatal(color.HiRedString("You can not use poolname and --all flag together."))
			}
		} else {
			poolName = args[0]
		}

		c.RestartWebAppPool(poolName, all)
	},
}

var poolListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Application Pools and properties on IIS",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		all, _ := cmd.Flags().GetBool("all")
		stopped, _ := cmd.Flags().GetBool("stopped")

		if all {
			util.MakeTable(c.ListAllAppPools())
		} else if stopped {
			util.MakeTable(c.ListAppPoolsByStatus("Stopped"))
		} else {
			util.MakeTable(c.ListAppPoolsByStatus("Started"))
		}

	},
}

var poolStateCmd = &cobra.Command{
	Use:     "state <poolname>",
	Short:   "Display the status of an Application Pool on IIS",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"st"},
	Run: func(cmd *cobra.Command, args []string) {

		poolName = args[0]
		c := setRemoteComputerDetails()
		fmt.Println(c.GetAppPoolState(poolName))

	},
}

var poolRemoveCmd = &cobra.Command{
	Use:   "remove <poolname>",
	Short: "Remove an Application Pool on IIS",
	Long: `This command removes an application pool with the given name from IIS.
If there is a Web Site belonging to the Application Pool that you want to delete, you should use the force flag to delete the Application Pool.`,
	Example: "iis pool remove <poolname>\niis pool remove <poolname> -f",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm"},
	Run: func(cmd *cobra.Command, args []string) {

		isForce, _ := cmd.Flags().GetBool("force")
		c := setRemoteComputerDetails()
		c.RemoveWebAppPool(isForce, args[0])

	},
}

var poolCreateCmd = &cobra.Command{
	Use:     "create <poolname>",
	Short:   "Create an Application Pool on IIS",
	Example: "iis pool create <poolname>",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := setRemoteComputerDetails()

		poolName = args[0]
		autostart, _ := cmd.Flags().GetString("autostart")
		runtime, _ := cmd.Flags().GetString("clr")
		startmode, _ := cmd.Flags().GetString("startmode")
		idleaction, _ := cmd.Flags().GetString("idleaction")
		idleminute, _ := cmd.Flags().GetString("idleminute")
		pipeline, _ := cmd.Flags().GetString("pipeline")

		if poolName != "" {
			c.AppPoolCreate(poolName, startmode, autostart, runtime, pipeline, idleaction, idleminute)
		} else {
			log.Fatal(color.HiGreenString("Please specify a poolname.\ne.g:\niis-hero pool create <poolname>"))

		}
	},
}

var poolChangeCmd = &cobra.Command{
	Use:   "change <poolname>",
	Short: "Change the properties of an Application Pool on IIS",
	Example: `iis pool change <poolname> --runtime 0 --pipeline Integrated
iis pool change <poolname> --startmode AlwaysRunning
iis pool change <poolname> --startmode OnDemand --runtime v2 --idleminute 0  `,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		poolName = args[0]
		autostart, _ := cmd.Flags().GetString("autostart")
		runtime, _ := cmd.Flags().GetString("clr")
		startmode, _ := cmd.Flags().GetString("startmode")
		idleaction, _ := cmd.Flags().GetString("idleaction")
		idleminute, _ := cmd.Flags().GetString("idleminute")
		pipeline, _ := cmd.Flags().GetString("pipeline")
		rename, _ := cmd.Flags().GetString("rename")

		c := setRemoteComputerDetails()
		c.AppPoolChange(poolName, runtime, startmode, autostart, idleaction, idleminute, pipeline, rename)

	},
}

func init() {
	poolCmd.AddCommand(poolStartCmd)
	poolCmd.AddCommand(poolStopCmd)
	poolCmd.AddCommand(poolRestartCmd)
	poolCmd.AddCommand(poolListCmd)
	poolCmd.AddCommand(poolStateCmd)
	poolCmd.AddCommand(poolCreateCmd)
	poolCmd.AddCommand(poolRemoveCmd)
	poolCmd.AddCommand(poolChangeCmd)

	poolListCmd.Flags().BoolP("all", "a", false, "Lists all Application Pools")
	poolListCmd.Flags().BoolP("stopped", "s", false, "Lists only stopped Application Pools")

	poolRestartCmd.Flags().BoolP("all", "a", false, "Restart all Application Pools")

	poolCreateCmd.Flags().String("autostart", "true", "Set Autostart value for Application Pool")
	poolCreateCmd.Flags().String("startmode", "AlwaysRunning", "Set Start Mode value for Application Pool")
	poolCreateCmd.Flags().String("clr", "", `Set Runtime Version value for Application Pool e.g: --runtime v4, --runtime v2 (default "No Managed Code")`)
	poolCreateCmd.Flags().String("idleaction", "Suspend", "Set Idle-Timeout Action value for Application Pool")
	poolCreateCmd.Flags().String("idleminute", "0", "Set Idle-Timeout(minutes) value for Application Pool")
	poolCreateCmd.Flags().String("pipeline", "Integrated", "Set Managed Pipeline Mode value for Application Pool")
	poolRemoveCmd.Flags().BoolP("force", "f", false, "Force flag to delete an Application Pool")

	poolChangeCmd.Flags().String("autostart", "", "Set new Autostart value for Application Pool")
	poolChangeCmd.Flags().String("startmode", "", "Set new Start Mode value for Application Pool")
	poolChangeCmd.Flags().String("clr", "", `Set new Runtime Version value for Application Pool`)
	poolChangeCmd.Flags().String("idleaction", "", "Set new Idle-Timeout Action value for Application Pool")
	poolChangeCmd.Flags().String("idleminute", "", "Set new Idle-Timeout(minutes) value for Application Pool")
	poolChangeCmd.Flags().String("pipeline", "", "Set new Managed Pipeline Mode value for Application Pool")
	poolChangeCmd.Flags().String("rename", "", "Set new name for Application Pool")

}
