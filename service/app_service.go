package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/model"
	"github.com/nub06/iis-hero/util"
)

func (r RemoteComputer) IsAppExist(app string, site string) bool {

	psCommand := fmt.Sprintf(`Import-Module WebAdministration
	$siteName = "%s"
	$appName = "%s"
	$app = Get-WebApplication -Site $siteName -Name $appName -ErrorAction SilentlyContinue
	if ($app) {
		Write-Output "true"
	} else {
		Write-Output "false"
	}
	
	`, site, app)

	res := r.RunCommandPlain(psCommand)

	return strings.Contains(res, "true")

}
func (r RemoteComputer) CreateSiteApp(app string, site string, path string, pool string) {

	if r.isSiteExist(site) && path != "" && pool != "" {

		psCommand := fmt.Sprintf(`
		Import-Module WebAdministration
		$siteName = "%s"
		$appName = "/%s"
		$appPath = "%s"
		$appPool = "%s"	
		$site = Get-Item IIS:\Sites\$siteName	
		New-WebApplication -Site $site.Name -Name $appName -PhysicalPath $appPath -ApplicationPool $appPool	
		`, site, app, path, pool)

		fmt.Println(r.RunCommandPlain(psCommand))

	}

}

func (r RemoteComputer) RemoveSiteApplicationForce(app string, site string, force bool) {

	if site == "" {

		site = r.findSiteOfApplication(app)
	}

	if app != "" && site != "" {

		subApps := r.FindSubApps(app)

		for _, s := range subApps {

			r.RemoveSiteApplication(s, site, force)
		}

		r.RemoveSiteApplication(app, site, force)

	}
}

func (r RemoteComputer) RemoveSiteApplication(app string, site string, force bool) {

	if site == "" {
		site = r.findSiteOfApplication(app)

		if site == "" {

			fmt.Println(color.HiRedString("Web Site of Application: %s couldn't find", app))
			msg := fmt.Sprintf(`Specify to site of application: "%s"`, app)
			log.Fatal(color.HiRedString(msg))

		}

	} else {
		if !r.IsAppExist(app, site) {

			msg := fmt.Sprintf(`Application "%s" is couldn't find on site "%s".`, app, site)

			log.Fatal(color.HiRedString(msg))
		}
	}

	if app != "" && site != "" && !force {

		subAppList := r.FindSubApps(app)

		if len(subAppList) > 0 {

			var sb strings.Builder
			ms := fmt.Sprintf("The Application %s includes sub applications: %s\nFor remove the Application and all sub applications\n", app, strings.Join(subAppList, ","))
			sb.WriteString(color.HiRedString(ms))
			sb.WriteString(color.HiRedString("e.g iis-hero site app remove <app> -f"))

			log.Fatal(sb.String())

		}

	}

	msg := fmt.Sprintf(`Application "%s" is removing...`, app)
	fmt.Println(color.HiGreenString(msg))

	psCommand := fmt.Sprintf(`		
	Import-Module WebAdministration
	Remove-WebApplication -Name "%s" -Site "%s"
		`, app, site)

	r.ExecuteCommand(psCommand)

}

func (r RemoteComputer) SiteAppList() []model.SitesApplications {
	psCommand := `Import-Module WebAdministration;
	$Websites = Get-ChildItem IIS:\Sites;
	$output = @()
	
	foreach ($Site in $Websites) {
		$WebApplications = Get-WebApplication -Site $Site.name
	
		foreach ($WebApp in $WebApplications) {
			$ApplicationName = $WebApp.path.TrimStart("/")
			$OutputObject = [ordered]@{
				Website = $Site.name
				ApplicationName = $ApplicationName
				PhysicalPath = $WebApp.physicalPath
				ApplicationPool = $WebApp.applicationPool
				PreloadEnabled = $WebApp.preloadEnabled
			}
	
			$Output += $OutputObject
		}
	}
	
	$Output	
	`
	res := r.RunCommandJSON(psCommand)
	return util.JsonToStructApp(res)

}

