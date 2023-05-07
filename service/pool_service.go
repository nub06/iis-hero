package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/model"
	"github.com/nub06/iis-hero/util"
)

func (r RemoteComputer) IsPoolExist(poolName string) bool {

	var res string
	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration;

	if (Get-IISAppPool -Name "%s" -ErrorAction SilentlyContinue) {
		Write-Host "true"
	} else {
		Write-Host "false"
	}`, poolName)

	res = r.RunCommandPlain(psCommand)
	res = strings.TrimRight(res, "\r\n")

	return res == "true"

}

func (r RemoteComputer) GetPoolInfo(poolName string) model.AppPool {

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration;
    $AppPoolName = "%s";
    $AppPool = Get-Item "IIS:\AppPools\$AppPoolName";
    $Sites = (Get-ChildItem "IIS:\Sites" | Where-Object {$_.ApplicationPool -eq $AppPoolName}).Name -join ", "
    $AppPoolObj = [PSCustomObject]@{
    Name = $AppPool.Name
    State = $AppPool.State
    AutoStart = $AppPool.AutoStart
    ManagedRuntimeVersion = $AppPool.ManagedRuntimeVersion
    ManagedPipelineMode = $AppPool.ManagedPipelineMode
    StartMode = $AppPool.StartMode
    IdleTimeout = $AppPool.ProcessModel.IdleTimeout.TotalMinutes
    IdleTimeoutAction = $AppPool.ProcessModel.IdleTimeoutAction
    Sites = $Sites
   } 
   Write-Output $AppPoolObj	
	`, poolName)

	var res string

	if r.IsPoolExist(poolName) {

		res = r.RunCommandJSON(psCommand)
	} else {

		log.Fatalf(color.HiRedString(`Application Pool  "%s" couldn't find on computer "%s"`), poolName, r.ComputerName)
	}

	m := util.JsonToStructAppPool(res)

	return m[0]

}

func (r RemoteComputer) GetAppPoolState(appName string) string {

	psCommand := fmt.Sprintf(`Get-IISAppPool "%s" | Select-Object -ExpandProperty state`, appName)

	res := r.RunCommandPlain(psCommand)

	res = strings.TrimRight(res, "\n")

	if strings.Contains(res, "Started") {

		return util.MakeColored("Started")

	} else if strings.Contains(res, "Stopped") {

		return util.MakeColored("Stopped")
	}

	return (color.HiCyanString("Error"))
}

func (r RemoteComputer) StopWebAppPool(appName string) {

	if appName == "" {

		log.Fatal("Application Pool name cannot be empty.")
	}

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration;
	Stop-WebAppPool -Name "%s"`, appName)

	r.ExecuteCommand(psCommand)

}

func (r RemoteComputer) StartWebAppPool(appName string) {

	if appName == "" {

		log.Fatal("Application Pool name cannot be empty.")
	}

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration;
	Start-WebAppPool -Name "%s"`, appName)

	r.ExecuteCommand(psCommand)

}

func (r RemoteComputer) RestartWebAppPool(appName string, all bool) {

	if appName == "" && !all {

		log.Fatal(color.HiGreenString("Application Pool name can not be empty. Please specify a Application Pool name or use --all flag.\ne.g: iis-hero pool restart <pool>\niis-hero pool restart -a"))
	}

	var psCommand string

	if all {

		psCommand = `
		Import-Module WebAdministration;
		Get-ChildItem IIS:\AppPools | ForEach-Object { Restart-WebAppPool $_.Name }
		`

		fmt.Println(color.HiGreenString("All Application Pools are restarting..."))

	} else {

		psCommand = fmt.Sprintf(`
		Import-Module WebAdministration;
		Restart-WebAppPool -Name "%s"`, appName)

		msg := fmt.Sprintf(`Application Pool "%s" is restarting...`, appName)
		fmt.Println(color.HiGreenString(msg))
	}

	r.ExecuteCommand(psCommand)
}

