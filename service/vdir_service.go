package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/model"
	"github.com/nub06/iis-hero/util"
)

func (r RemoteComputer) VirtualDirCreate(sitename string, virtualdirName string, path string, app string) {

	if isStrNotEmpty(virtualdirName) {

		if r.isSiteExist(sitename) {

			if app != "" {

				r.VdirAddToApp(virtualdirName, sitename, path, app)

			} else {

				psCommand := fmt.Sprintf(`Import-Module WebAdministration
			New-WebVirtualDirectory -Site "%s" -Name "%s" -PhysicalPath "%s"
			`, sitename, virtualdirName, path)
				r.ExecuteCommand(psCommand)
			}
		} else {

			log.Fatal(color.HiRedString(fmt.Sprintf("WebSite %s couldn't find", sitename)))
		}

	}

}

func (r RemoteComputer) VirtualDirRemove(vdir string, sitename string, appname string) {

	if sitename == "" && appname == "" {

		fmt.Println(color.HiGreenString("%s%s", "Finding virtual directory: ", vdir))
		sitename, appname = r.findSiteAndApplicationName(vdir)

	}

	var psCommand string

	fmt.Println(color.HiGreenString("%s%s", "Removing virtual directory: ", vdir))

	if appname == "" {

		psCommand = fmt.Sprintf(`Import-Module WebAdministration
			Remove-Item IIS:\Sites\%s\%s -Recurse
			`, sitename, vdir)

	} else {

		psCommand = fmt.Sprintf(`Import-Module WebAdministration
			Remove-Item IIS:\Sites\%s\%s\%s -Recurse
			`, sitename, appname, vdir)
	}

	r.ExecuteCommand(psCommand)

}

func (r RemoteComputer) VirtualDirList() []model.VirtualDir {
	psCommand := `	Import-Module WebAdministration;

	$sites = Get-ChildItem IIS:\Sites;
	
	$output = foreach ($site in $sites) {
		$webvdirectories = Get-WebVirtualDirectory -Site $site.name;
		
		foreach ($webvdirectory in $webvdirectories) {
			[PSCustomObject]@{
				SiteName = $site.name;
				ApplicationName = "";
				VirtualDirectory = $webvdirectory.path.split("/")[-1];
				PhysicalPath = $webvdirectory.physicalPath;
			}
		}
		
		$applications = Get-WebApplication -Site $site.Name -ErrorAction SilentlyContinue;
	
		foreach ($application in $applications) {
			$vdirs = Get-WebVirtualDirectory -Application $application.Path;
	
			foreach ($vdir in $vdirs) {
				[PSCustomObject]@{
					SiteName = $site.Name;
					ApplicationName = $application.Path.Split('/', 2)[-1];
					VirtualDirectory = $vdir.Path.Split('/', 2)[-1];
					PhysicalPath = $vdir.PhysicalPath;
				}
			}
		}
	}
	
	$output 
	`

	res := r.RunCommandJSON(psCommand)
	return util.JsonToStructVirtualDir(res)

}

func (r RemoteComputer) VdirChangeProps(vdir string, site string, pool string, path string, rename string) {

	if site == "" {
		fmt.Println(color.HiGreenString(`Finding the "%s"...`, vdir))
		site = r.FindSiteNameByVdir(vdir)
		fmt.Println(color.HiGreenString(`Virtual Directory: "%s" Found on the Web Site: "%s"`, vdir, site))
	} else if vdir == "" {
		log.Fatal(color.HiRedString("Specify a virtual directory"))
	} else {
		if !r.isSiteExist(site) {
			log.Fatal(color.HiRedString(`Web Site "%s" cannot found`, pool))
		}
	}

	r.vDirsetAppPool(vdir, site, pool)
	r.vDirsetName(vdir, site, rename)
	r.vDirsetPath(vdir, site, path)

}

func (r RemoteComputer) vDirsetPath(vdir string, site string, path string) {
	if isStrNotEmpty(path) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration
		Set-WebConfiguration -Filter "/system.applicationHost/sites/site[@name='%s']/application/virtualDirectory[@path='/%s']" -Value @{physicalPath="%s"}
		`, site, vdir, path)
		r.ExecuteCommand(psCommand)

	}

}

func (r RemoteComputer) vDirsetName(vdir string, site string, rename string) {
	if isStrNotEmpty(rename) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration
		Set-WebConfiguration -Filter "/system.applicationHost/sites/site[@name='%s']/application/virtualDirectory[@path='/%s']" -Value @{path="/%s"}
		`, site, vdir, rename)
		r.ExecuteCommand(psCommand)

	}

}

