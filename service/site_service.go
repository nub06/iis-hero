package service

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/model"
	"github.com/nub06/iis-hero/util"
)

func (r RemoteComputer) isSiteExist(siteName string) bool {

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration;
	if (Get-IISSite -Name "%s" -ErrorAction SilentlyContinue) {
		Write-Host "true"
	} else {
		Write-Host "false"
	}`, siteName)

	res := r.RunCommandPlain(psCommand)

	res = strings.TrimRight(res, "\r\n")

	return res == "true"

}

func (r RemoteComputer) RemoveSite(siteName string) {

	psCommand := fmt.Sprintf(`Remove-IISSite -Name "%s" -Confirm:$false`, siteName)

	if r.isSiteExist(siteName) {

		r.ExecuteCommand(psCommand)

		fmt.Printf(color.HiGreenString(`Web Site "%s" is removing...`), siteName)

	} else {

		errorGenerator(siteName, r.ComputerName)

	}

}

func (r RemoteComputer) CreateSite(siteName string, siteAddress string, sitePath string, poolName string, isPreloadEnabled string) {

	if strings.Contains(siteAddress, "https") {

		log.Fatal(color.HiGreenString("Use a http address"))

	}

	if siteName == "" || siteAddress == "" || sitePath == "" || poolName == "" {

		message := color.HiGreenString("Please make sure you specify all the necessary parameters to create a website.\n")
		usage := color.HiCyanString(" e.g: iis-hero site create <sitename> --pool <applicationpool name> --bind <bindaddress> --path <physicalpath> --preload <default true>")

		log.Fatal(message + usage)

	}

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration;
	$siteName = "%s";
    $physicalPath = "%s";
    $appPoolName = "%s";
    $bindings = @{Protocol="http";BindingInformation="*:80:%s"};
    New-Item IIS:\Sites\$siteName -bindings $bindings -physicalPath $physicalPath -applicationPool $appPoolName;
    Set-ItemProperty IIS:\Sites\$siteName -name applicationDefaults.preloadEnabled -value %s;
	`, siteName, sitePath, poolName, siteAddress, isPreloadEnabled)

	if r.isSiteExist(siteName) {

		fmt.Println(color.HiGreenString(`Web Site "%s" is already exist`, siteName))

	} else {

		s := color.HiGreenString("Creating your site ")
		v := color.HiMagentaString(siteName)
		fmt.Println(s + v)
		r.ExecuteCommand(psCommand)

	}

}

func (r RemoteComputer) GetSiteState(siteName string) string {

	psCommand := fmt.Sprintf(`Import-Module WebAdministration;
	if (Get-IISSite -Name "%s" -ErrorAction SilentlyContinue) {
		Get-IISSite "%s" | Select-Object -ExpandProperty state
		}`, siteName, siteName)

	res := r.RunCommandPlain(psCommand)

	res = strings.TrimRight(res, "\n")

	if strings.Contains(res, "Started") {

		return util.MakeColored("Started")

	} else if strings.Contains(res, "Stopped") {

		return util.MakeColored("Stopped")

	}

	return res
}

func (r RemoteComputer) StopIISSite(siteName string) {

	psCommand := fmt.Sprintf(`
		Import-Module WebAdministration;
	if (Get-IISSite -Name "%s" -ErrorAction SilentlyContinue) {
		Stop-WebSite -Name "%s"
	} 
	`, siteName, siteName)

	if r.isSiteExist(siteName) {

		msg := fmt.Sprintf(`Web Site "%s" is stopping...`, siteName)
		fmt.Println(color.HiGreenString(msg))
		r.ExecuteCommand(psCommand)

	} else {

		errorGenerator(siteName, r.ComputerName)

	}

}

func (r RemoteComputer) StartIISSite(siteName string) {

	psCommand := fmt.Sprintf(`
		Import-Module WebAdministration;
	if (Get-IISSite -Name "%s" -ErrorAction SilentlyContinue) {
		Start-WebSite -Name "%s"
	} 
	`, siteName, siteName)

	if r.isSiteExist(siteName) {

		msg := fmt.Sprintf(`Web Site "%s" is starting...`, siteName)
		fmt.Println(color.HiGreenString(msg))
		r.ExecuteCommand(psCommand)

	} else {

		errorGenerator(siteName, r.ComputerName)

	}

}

