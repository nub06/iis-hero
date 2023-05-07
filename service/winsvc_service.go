package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/model"
	"github.com/nub06/iis-hero/util"
)

func (r RemoteComputer) isServiceExist(serviceName string) bool {

	psCommand := fmt.Sprintf(`
	$name = "%s"
    $service = Get-CimInstance -ClassName Win32_Service -Filter "Name='$name' OR DisplayName='$name'"
	if ($service) {
		Write-Host "True"
	} else {
		Write-Host "False"
	}`, serviceName)

	res := r.RunCommandPlain(psCommand)

	res = strings.TrimRight(res, "\r\n")

	return res == "True"

}

func (r RemoteComputer) GetWindowsServiceStats(serviceName string) model.WindowsService {

	psCommand := fmt.Sprintf(`
	$serviceName = "%s"
    $service = Get-CimInstance -ClassName Win32_Service | ?{$_.Name -like $serviceName -or $_.DisplayName -like $serviceName}
    if ($service) {
    $serviceDetails = $service | Select-Object Name, DisplayName, Description, State, StartMode, @{Name='Path';Expression={$_.PathName |Split-Path -Parent}}, @{Name='ExecutablePath';Expression={$_.PathName}}
    Write-Output $serviceDetails
    } else {
    Write-Host "Error"
    }`, serviceName)

	res := r.RunCommandJSON(psCommand)

	if strings.Contains(res, "Error") {

		message := fmt.Sprintf("Service %s couldn't found on computer %s", serviceName, r.ComputerName)

		log.Fatal(color.HiRedString(message))

	}
	resp := util.JsonStructToWindowsService(res)

	return resp

}

func (r RemoteComputer) StartWindowsService(serviceName string) {

	_, isStarted := r.GetWindowsServiceState(serviceName)

	if strings.Contains(isStarted, "Running") {

		log.Fatalf(fmt.Sprintf("%s%s%s", color.HiGreenString("The service -> "), color.HiCyanString(serviceName), color.HiGreenString(" is already running. ")))

	}

	psCommand := fmt.Sprintf(`
	$serviceName = "%s"
    if (Get-Service -Name $serviceName -ErrorAction SilentlyContinue) {
    Start-Service -Name $serviceName
	Write-Host "OK"
    } elseif (Get-Service -DisplayName $serviceName -ErrorAction SilentlyContinue) {
    Start-Service -DisplayName $serviceName
	Write-Host "OK"
    } else {
    Write-Host "Error"
    }`, serviceName)

	res := r.RunCommandPlain(psCommand)

	responseGenerate(res, serviceName, "Starting...")

}

func (r RemoteComputer) StopWindowsService(serviceName string) {

	_, isStopped := r.GetWindowsServiceState(serviceName)

	if strings.Contains(isStopped, "Stopped") {

		log.Fatalf(fmt.Sprintf("%s%s%s", color.HiRedString("The service -> "), color.HiCyanString(serviceName), color.HiRedString(" is already stopped. ")))

	}

	psCommand := fmt.Sprintf(`
	$serviceName = "%s"
    if (Get-Service -Name $serviceName -ErrorAction SilentlyContinue) {
    Stop-Service -Name $serviceName
	Write-Host "OK"
    } elseif (Get-Service -DisplayName $serviceName -ErrorAction SilentlyContinue) {
    Stop-Service -DisplayName $serviceName
	Write-Host "OK"
    } else {
    Write-Host "Error"
    }`, serviceName)

	res := r.RunCommandPlain(psCommand)

	responseGenerate(res, serviceName, "Stopping...")

}

func (r RemoteComputer) RestartWindowsService(serviceName string) {

	psCommand := fmt.Sprintf(`Restart-Service -Name "%s"`, serviceName)
	r.ExecuteCommand(psCommand)

}

func (r RemoteComputer) GetWindowsServiceState(serviceName string) (string, string) {

	psCommand := fmt.Sprintf(`
	$serviceName = "%s"
    if (Get-Service -Name $serviceName -ErrorAction SilentlyContinue) {
    $serviceStatus = Get-Service -Name $serviceName | Select-Object -ExpandProperty status
    Write-Host $serviceStatus
	}elseif (Get-Service | Where-Object { $_.DisplayName -eq $serviceName } -ErrorAction SilentlyContinue) {
    $serviceStatus = Get-Service | Where-Object { $_.DisplayName -eq $serviceName } | Select-Object -ExpandProperty status
    Write-Host $serviceStatus
    } else {
    Write-Host "Error" 
   }`, serviceName)

	res := r.RunCommandPlain(psCommand)
	res = strings.TrimRight(res, "\n")

	responseErrorCheck(serviceName, res)

	var sb strings.Builder
	sb.WriteString(color.HiCyanString("Windows Service: "))
	sb.WriteString(color.HiBlueString(serviceName))
	sb.WriteString(color.HiCyanString(" -> "))
	sb.WriteString(util.MakeColored(res))

	return sb.String(), res

}

func (r RemoteComputer) ChangeStartMode(serviceName string, startMode string) {

	var psCommand string

	if isStrNotEmpty(startMode) {

		startMode = strings.ToLower(startMode)

		if strings.Contains(startMode, "auto") {

			psCommand = fmt.Sprintf(`Set-Service -Name "%s" -StartupType Auto`, serviceName)

		} else if strings.Contains(startMode, "manual") {

			psCommand = fmt.Sprintf(`Set-Service -Name "%s" -StartupType Manual`, serviceName)

		} else if strings.Contains(startMode, "disabled") {

			psCommand = fmt.Sprintf(`Set-Service -Name "%s" -StartupType Disabled`, serviceName)

		}

		fmt.Println(color.HiGreenString(`Changing start mode of Windows Service "%s"`, serviceName))

		r.ExecuteCommand(psCommand)
	}

}

