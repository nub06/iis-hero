package cmd

import (
	"github.com/nub06/iis-hero/util"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage IIS Applications",
	Args:  cobra.ExactArgs(0),
}

var appCreateCmd = &cobra.Command{
	Use:   "create <appname>",
	Short: "Create an Application on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		app := args[0]
		site, _ := cmd.Flags().GetString("site")
		pool, _ := cmd.Flags().GetString("pool")
		path, _ := cmd.Flags().GetString("path")

		c := setRemoteComputerDetails()
		c.CreateSiteApp(app, site, path, pool)

	},
}

var appListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Applications and properties on IIS",
	Args:    cobra.ExactArgs(0),
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		util.MakeTable(c.SiteAppList())

	},
}

var appRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an Application on IIS",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm"},

	Run: func(cmd *cobra.Command, args []string) {

		app := args[0]
		site, _ := cmd.Flags().GetString("site")
		force, _ := cmd.Flags().GetBool("force")

		c := setRemoteComputerDetails()

		if force {
			c.RemoveSiteApplicationForce(app, site, force)
		} else {
			c.RemoveSiteApplication(app, site, force)
		}

	},
}

var appChangeCmd = &cobra.Command{
	Use:   "change <appname>",
	Short: "Change the properties of an Application on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		app := args[0]
		site, _ := cmd.Flags().GetString("site")
		pool, _ := cmd.Flags().GetString("pool")
		preload, _ := cmd.Flags().GetString("preload")
		path, _ := cmd.Flags().GetString("path")
		rename, _ := cmd.Flags().GetString("rename")

		c := setRemoteComputerDetails()
		c.ChangeApplicationProps(site, app, pool, preload, path, rename)
	},
}

func init() {

	appCmd.AddCommand(appCreateCmd)
	appCmd.AddCommand(appRemoveCmd)
	appCmd.AddCommand(appChangeCmd)

	appCreateCmd.Flags().String("site", "", "Identify Site Name for Application")
	appCreateCmd.Flags().String("pool", "", "Identify Application Pool for Application")
	appCreateCmd.Flags().String("path", "", "Identify Physical Path for Application")

	appChangeCmd.Flags().String("site", "", "Identify Site Name for Application")
	appChangeCmd.Flags().String("pool", "", "Set a new Application Pool for Application")
	appChangeCmd.Flags().String("path", "", "Set a new Physical Path for Application")
	appChangeCmd.Flags().String("preload", "", "Set a PreloadEnabled value for Application")
	appChangeCmd.Flags().String("rename", "", "Set a new name for Application")

	appRemoveCmd.Flags().String("site", "", "Identify Site Name for Application")
	appRemoveCmd.Flags().BoolP("force", "f", false, "Force flag for the Remove Application")

}