func (r RemoteComputer) FindSubApps(app string) []string {
	res := r.SiteAppList()
	var subAppList []string
	app = strings.ToLower(app)
	for _, s := range res {
		if strings.HasPrefix(strings.ToLower(s.ApplicationName), app+"/") {
			if strings.ToLower(s.ApplicationName) != app {
				subAppList = append(subAppList, s.ApplicationName)
			}
		}
	}

	return subAppList
}

func (r RemoteComputer) ChangeApplicationProps(site string, app string, pool string, preload string, path string, rename string) {

	if site == "" {

		site = r.findSiteOfApplication(app)

		if site == "" {

			log.Fatal(color.HiRedString("Please specify a website for application %s", app))

		}
	}

	if app != "" && site != "" {
		if !r.IsAppExist(app, site) {
			log.Fatal(color.HiRedString(fmt.Sprintf("%s%s", "Application", "couldn't find on web site")))
		} else {
			r.AppChangePath(site, app, path)
			r.AppChangePool(site, app, pool)
			r.AppChangePreload(site, app, preload)
			r.AppRename(site, app, rename)
		}

	} else {
		log.Fatal(color.HiRedString("Identify application and website name"))
	}

}

func (r RemoteComputer) AppChangePool(site string, app string, pool string) {

	if pool != "" {
		psCommand := fmt.Sprintf(
			`Import-Module WebAdministration
			$site = "%s"
			$app = "%s"
			$appPath = "IIS:\Sites\$site\$app"
			Set-ItemProperty "IIS:\Sites\$site\$app" -Name "applicationPool" -Value "%s"
			`, site, app, pool)

		if !r.IsPoolExist(pool) {

			fmt.Println(color.HiRedString(`The ApplicationPool "%s" is couldn't find. Try a validate Application Pool`, pool))
		}

		fmt.Println(color.HiGreenString(`Changing pool of application "%s"`, app))

		r.RunCommandPlain(psCommand)

	}

}

func (r RemoteComputer) AppChangePreload(site string, app string, preload string) {

	if preload != "" {
		psCommand := fmt.Sprintf(`	
			Import-Module WebAdministration
			$site = "%s"
			$app = "%s"
			$appPath = "IIS:\Sites\$site\$app"
			Set-ItemProperty "IIS:\Sites\$site\$app" -Name "preloadEnabled" -Value %s
			`, site, app, preload)

		fmt.Println(color.HiGreenString(`Changing preloadenabled value of application "%s"`, app))

		r.RunCommandPlain(psCommand)

	}

}

func (r RemoteComputer) AppChangePath(site string, app string, path string) {

	if path != "" {
		psCommand := fmt.Sprintf(`	
		Import-Module WebAdministration
		$siteName = "%s
		$appPath = "/%s"
		$newPath = "%s"
		Set-WebConfiguration -Filter "/system.applicationHost/sites/site[@name='$siteName']/application[@path='$appPath']/virtualDirectory[@path='/']" -Value @{physicalPath="$newPath"}
			`, site, app, path)

		fmt.Println(color.HiGreenString(`Changing physical path of application "%s"`, app))
		r.RunCommandPlain(psCommand)

	}

}

func (r RemoteComputer) AppRename(site string, app string, rename string) {

	if rename != "" {
		psCommand := fmt.Sprintf(`	
	Import-Module WebAdministration
	$siteName = "%s"
	$appPath = "/%s"
	$rename = "%s"
	Set-WebConfigurationProperty -Filter "/system.applicationHost/sites/site[@name='$siteName']/application[@path='$appPath']" -Name "path" -Value "/$rename"
		`, site, app, rename)

		fmt.Println(color.HiGreenString(`Renaming application "%s"`, app))

		r.RunCommandPlain(psCommand)

	}

}

func (r RemoteComputer) findSiteOfApplication(app string) string {

	psCommand := fmt.Sprintf(`	
	Import-Module WebAdministration;
	$applicationName = "%s"
	$sites = Get-ChildItem IIS:\Sites	
	foreach ($site in $sites) {
		$applications = Get-ChildItem "IIS:\Sites\$($site.name)\"
	
		foreach ($application in $applications) {
			if ($application.Path -eq "/$applicationName") {
				Write-Host "$($site.name)"
			}
		}
	}`, app)

	res := r.RunCommandPlain(psCommand)

	return strings.TrimRight(res, "\r\n")

}