func (r RemoteComputer) ListAllSites() []model.SitesIIS {

	psCommand := `
	Import-Module WebAdministration
    $sites = Get-ChildItem -Path IIS:\Sites | ForEach-Object {
    $site = $_
    $virtualDirectories = (Get-WebVirtualDirectory -Site $site.Name).Path -replace "^/", "" -join ","
    $applications = (Get-WebApplication -Site $site.Name).Path -replace "^/", ""

    $preloadEnabled = $site.Collection.preloadEnabled
    if ($applications) {
        $preloadEnabled += (Get-WebApplication -Site $site.Name).Collection.preloadEnabled
    }

    $preloadEnabled = ($preloadEnabled | Select-Object -Unique) -join ", "
    
    $protocol = ($site.bindings.collection.protocol | Select-Object -Unique) -join ", "
    $bindings = ($site.bindings.collection.bindingInformation | Select-Object -Unique) -join ", "

    [PSCustomObject]@{
        Name = $site.name
        ApplicationPool = $site.applicationPool
        State = $site.state
        PhysicalPath = $site.physicalPath
        Bindings = $bindings
        Protocol = $protocol
        PreloadEnabled = $preloadEnabled
        VirtualDirectories = $virtualDirectories
        Applications = ($applications -join ",")
    }
}
Write-Output $sites
`

	res := r.RunCommandJSON(psCommand)
	res = strings.TrimRight(res, "\n")
	siteModel := util.JsonToStructIISSite(res)

	//If application pool remove with force command. Sites States are coming null.
	r.prepModel(siteModel)

	return setBindings(siteModel)
}

func (r RemoteComputer) WebSitesListByStatus(status string) []model.SitesIIS {

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration

	$sites = Get-ChildItem -Path IIS:\Sites  | Where-Object { $_.state -eq '%s'  } | ForEach-Object {
		$site = $_
		$virtualDirectories = (Get-WebVirtualDirectory -Site $site.Name).Path -replace "^/", "" -join ","
		$applications = (Get-WebApplication -Site $site.Name).Path -replace "^/", ""
	
		$preloadEnabled = $site.Collection.preloadEnabled
		if ($applications) {
			$preloadEnabled += (Get-WebApplication -Site $site.Name).Collection.preloadEnabled
		}
	
		$preloadEnabled = ($preloadEnabled | Select-Object -Unique) -join ", "
		
		$protocol = ($site.bindings.collection.protocol | Select-Object -Unique) -join ", "
		$bindings = ($site.bindings.collection.bindingInformation | Select-Object -Unique) -join ", "
	
		[PSCustomObject]@{
			Name = $site.name
			ApplicationPool = $site.applicationPool
			State = $site.state
			PhysicalPath = $site.physicalPath
			Bindings = $bindings
			Protocol = $protocol
			PreloadEnabled = $preloadEnabled
			VirtualDirectories = $virtualDirectories
			Applications = ($applications -join ",")
		}
	}
	Write-Output $sites
`, status)

	res := r.RunCommandJSON(psCommand)
	res = strings.TrimRight(res, "\n")
	siteModel := util.JsonToStructIISSite(res)
	r.prepModel(siteModel)

	return setBindings(siteModel)
}

func (r RemoteComputer) ListSingleSite(siteName string) []model.SitesIIS {

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration

$sites = Get-ChildItem -Path IIS:\Sites  | Where-Object { $_.Name -eq '%s' } | ForEach-Object {
    $site = $_
    $virtualDirectories = (Get-WebVirtualDirectory -Site $site.Name).Path -replace "^/", "" -join ","
    $applications = (Get-WebApplication -Site $site.Name).Path -replace "^/", ""

    $preloadEnabled = $site.Collection.preloadEnabled
    if ($applications) {
        $preloadEnabled += (Get-WebApplication -Site $site.Name).Collection.preloadEnabled
    }

    $preloadEnabled = ($preloadEnabled | Select-Object -Unique) -join ", "
    
    $protocol = ($site.bindings.collection.protocol | Select-Object -Unique) -join ", "
    $bindings = ($site.bindings.collection.bindingInformation | Select-Object -Unique) -join ", "

    [PSCustomObject]@{
        Name = $site.name
        ApplicationPool = $site.applicationPool
        State = $site.state
        PhysicalPath = $site.physicalPath
        Bindings = $bindings
        Protocol = $protocol
        PreloadEnabled = $preloadEnabled
        VirtualDirectories = $virtualDirectories
        Applications = ($applications -join ",")
    }
}
Write-Output $sites`, siteName)
	res := r.RunCommandJSON(psCommand)
	res = strings.TrimRight(res, "\n")

	siteModel := util.JsonToStructIISSite(res)
	r.prepModel(siteModel)

	return setBindings(siteModel)

}

