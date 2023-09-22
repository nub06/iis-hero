package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/abdfnx/gosh"
	"github.com/fatih/color"
	"github.com/nub06/iis-hero/model"
	"github.com/nub06/iis-hero/util"
)

func SaveConfig(tag string) {
	name := tag + ".yaml"
	source := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "config.yaml")
	destDir := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "saved_objects")
	destination := filepath.Join(destDir, name)
	psCommand := fmt.Sprintf(`
	try{
	$sourcePath = "%s"
	$destinationPath = "%s"
	$destinationFolder = "%s"
	$reName = "saved_%s"
	$isExist = Join-Path -Path $destinationFolder -ChildPath $reName
	
	if (-not (Test-Path -Path $destinationFolder -PathType Container)) {
		New-Item -Path $destinationFolder -ItemType Directory
	}

	Copy-Item -Path $sourcePath -Destination $destinationPath


	if (Test-Path -Path $isExist -PathType Leaf) {

		Remove-Item -Path $isExist -Force
	}
	
	Rename-Item -Path $destinationPath -NewName $reName
	Write-Output "Success"
	}catch{	
	Write-Host "Error"
	}
	
`, source, destination, destDir, name)

	_, res, _ := gosh.PowershellOutput(psCommand)

	if strings.Contains(res, "Success") {

		fmt.Println(color.HiGreenString("Current credentials has been saved as a configuration profile '%s' successfully.", tag))
	} else {

		fmt.Println(color.HiRedString("Failed to save the configuration profile: '%s'.", tag))

	}
}

func UseConfig(tag string) {

	conf := "saved_" + tag + ".yaml"
	source := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "saved_objects", conf)
	destination := filepath.Join(os.Getenv("APPDATA"), "iis-hero", conf)

	configPath := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "config.yaml")

	currentConfigPath := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "saved_objects", "current")
	currentConfig := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "saved_objects", "current", conf)

	psCommand := fmt.Sprintf(`
	$sourcePath = "%s"
	$destinationPath = "%s"	
	$configPath = "%s"
	$currentConfigPath = "%s"
	$currentConfig = "%s"

	if (Test-Path -Path $sourcePath) {
	Copy-Item -Path $sourcePath -Destination $destinationPath

	if (-not (Test-Path -Path $currentConfigPath -PathType Container)) {
		New-Item -Path $currentConfigPath -ItemType Directory
	} 

	Remove-Item -Path $currentConfigPath\* -Recurse -Force
	
	Copy-Item -Path $sourcePath -Destination $currentConfig

	if (Test-Path -Path $configPath) {
		Remove-Item -Path $configPath -Force
	}
	
	Rename-Item -Path $destinationPath -NewName "config.yaml"
}else {
    Write-Host "Error"
}
`, source, destination, configPath, currentConfigPath, currentConfig)

	_, resp, _ := gosh.PowershellOutput(psCommand)

	if strings.Contains(resp, "Error") {

		log.Fatalf(color.HiRedString("Configuration profile not found: %s\nYou can try listing all saved configuration profiles, e.g.:  %s", tag, color.HiGreenString("\niis-hero profile ls")))

	} else {

		fmt.Println(color.HiGreenString("Configuration profile change successful.\nCurrent Profile: '%s'.", tag))
	}

}

func ShowSavedConfigs() (error, model.ConfigInfo) {

	source := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "saved_objects")

	psCommand := fmt.Sprintf(`

	$dirPath = "%s"

$files = Get-ChildItem -Path $dirPath -File -Filter "saved_*.yaml" | ForEach-Object {

	$lastAccessTime = $_.LastAccessTime.ToString("yyyyMMdd_HHmm")
	$lastWriteTime = $_.LastWriteTime.ToString("yyyyMMdd_HHmm")

    [PSCustomObject]@{
        Profile = $_.Name -replace '^saved_(.*?)\.yaml$', '$1'
		lastWriteTime = $LastWriteTime
		lastAccessTime = $lastAccessTime

    }
}

$files | ConvertTo-Json
	
`, source)

	err, response, stder := gosh.PowershellOutput(psCommand)

	if err != nil {

		log.Fatal(err, "\n", stder)
	}

	response = strings.TrimRight(response, "\r\n")

	error, configList := util.JsonStructToConfigInfo(response)

	return error, configList

}

func ShowCurrentConfig() (error, model.ConfigInfo) {

	source := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "saved_objects", "current")

	psCommand := fmt.Sprintf(`

	$dirPath = "%s"

$files = Get-ChildItem -Path $dirPath -File -Filter "saved_*.yaml" | ForEach-Object {

	$lastAccessTime = $_.LastAccessTime.ToString("yyyyMMdd_HHmm")
	$lastWriteTime = $_.LastWriteTime.ToString("yyyyMMdd_HHmm")

    [PSCustomObject]@{
        Profile = $_.Name -replace '^saved_(.*?)\.yaml$', '$1'
		lastWriteTime = $LastWriteTime
		lastAccessTime = $lastAccessTime

    }
}

$files | ConvertTo-Json
	
`, source)

	err, response, stder := gosh.PowershellOutput(psCommand)

	if err != nil {

		log.Fatal(err, "\n", stder)
	}

	response = strings.TrimRight(response, "\r\n")
	err, configList := util.JsonStructToConfigInfo(response)
	return err, configList
}

func RemoveCurrentConf() {

	path := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "saved_objects", "current")

	psCommand := fmt.Sprintf(`	
	$sourcePath = "%s"	

	if (Test-Path -Path $sourcePath -PathType Container) {
		Remove-Item -Path $sourcePath -Recurse -Force
	}
`, path)

	gosh.PowershellCommand(psCommand)

}

func DeleteConfiguration(tag string, isAll bool) {

	source := filepath.Join(os.Getenv("APPDATA"), "iis-hero", "saved_objects")

	var psCommand string

	if isAll {

		psCommand = fmt.Sprintf(`
	
		$dirPath = "%s"
		if (Test-Path -Path $dirPath -PathType Container) {
		Remove-Item -Path $dirPath -Force -Recurse
		Write-Output "Success"
		}else{
			Write-Output "Error"
		}`, source)

	} else {

		conf := fmt.Sprintf("saved_%s.yaml", tag)

		source = filepath.Join(source, conf)

		psCommand = fmt.Sprintf(`
		$confPath= "%s"	
		if (Test-Path -Path $confPath  -PathType Leaf) {
			Remove-Item -Path $confPath -Force
			Write-Output "Success"
			}else{
				Write-Output "Error"
			}`, source)

	}

	_, res, _ := gosh.PowershellOutput(psCommand)

	if strings.Contains(res, "Error") {

		if !isAll {

			fmt.Println(color.HiRedString("Cannot find Configuration Profile '%s'.", tag))
		} else {
			fmt.Println(color.HiRedString("Cannot find any Configuration Profile."))

		}
	} else {

		if isAll {
			fmt.Println(color.HiGreenString("Deleting all Configuration Profiles..."))
		} else {
			fmt.Println(color.HiGreenString("Deleting Configuration Profile: '%s'...", tag))
		}
	}

}