func responseGenerate(response string, serviceName, status string) string {

	responseErrorCheck(serviceName, response)

	var sb strings.Builder
	sb.WriteString(color.HiCyanString("Windows Service: "))
	sb.WriteString(color.HiBlueString(serviceName))
	sb.WriteString(color.HiCyanString(" -> is "))
	sb.WriteString(color.HiCyanString(status))

	if status != "" {
		fmt.Println(sb.String())
	}

	return sb.String()

}

func (r RemoteComputer) DeleteWinSvc(serviceName string) {

	if serviceName == "" {
		log.Fatal(color.HiRedString("Specify a service name"))
	} else {
		if !r.isServiceExist(serviceName) {
			log.Fatal(color.HiRedString("Windows Service -> %s cannot found", serviceName))
		}
	}

	psCommand := fmt.Sprintf(`
	$name = "%s"
    $service = Get-CimInstance -ClassName Win32_Service -Filter "Name='$name' OR DisplayName='$name'"
    if ($service) {
	Stop-Service $name
	Remove-CimInstance $service
    } else {
    Write-Host "Error"
    }`, serviceName)

	fmt.Println(color.HiGreenString(`Windows Service "%s" is removing...`, serviceName))
	res := r.RunCommandPlain(psCommand)

	responseErrorCheck(serviceName, res)

}

func (r RemoteComputer) CreateWinsvc(serviceName string, displayName string, description string, exePath string, startupType string) {

	fmt.Println("Exepath:", exePath)

	if serviceName != "" && exePath != "" {

		if description == "Description of <service name>" {

			description = fmt.Sprintf("Description of %s", serviceName)

		}

		if displayName == "<service name>" {

			displayName = serviceName
		}
		psCommand := fmt.Sprintf(`New-Service -Name "%s" -DisplayName "%s" -Description "%s" -BinaryPathName "%s" -StartupType "%s"`, serviceName, displayName, description, exePath, startupType)

		//r.ExecuteCommand(psCommand)

		s := r.RunCommandPlain(psCommand)

		fmt.Println(s)

	} else {

		log.Fatal(color.HiRedString("Please specify service name and executable path for create windows service.\ne.g: iis-hero winsvc create <service name> --exepath <executable path>"))
	}

}

func (r RemoteComputer) ChangeWinSvc(serviceName string, displayName string, description string, exePath string, startupType string) {

	if serviceName == "" {
		log.Fatal(color.HiRedString("Specify a service name"))
	} else {
		if !r.isServiceExist(serviceName) {
			log.Fatal(color.HiRedString("Windows Service -> %s cannot found", serviceName))
		}
	}

	r.ChangeStartMode(serviceName, startupType)
	r.setDisplayName(serviceName, displayName)
	r.setExePath(serviceName, exePath)
	r.setDescription(serviceName, description)
}

func (r RemoteComputer) setDescription(serviceName string, description string) {

	if isStrNotEmpty(description) {

		psCommand := fmt.Sprintf(`
	$name = "%s"
    $service = Get-CimInstance -Class Win32_Service -Filter "Name='$name' OR DisplayName='$name'"
	if ($service) {
		Set-Service -Name "%s" -Description "%s"
	} else {
		Write-Host "Error"
	}`, serviceName, serviceName, description)

		fmt.Println(color.HiGreenString(`Changing Description of Windows Service "%s"`, serviceName))

		res := r.RunCommandPlain(psCommand)

		responseErrorCheck(serviceName, res)

	}

}

func (r RemoteComputer) setDisplayName(serviceName string, displayName string) {

	if isStrNotEmpty(displayName) {

		psCommand := fmt.Sprintf(`
	$name = "%s"
    $service = Get-CimInstance -Class Win32_Service -Filter "Name='$name' OR DisplayName='$name'"
	if ($service) {
		Set-Service -Name $name -DisplayName "%s"
	} else {
		Write-Host "Error"
	}`, serviceName, displayName)

		fmt.Println(color.HiGreenString(`Changing DisplayName of Windows Service "%s"`, serviceName))

		res := r.RunCommandPlain(psCommand)

		responseErrorCheck(serviceName, res)

	}
}

func (r RemoteComputer) setExePath(serviceName string, exePath string) {

	if isStrNotEmpty(exePath) {

		_, serviceState := r.GetWindowsServiceState(serviceName)

		psCommand := fmt.Sprintf(`
	$name = "%s"
    $service = Get-CimInstance -Class Win32_Service -Filter "Name='$name' OR DisplayName='$name'"
	if ($service) {
		$currentPath = (Get-WmiObject -Class Win32_Service -Filter"Name='$name' OR DisplayName='$name'").PathName
        $newPath = "%s"
        Stop-Service $name
        Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\$name" -Name "ImagePath" -Value $newPath
	} else {
		Write-Host "Error"
	}`, serviceName, exePath)

		if serviceState == "Running" {

			r.StartWindowsService(serviceName)

		}

		fmt.Println(color.HiGreenString(`Changing Executable Path of Windows Service "%s"`, serviceName))

		res := r.RunCommandPlain(psCommand)

		responseErrorCheck(serviceName, res)

	}
}

func responseErrorCheck(serviceName, response string) {

	if strings.Contains(response, "Error") {

		serviceName := color.HiCyanString(serviceName)

		var sb strings.Builder
		sb.WriteString(color.HiRedString("Windows service: "))
		sb.WriteString(serviceName)
		sb.WriteString(color.HiRedString(" couldn't find on computer "))
		log.Fatal(sb.String())

	}
}