func (r RemoteComputer) RemoveWebAppPool(isForce bool, appName string) {

	if appName == "" {

		log.Fatal(color.HiRedString("Application Pool name cannot be empty."))
	}

	psCommand := fmt.Sprintf(`Remove-WebAppPool -Name "%s"`, appName)

	poolInfo := r.GetPoolInfo(appName)

	siteInfo := poolInfo.Sites

	if isForce {
		fmt.Printf(color.HiGreenString(`The Application Pool "%s" is removing...`), appName)

		r.ExecuteCommand(psCommand)

	} else {
		if siteInfo != "" {

			//91m is coming from HiRedString
			if strings.Contains(siteInfo, "91m") {
				fmt.Printf(color.HiGreenString(`The Application Pool "%s" is removing...`), appName)

				r.ExecuteCommand(psCommand)

			} else {

				var sb strings.Builder
				message := fmt.Sprintf(`The Application Pool "%s" includes sites: "%s" `, appName, poolInfo.Sites)
				sb.WriteString(color.HiRedString(message))
				sb.WriteString(color.HiRedString("if you still want to delete this pool "))
				sb.WriteString(color.HiGreenString("\ne.g:\niis-hero pool remove <poolname> -f\niis-hero pool remove <poolname> --force"))
				log.Fatal(sb.String())
			}

		}
	}

}

func (r RemoteComputer) ListAllAppPools() []model.AppPool {

	psCommand := `
	Import-Module WebAdministration;
    $AppPools = Get-ChildItem IIS:\AppPools;
    $AppPoolInfo = @()
    foreach ($AppPool in $AppPools) {
    $Sites = (Get-ChildItem "IIS:\Sites" | Where-Object {$_.ApplicationPool -eq $AppPool.Name}).Name -join ", "
    $AppPoolObj = [PSCustomObject]@{
        Name = $AppPool.Name
        State = $AppPool.State
        AutoStart = $AppPool.AutoStart
        ManagedRuntimeVersion = $AppPool.ManagedRuntimeVersion
        ManagedPipelineMode = $AppPool.ManagedPipelineMode
        StartMode = $AppPool.StartMode
        IdleTimeout = $AppPool.ProcessModel.IdleTimeout.TotalMinutes
        IdleTimeoutAction = $AppPool.ProcessModel.IdleTimeoutAction
        Sites = $Sites
    }

    $AppPoolInfo += $AppPoolObj
    }

   Write-Output $AppPoolInfo`

	res := r.RunCommandJSON(psCommand)
	res = strings.TrimRight(res, "\n")

	apps := util.JsonToStructAppPool(res)

	return apps
}

func (r RemoteComputer) ListAppPoolsByStatus(status string) []model.AppPool {

	psCommand := fmt.Sprintf(`Import-Module WebAdministration;
    $AppPools = Get-ChildItem IIS:\AppPools;
    $AppPoolInfo = @()
    foreach ($AppPool in $AppPools) {
        if($AppPool.State -eq "%s"){
            $Sites = (Get-ChildItem "IIS:\Sites" | Where-Object {$_.ApplicationPool -eq $AppPool.Name}).Name -join ", "
            $AppPoolObj = [PSCustomObject]@{
                Name = $AppPool.Name
                State = $AppPool.State
                AutoStart = $AppPool.AutoStart
                ManagedRuntimeVersion = $AppPool.ManagedRuntimeVersion
                ManagedPipelineMode = $AppPool.ManagedPipelineMode
                StartMode = $AppPool.StartMode
                IdleTimeout = $AppPool.ProcessModel.IdleTimeout.TotalMinutes
                IdleTimeoutAction = $AppPool.ProcessModel.IdleTimeoutAction
                Sites = $Sites
            }
            $AppPoolInfo += $AppPoolObj
        }
    }
    Write-Output $AppPoolInfo
    `, status)

	res := r.RunCommandJSON(psCommand)
	res = strings.TrimRight(res, "\n")

	apps := util.JsonToStructAppPool(res)

	return apps
}

