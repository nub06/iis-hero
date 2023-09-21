package service

import (
	"fmt"
	"strings"

	"github.com/nub06/iis-hero/model"
	"github.com/nub06/iis-hero/util"
)

func (r RemoteComputer) CreateBackupDir(source string, destination string) {

	//Test-Path kontrolü eklenecek

	psCommand := fmt.Sprintf(`
	$date = (Get-Date).ToString('yyyyMMdd');
	$time = (Get-Date).ToString('hhmm');
	$getDate = -join ($date, "_", $time);
	$appDir = "%s";
	$backupDir="%s\";
	if (-not (Test-Path -Path $backupDir -PathType Container)) {
		New-Item -ItemType Directory -Path $backupDir -Name $getDate;
	}
	[string]$appSource = "$appDir\*";
	[string]$backupDestination = "$backupDir$getDate";
	Copy-item -Force -Recurse $appSource -Destination $backupDestination;`, source, destination)

	r.ExecuteCommand(psCommand)
}

func (r RemoteComputer) FolderList(source string) model.FolderList {

	psCommand := fmt.Sprintf(
		`Get-ChildItem -Path "%s" -Directory -Depth 1 | ForEach-Object {
			$folder = $_
			$hostname = hostname
			$folderSizeBytes = (Get-ChildItem $folder.FullName -Recurse -File | Measure-Object -Property Length -Sum).Sum
			$folderSizeMB = $folderSizeBytes / 1MB
			$lastAccessTime = $folder.LastAccessTime.ToString("yyyyMMdd_HHmm")
			$lastWriteTime = $folder.LastWriteTime.ToString("yyyyMMdd_HHmm")
			
				[PSCustomObject]@{
					Name = $folder.Name
					Path = $folder.FullName
					Size = "{0:N2} MB" -f $folderSizeMB
					LastAccessTime = $lastAccessTime
					LastWriteTime = $lastWriteTime
					ComputerName = $hostname
		
				
			}
		}`, source)

	res := r.RunCommandJSON(psCommand)
	res = strings.TrimRight(res, "\r\n")

	folderList := util.JsonToStructFolderList(res)

	return folderList
}

func (r RemoteComputer) RemoveFolder(path string) {

	psCommand := fmt.Sprintf(`
	$folderPath = "%s"
	Remove-Item $folderPath -Recurse -Force`, path)

	r.ExecuteCommand(psCommand)

}

func (r RemoteComputer) isSourceIsDirectory(path string) bool {

	psCommand := fmt.Sprintf(`
	$source = "%s"
	if (Test-Path $source -PathType Container) {	
	Write-Host "true"
	}else{
	Write-Host "false"}`, path)

	res := r.RunCommandPlain(psCommand)

	return strings.Contains(res, "true")

}

func (r RemoteComputer) CopyFromTarget(remotepath string, localpath string) {

	if r.isSourceIsDirectory(remotepath) {

		psCommand := fmt.Sprintf(`Copy-Item -Path "%s" -Destination "%s" -FromSession $session -Recurse;`, remotepath, localpath)

		r.OpenSession(psCommand)
	} else {
		psCommand := fmt.Sprintf(`
		$source="%s";
		$destination="%s";		
	
		if (!(Test-Path $destination -PathType Container)){
			New-Item -Path $destination -ItemType Directory -Force 
		}
	
		$fileName = Split-Path -Path $source -Leaf	
		$targetPath = Join-Path $destination $fileName	
		Copy-Item -Path $source -Destination $targetPath -FromSession $session
		
		`, remotepath, localpath)

		r.OpenSession(psCommand)
	}
}

func (r RemoteComputer) CopyToTarget(remotepath string, localpath string) {

	psCommand := fmt.Sprintf(`
	$source="%s";
	$destination="%s";	

	Invoke-Command -Session $session -ScriptBlock {
	$destination="%s";
	
	if (!(Test-Path $destination -PathType Container)){
		New-Item -Path $destination -ItemType Directory -Force 
	}}

	if (Test-Path $source -PathType Container) {	
		Copy-Item -Path $source -Destination $destination -ToSession $session -Recurse -Force
	}else{
		$fileName = Split-Path -Path $source -Leaf	
		$targetPath = Join-Path $destination $fileName	
		Copy-Item -Path $source -Destination $targetPath -ToSession $session
	}`, localpath, remotepath, remotepath)

	r.OpenSession(psCommand)
}
