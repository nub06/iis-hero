package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/util"
	"github.com/spf13/cobra"
)

var siteName string

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Manage Web Sites on IIS",
	Args:  cobra.ExactArgs(0),
}

var siteStartCmd = &cobra.Command{
	Use:   "start <sitename>",
	Short: "Start the Web Site on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		siteName = args[0]
		if siteName == "" {
			log.Fatal(color.HiGreenString("Specify a site name to run this command \ne.g: iis-hero site start <sitename>"))
		}
		c := setRemoteComputerDetails()
		c.StartIISSite(siteName)
	},
}

var siteStopCmd = &cobra.Command{
	Use:   "stop <sitename>",
	Short: "Stop the Web Site on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		siteName = args[0]

		if siteName == "" {
			log.Fatal(color.HiGreenString("Specify a site name to run this command \ne.g: iis site stop <sitename>"))
		}
		c := setRemoteComputerDetails()
		c.StopIISSite(siteName)
	},
}

var siteRemoveCmd = &cobra.Command{
	Use:     "remove <sitename>",
	Short:   "Remove a Web Site on IIS",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		siteName = args[0]
		c := setRemoteComputerDetails()
		c.RemoveSite(siteName)

	},
}

var siteListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Web Sites and properties on IIS",
	Args:    cobra.ExactArgs(0),
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()

		all, _ := cmd.Flags().GetBool("all")
		stopped, _ := cmd.Flags().GetBool("stopped")

		if siteName != "" {

			util.MakeTable(c.ListSingleSite(siteName))

		} else {
			if all {

				util.MakeTable(c.ListAllSites())
			} else if stopped {

				util.MakeTable(c.WebSitesListByStatus("Stopped"))
			} else {

				util.MakeTable(c.WebSitesListByStatus("Started"))
			}
		}

	},
}

var siteStateCmd = &cobra.Command{
	Use:     "state <sitename>",
	Short:   "Display the status of a Web Site on IIS",
	Aliases: []string{"st"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()

		siteName = args[0]
		s := (c.GetSiteState(siteName))

		fmt.Println(s)

	},
}

var siteCreateCmd = &cobra.Command{
	Use:   "create <sitename>",
	Short: "Create a Web Site on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()

		siteName = args[0]
		binding, _ := cmd.Flags().GetString("bind")
		preload, _ := cmd.Flags().GetString("preload")
		pool, _ := cmd.Flags().GetString("pool")
		path, _ := cmd.Flags().GetString("path")

		c.CreateSite(siteName, binding, path, pool, preload)

	},
}

var siteChangeCmd = &cobra.Command{
	Use:   "change <sitename>",
	Short: "Change the properties of Web Site on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		siteName = args[0]
		rename, _ := cmd.Flags().GetString("rename")
		binding, _ := cmd.Flags().GetString("bind")
		preload, _ := cmd.Flags().GetString("preload")
		pool, _ := cmd.Flags().GetString("pool")
		path, _ := cmd.Flags().GetString("path")

		c := setRemoteComputerDetails()
		c.ChangeWebSite(siteName, rename, preload, pool, binding, path)

	},
}

func init() {

	siteListCmd.Flags().BoolP("all", "a", false, "-a , --all returns all sites")
	siteListCmd.Flags().BoolP("stopped", "s", false, "-s, --stopped returns stopped sites")
	siteListCmd.Flags().StringVarP(&siteName, "name", "n", "", "Specify the name of the website you want to list")

	siteCreateCmd.Flags().String("path", "", "Identify Physical Path Value for the Web Site")
	siteCreateCmd.Flags().String("bind", "", "Identify Binding Path Value for the Web Site")
	siteCreateCmd.Flags().String("pool", "", "Identify Application Pool Name for the Web Site")
	siteCreateCmd.Flags().String("preload", "true", "Identify PreLoadEnabled Value for the Web Site")

	siteChangeCmd.Flags().String("path", "", "Set new Physical Path Value for the Web Site")
	siteChangeCmd.Flags().String("bind", "", "Set new Binding Value for the Web Site")
	siteChangeCmd.Flags().String("pool", "", "Set new Application Pool Name for the Web Site")
	siteChangeCmd.Flags().String("preload", "", "Set new PreLoadEnabled Value for the Web Site")
	siteChangeCmd.Flags().String("rename", "", "Rename Web Site")

	siteCmd.AddCommand(siteStartCmd)
	siteCmd.AddCommand(siteStopCmd)
	siteCmd.AddCommand(siteStateCmd)
	siteCmd.AddCommand(siteListCmd)
	siteCmd.AddCommand(siteCreateCmd)
	siteCmd.AddCommand(siteRemoveCmd)
	siteCmd.AddCommand(siteChangeCmd)

}