func (r RemoteComputer) AppPoolCreate(poolName string, startmode string, autostart string, runtime string,
	pipelinemode string, idleaction string, idleminute string) {

	if strings.Contains(runtime, "v4") {

		runtime = "v4.0"
	} else if strings.Contains(runtime, "v2") {

		runtime = "v2.0"

	}

	psCommand := fmt.Sprintf(`
	Import-Module WebAdministration;
	$poolName="%s"
	New-Item IIS:\AppPools\$poolName;
	Set-ItemProperty IIS:\AppPools\$poolName -Name "autoStart" -Value "%s";
	Set-ItemProperty IIS:\AppPools\$poolName -Name "startMode" -Value "%s";
	Set-ItemProperty IIS:\AppPools\$poolName -Name "managedRuntimeVersion" -Value "%s";
	Set-ItemProperty IIS:\AppPools\$poolName -Name "managedPipelineMode" -Value "%s";
	Set-ItemProperty IIS:\AppPools\$poolName -Name "processModel.idleTimeout" -Value "%s";
	Set-ItemProperty IIS:\AppPools\$poolName -Name "processModel.idleTimeoutAction" -Value "%s";
	`, poolName, autostart, startmode, runtime, pipelinemode, idleminute, idleaction)

	if r.IsPoolExist(poolName) {

		message := fmt.Sprintf(`Application Pool "%s" is already exist...`, poolName)
		fmt.Println(color.HiRedString(message))
		util.MakeTableFromStruct(r.GetPoolInfo(poolName))

	} else {

		message := fmt.Sprintf(`Application Pool "%s" is creating...`, (poolName))
		fmt.Println(color.HiGreenString(message))
		r.ExecuteCommand(psCommand)
	}

}

func (r RemoteComputer) AppPoolChange(pool string, runtime string, startmode string, autostart string,
	idleaction string, idleminute string, pipelinemode string, rename string) {

	if pool == "" {
		log.Fatal(color.HiRedString("Specify a pool name"))
	} else {
		if !r.IsPoolExist(pool) {
			log.Fatal(color.HiRedString(`Application Pool "%s" cannot found`, pool))
		}
	}
	r.setRuntimeVersion(pool, runtime)
	r.setPipelineMode(pool, pipelinemode)
	r.setAutoStart(pool, autostart)
	r.setStartMode(pool, startmode)
	r.setIdleTimeAction(pool, idleaction)
	r.setIdleTimeMinute(pool, idleminute)
	r.renameAppPool(pool, rename)
}

func (r RemoteComputer) renameAppPool(pool string, rename string) {

	if isStrNotEmpty(rename) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
		Rename-Item 'IIS:\AppPools\%s' '%s'`, pool, rename)
		r.ExecuteCommand(psCommand)
	}
}

func (r RemoteComputer) setRuntimeVersion(pool string, runtime string) {

	if isStrNotEmpty(runtime) {
		if runtime == "0" {
			runtime = ""
		} else if strings.Contains(runtime, "v4") {
			runtime = "v4.0"
		} else if strings.Contains(runtime, "v2") {
			runtime = "v2.0"
		}
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
		Set-ItemProperty IIS:\AppPools\"%s" -Name "managedRuntimeVersion" -Value "%s"`, pool, runtime)
		r.ExecuteCommand(psCommand)
	}

}

func (r RemoteComputer) setStartMode(pool string, startmode string) {

	if isStrNotEmpty(startmode) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
		Set-ItemProperty IIS:\AppPools\"%s" -Name "startMode" -Value "%s"`, pool, startmode)
		r.ExecuteCommand(psCommand)
	}

}

func (r RemoteComputer) setAutoStart(pool string, autostart string) {

	if isStrNotEmpty(autostart) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
		Set-ItemProperty IIS:\AppPools\"%s" -Name "autoStart" -Value "%s"`, pool, autostart)
		r.ExecuteCommand(psCommand)
	}

}

func (r RemoteComputer) setIdleTimeAction(pool string, idleaction string) {

	if isStrNotEmpty(idleaction) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
	Set-ItemProperty IIS:\AppPools\"%s" -Name "processModel.idleTimeoutAction" -Value "%s"`, pool, idleaction)
		r.ExecuteCommand(psCommand)
	}
}
func (r RemoteComputer) setIdleTimeMinute(pool string, idleminute string) {

	if isStrNotEmpty(idleminute) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
	Set-ItemProperty IIS:\AppPools\"%s" -Name "processModel.idleTimeout" -Value "00:%s"`, pool, idleminute)
		r.ExecuteCommand(psCommand)
	}

}
func (r RemoteComputer) setPipelineMode(pool string, pipelinemode string) {

	if isStrNotEmpty(pipelinemode) {
		psCommand := fmt.Sprintf(`Import-Module WebAdministration;
	Set-ItemProperty IIS:\AppPools\"%s" -Name "managedPipelineMode" -Value "%s"`, pool, pipelinemode)
		r.ExecuteCommand(psCommand)
	}
}

func isStrNotEmpty(str string) bool {
	return str != ""
}
