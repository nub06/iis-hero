package cmd

import (
	"github.com/nub06/iis-hero/util"
	"github.com/spf13/cobra"
)

var virtualDir = &cobra.Command{
	Use:   "vdir",
	Short: "Manage Virtual Directories on IIS",
	Args:  cobra.ExactArgs(0),
}

var virtualDirList = &cobra.Command{
	Use:     "list",
	Short:   "List Virtual Directories and properties on IIS",
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		c := setRemoteComputerDetails()
		util.MakeTable(c.VirtualDirList())

	},
}

var virtualDirCreate = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a Virtual Directory on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		site, _ := cmd.Flags().GetString("site")
		path, _ := cmd.Flags().GetString("path")
		app, _ := cmd.Flags().GetString("app")

		c := setRemoteComputerDetails()
		c.VirtualDirCreate(site, name, path, app)

	},
}

var virtualDirRemove = &cobra.Command{
	Use:     "remove <name>",
	Short:   "Remove a Virtual Directory on IIS",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		site, _ := cmd.Flags().GetString("site")
		app, _ := cmd.Flags().GetString("app")

		c := setRemoteComputerDetails()

		c.VirtualDirRemove(name, site, app)

	},
}

var virtualDirChange = &cobra.Command{
	Use:   "change <name>",
	Short: "Change the properties of Virtual Directory on IIS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		site, _ := cmd.Flags().GetString("site")
		path, _ := cmd.Flags().GetString("path")
		rename, _ := cmd.Flags().GetString("rename")
		pool, _ := cmd.Flags().GetString("pool")

		c := setRemoteComputerDetails()
		c.VdirChangeProps(name, site, pool, path, rename)

	},
}

var virtualDirFind = &cobra.Command{
	Use:   "find <name>",
	Short: "Find a Virtual Directory ",
	//Long:  `change the site for the application`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		c := setRemoteComputerDetails()
		c.FindSiteNameByVdir(name)

	},
}

func init() {

	virtualDirCreate.Flags().String("path", "", "Identify Physical Path Value for the Virtual Directory")
	virtualDirCreate.Flags().String("site", "", "Identify Web Site Name Value for the Virtual Directory")
	virtualDirCreate.Flags().String("app", "", "Use this flag to Identify a Application Name If you want to create Virtual Directory under of an Application")

	virtualDirRemove.Flags().String("site", "", "Specify the Web Site Name of the Virtual Directory")
	virtualDirRemove.Flags().String("app", "", "Specify the Application Name of the Virtual Directory")

	virtualDirChange.Flags().String("path", "", "Set a new Physical Path Value for the Virtual Directory")
	virtualDirChange.Flags().String("site", "", "Specify a Web Site Name  Value for the Virtual Directory")
	virtualDirChange.Flags().String("rename", "", "Set a new Name  Value for the Virtual Directory")
	virtualDirChange.Flags().String("pool", "", "Set a new Application  Value for the Virtual Directory")

}
