package service

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/util"
)

type RemoteComputer struct {
	ComputerName string
	Username     string
	Password     string
	DomainName   string
}

func (r RemoteComputer) IsIISRunning() bool {

	psCommand := "iisreset /status"
	status := r.RunCommandPlain(psCommand)

	return strings.Contains(status, "Running")
}

func (r RemoteComputer) StartIIS() {
	psCommand := "iisreset /start"
	r.ExecuteCommand(psCommand)

}

func (r RemoteComputer) StopIIS() {
	psCommand := "iisreset /stop"
	r.ExecuteCommand(psCommand)
}

func (r RemoteComputer) ResetIIS() {
	psCommand := "iisreset"
	r.ExecuteCommand(psCommand)
}

func (r RemoteComputer) BackupIISConfig(foldername string) {

	psCommand := fmt.Sprintf(`
	$specialFolder = "%s"
	$date=(Get-Date -Format 'yyyyMMdd_HHmm')
    if ($specialFolder -eq "") {
    $folderName = -join("IISConfigBackup_", $date)
    }else {
    $folderName = $specialFolder
    }
    Backup-WebConfiguration -Name $folderName
	`, foldername)

	fmt.Println(r.RunCommandPlain(psCommand))
}

func (r RemoteComputer) RestoreIISConfig(isLatest bool, backupFolderName string) {

	if backupFolderName == "" && !isLatest {

		log.Fatal(color.HiRedString("You have to specify backup folder name to restore"))
	}

	backupFolderName = filepath.Base(backupFolderName)

	if isLatest {
		psCommand := `
		Get-ChildItem C:\Windows\System32\inetsrv\backup| Where-Object { $_.PSIsContainer } | Sort-Object CreationTime -Descending | Select-Object -First 1 | Select-Object -ExpandProperty name
		`
		latestFolderName := r.RunCommandPlain(psCommand)
		backupFolderName = strings.TrimRight(latestFolderName, "\r\n")
	}

	fmt.Println(color.HiGreenString("Restoring IIS configuration to %s", backupFolderName))
	psCommand := fmt.Sprintf(`Restore-WebConfiguration -Name %s
	Restart-Service W3SVC`, backupFolderName)
	r.ExecuteCommand(psCommand)
}

func (r RemoteComputer) ConfigBackupLists() {

	util.MakeTable(r.FolderList("C:\\Windows\\System32\\inetsrv\\backup"))
}

func (r RemoteComputer) ConfigClearAll() {

	psCommand := `
	Import-Module WebAdministration
	Clear-WebConfiguration -Recurse
	Remove-Item -Path "IIS:\Sites\*" -Recurse
	Remove-WebAppPool -Name *`

	res := r.RunCommandPlain(psCommand)

	fmt.Println(res)

}

func (r RemoteComputer) RemoveIISBackup(folder string, isForce bool) {

	folder = filepath.Base(folder)
	var folderPath string
	if isForce {
		if folder == "" {
			folderPath = "C:\\Windows\\System32\\inetsrv\\backup\\*"
		} else {

			log.Fatal(color.HiRedString(`Can not use a folder name and --all flag together
e.g: iis-hero config backup remove <foldername>
iis-hero config backup remove -a`))
		}

	} else {
		folderPath = "C:\\Windows\\System32\\inetsrv\\backup" + "\\" + folder
	}
	r.RemoveFolder(folderPath)
}