func setBindings(siteList []model.SitesIIS) []model.SitesIIS {

	// TO DO HTTPS binding

	for i := range siteList {
		binding := siteList[i].Bindings

		if strings.HasPrefix(binding, "*:80:") {
			if binding == "*:80:" {

				siteList[i].Bindings = "http://localhost:80"

			} else {

				siteList[i].Bindings = strings.TrimPrefix(binding, "*:80:")
				siteList[i].Bindings = "http://" + strings.TrimPrefix(binding, "*:80:")
			}

		} else {
			re := regexp.MustCompile(`(\d+)`)
			matches := re.FindAllStringSubmatch(binding, -1)

			for _, match := range matches {
				siteList[i].Bindings = "http://localhost:" + match[1]
			}
		}
	}
	return siteList
}

func (r RemoteComputer) prepModel(siteModel []model.SitesIIS) {

	isIISRunning := r.IsIISRunning()

	for i, s := range siteModel {

		if s.State == "" {

			if isIISRunning {

				msg := color.HiRedString(fmt.Sprintf(`%s "%s" %s`, "Undefined State\nApplication Pool", s.ApplicationPool, "\ndoesn't exist."))
				siteModel[i].State = msg
			} else {
				magenta := color.New(color.BgHiMagenta)
				coloredStr := magenta.Sprint("IIS is stopped        \ne.g: iis-hero start")
				siteModel[i].State = coloredStr
			}
		}

		if s.PSComputerName == "" {

			siteModel[i].PSComputerName = r.ComputerName
		}

		/*if s.Applications == "" {
			siteModel[i].Applications = color.HiRedString("No applications.")
		}

		if s.VirtualDirectories == "" {

			siteModel[i].VirtualDirectories = color.HiRedString("No virtual directory.")
		}*/

	}
}

func (r RemoteComputer) ChangeWebSite(sitename string, rename string, preload string, pool string, binding string, path string) {

	if sitename == "" {
		log.Fatal(color.HiRedString("Specify a site name"))
	} else {
		if !r.isSiteExist(sitename) {

			errorGenerator(sitename, r.ComputerName)
		}
	}

	r.changeAppPool(sitename, pool)
	r.renameSiteName(sitename, rename)
	r.setPreload(sitename, preload)
	r.setBindings(sitename, binding)
	r.setPhysicalPath(sitename, path)

}

func (r RemoteComputer) setPreload(sitename string, preload string) {

	if isStrNotEmpty(preload) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
		Set-ItemProperty IIS:\Sites\%s -name applicationDefaults.preloadEnabled -value %s;`, sitename, preload)
		r.ExecuteCommand(psCommand)
	}

}

func (r RemoteComputer) renameSiteName(sitename string, rename string) {

	if isStrNotEmpty(rename) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
	Rename-Item 'IIS:\Sites\%s' '%s'`, sitename, rename)
		r.ExecuteCommand(psCommand)

	}
}

func (r RemoteComputer) changeAppPool(sitename string, pool string) {

	if isStrNotEmpty(pool) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
	Set-ItemProperty "IIS:\Sites\%s" -Name "applicationPool" -Value "%s"`, sitename, pool)
		r.ExecuteCommand(psCommand)

	}

}

func (r RemoteComputer) setBindings(sitename string, binding string) {
	if isStrNotEmpty(binding) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
		Set-ItemProperty "IIS:\Sites\%s" -Name "bindings" -Value "%s"`, sitename, binding)
		r.ExecuteCommand(psCommand)

	}

}

func (r RemoteComputer) setPhysicalPath(sitename string, path string) {

	if isStrNotEmpty(path) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
		Set-ItemProperty "IIS:\Sites\%s" -Name "physicalPath" -Value "%s"`, sitename, path)
		r.ExecuteCommand(psCommand)
	}

}

func errorGenerator(site string, computer string) {

	log.Fatal(color.HiRedString(`Web Site "%s" cannot found on Computer "%s"`, site, computer))

}