func (r RemoteComputer) vDirsetAppPool(vdir string, site string, pool string) {

	if isStrNotEmpty(pool) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration
		Set-ItemProperty "IIS:\Sites\%s\%s" -Name Name -Value "%s"
		`, site, vdir, pool)
		r.ExecuteCommand(psCommand)

	}

}

func (r RemoteComputer) FindSiteNameByVdir(vdir string) string {

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration

	$virtualDirectoryName = "%s"
	$sites = Get-ChildItem -Path "IIS:\Sites"
	$siteNames = @()
	
	foreach ($site in $sites) {
		$virtualDirectories = Get-WebVirtualDirectory -Site $site.Name -Name $virtualDirectoryName -ErrorAction SilentlyContinue
		if ($virtualDirectories) {
			$siteNames += $site.Name
		}
		
		$applications = Get-WebApplication -Site $site.Name
		foreach ($application in $applications) {
			$virtualDirectories = Get-WebVirtualDirectory -Application $application.Path -Name $virtualDirectoryName -ErrorAction SilentlyContinue
			if ($virtualDirectories) {
				$siteNames += $site.Name
			}
		}
	}
	
	if ($siteNames.Count -eq 0) {
		Write-Output " Error! The virtual directory '$virtualDirectoryName' does not exist in any website or application."
	} else {
		Write-Output "$($siteNames -join ", ")"
	}`, vdir)

	fmt.Println(color.HiGreenString(`%s"%s"`, "Finding website of virtual directory: ", vdir))

	res := r.RunCommandPlain(psCommand)

	res = strings.TrimRight(res, "\r\n")

	if strings.Contains(res, "Error") {

		log.Fatal(color.HiRedString(res))
	}

	n := strings.Split(res, ",")

	if len(n) > 1 {

		log.Fatal(color.HiRedString("Multiple websites were found with the same virtual directory.\nPlease specify the name of the website."))
	}

	return res
}

func (r RemoteComputer) findSiteAndApplicationName(vdir string) (string, string) {

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration
	$virtualDirectoryName = "%s"
	$sites = Get-ChildItem -Path "IIS:\Sites"
	$siteInfos = @()
	
	foreach ($site in $sites) {
		$virtualDirectories = Get-WebVirtualDirectory -Site $site.Name -Name $virtualDirectoryName -ErrorAction SilentlyContinue
		if ($virtualDirectories) {
			$siteInfo = @{
				SiteName = $site.Name
				ApplicationName = $null
			}
			$siteInfos += $siteInfo
		}
	
		$applications = Get-WebApplication -Site $site.Name
		foreach ($application in $applications) {
			$virtualDirectories = Get-WebVirtualDirectory -Application $application.Path -Name $virtualDirectoryName -ErrorAction SilentlyContinue
			if ($virtualDirectories) {
				$siteInfo = @{
					SiteName = $site.Name
					ApplicationName = $application.Path.Substring(1)
				}
				$siteInfos += $siteInfo
			}
		}
	}
	
	if ($siteInfos.Count -eq 0) {
		$res= "Error! The virtual directory $virtualDirectoryName does not exist in any website or application."
		Write-Output $res
	} else {
		Write-Output $siteInfos
	}`, vdir)

	res := r.RunCommandJSON(psCommand)

	m := util.JsonToVdirInfo(res)

	return m.SiteName, m.ApplicationName

}

func (r RemoteComputer) VdirAddToApp(vdir string, site string, path string, app string) {

	if isStrNotEmpty(vdir) {
		psCommand := fmt.Sprintf(`
		$siteName = "%s"
		$appPath = "/%s"
		$vdPath = "%s"
		$physicalPath = "%s"		
		New-Item "IIS:\Sites\$siteName\$appPath\$vdPath" -type VirtualDirectory -physicalPath $physicalPath
		`, site, app, vdir, path)
		r.ExecuteCommand(psCommand)

	}

}

func (r RemoteComputer) IsVdirHasApplication(app string, site string) []string {
	vdir := r.VirtualDirList()
	vdirMap := make(map[string][]string)

	for _, s := range vdir {
		if s.ApplicationName != "" && s.SiteName == site {
			vdirMap[s.ApplicationName] = append(vdirMap[s.ApplicationName], s.VirtualDirectory)
		}
	}

	for k, v := range vdirMap {
		if strings.EqualFold(k, app) {
			//return strings.Join(v, ",")
			return v
		}
	}

	return []string{}
}
