package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abdfnx/gosh"
	"github.com/fatih/color"
)

func (r RemoteComputer) ExecuteCommand(command string) {

	if r.isLocal() {

		powershellCommand := fmt.Sprintf(`Invoke-Command  -ScriptBlock {%s};`, command)

		gosh.PowershellCommand(powershellCommand)

	} else if r.isCredentialRequired() {

		powershellCommand := fmt.Sprintf(`$username = "%s" 
			$Password = ConvertTo-SecureString -String "%s" -AsPlainText -Force 
			$Credential = [pscredential]::new($username,$Password);
			Invoke-Command -ComputerName %s -ScriptBlock {%s} -Credential $Credential;`,
			r.getUsername(),
			r.Password,
			r.ComputerName,
			command)

		gosh.PowershellCommand(powershellCommand)

	} else {

		powershellCommand := fmt.Sprintf(`Invoke-Command -ComputerName %s -ScriptBlock {%s};`, r.ComputerName, command)

		gosh.PowershellCommand(powershellCommand)
	}

}

func (r RemoteComputer) RunCommandPlain(command string) string {

	if r.isLocal() {

		powershellCommand := fmt.Sprintf(`
		Invoke-Command  -ScriptBlock {%s};
		`, command)

		err, response, stder := gosh.PowershellOutput(powershellCommand)

		if err != nil {

			log.Fatal(err, "\n", stder)
		}

		return response

	} else if r.isCredentialRequired() {

		powershellCommand := fmt.Sprintf(`$username = "%s" 
	    $Password = ConvertTo-SecureString -String "%s" -AsPlainText -Force 
	    $Credential = [pscredential]::new($username,$Password);
	    Invoke-Command -ComputerName %s -ScriptBlock {%s} -Credential $Credential;`,
			r.getUsername(),
			r.Password,
			r.ComputerName,
			command)

		err, response, stder := gosh.PowershellOutput(powershellCommand)

		if err != nil {

			log.Fatal(err, "\n", stder)
		}

		return response

	} else {

		powershellCommand := fmt.Sprintf(`
		Invoke-Command -ComputerName %s -ScriptBlock {%s};
		`, r.ComputerName, command)

		err, response, stder := gosh.PowershellOutput(powershellCommand)

		if err != nil {

			log.Fatal(err, "\n", stder)
		}

		return response

	}
}

func (r RemoteComputer) RunCommandJSON(command string) string {

	if r.isLocal() {

		powershellCommand := fmt.Sprintf(`
		try {
		$res=Invoke-Command  -ScriptBlock {%s};
		$resJson = $res | ConvertTo-Json;
		Write-Host $resJson
		}catch {
		Write-Host "Error!: $($_.Exception.Message)"
		}`, command)

		err, response, stder := gosh.PowershellOutput(powershellCommand)

		if err != nil {

			log.Fatal(err, "\n", stder)
		}

		response = strings.TrimRight(response, "\r\n")

		if response == "" {

			log.Fatal("\n" + color.HiGreenString("Could not find a result that matches your search criteria.\nIf the command has the --all flag. Try using it to list all the results \ne.g: iis-hero <your command> -a, --all"))
		} else if strings.Contains(response, "Error!") {

			log.Fatal(color.HiRedString(response))
			//log.Fatal("\n" + color.HiGreenString(response) + color.MagentaString("\nCheck your credentials. \ne.g: iis-hero login cred"))

		}

		return response

	} else if r.isCredentialRequired() {

		powershellCommand := fmt.Sprintf(`
        $username = "%s";
        $Password = ConvertTo-SecureString -String "%s" -AsPlainText -Force;
        $Credential = [pscredential]::new($username,$Password);

        try {
        $res = Invoke-Command -ComputerName %s -ScriptBlock {%s} -Credential $Credential -ErrorAction Stop
        $resJson = $res | ConvertTo-Json
        Write-Host $resJson
        } catch {
        Write-Host "Error!: $($_.Exception.Message)"
        }`, r.getUsername(),
			r.Password,
			r.ComputerName,
			command)

		err, response, stderr := gosh.PowershellOutput(powershellCommand)

		response = strings.TrimRight(response, "\r\n")

		if err != nil {
			log.Fatal(err, "\n", stderr)
		}

		if response == "" {

			log.Fatal("\n" + color.HiGreenString("Could not find a result that matches your search criteria.\nPlease try to list everything. With --all flag \ne.g:iis-hero <your command> -a, --all"))
		} else if strings.Contains(response, "Error!") {

			log.Fatal("\n" + color.HiGreenString(response) + color.MagentaString("\nCheck your credentials. \ne.g: iis-hero login cred"))

		}
		return response

	} else {
		powershellCommand := fmt.Sprintf(`
		$computerName = "%s";
			$result = Invoke-Command -ComputerName $computerName -ScriptBlock {%s}
			if($result) {
				$jsonResult = $result | ConvertTo-Json
				Write-Host $jsonResult
			} `, r.ComputerName, command)

		err, response, stder := gosh.PowershellOutput(powershellCommand)

		if err != nil {

			log.Fatal(err, "\n", stder)
		}
		response = strings.TrimRight(response, "\r\n")

		if response == "" {

			log.Fatal("\n" + color.HiGreenString("Could not find a result that matches your search criteria.\nPlease try to list everything. With --all flag \ne.g: iis <your command> -a, --all"))
		}
		return response

	}

}

func (r RemoteComputer) OpenSession(command string) {

	if r.isLocal() {

		log.Fatal("Error, do not open session on your own computer")
	} else if r.isCredentialRequired() {

		var sb strings.Builder

		baseCommand := fmt.Sprintf(`$username = "%s" 
			$Password = ConvertTo-SecureString -String "%s" -AsPlainText -Force 
			$Credential = [pscredential]::new($username,$Password);
			$computerName = "%s";
			$session = New-PSSession -ComputerName $computerName -Credential $Credential;	

		`, r.getUsername(), r.Password, r.ComputerName)

		closeSectionCommand := `
		Remove-PSSession $session`

		sb.WriteString(baseCommand)
		sb.WriteString(command)
		sb.WriteString(closeSectionCommand)

		powershellCommand := sb.String()
		gosh.PowershellCommand(powershellCommand)

	} else {

		var sb strings.Builder

		baseCommand := fmt.Sprintf(`
		$computerName = "%s"
		$session = New-PSSession -ComputerName $computerName;
		
		`, r.ComputerName)

		closeSectionCommand := "Remove-PSSession $session"

		sb.WriteString(baseCommand)
		sb.WriteString(command)
		sb.WriteString(closeSectionCommand)

		powershellCommand := sb.String()
		gosh.PowershellCommand(powershellCommand)
	}

}

func (r RemoteComputer) isCredentialRequired() bool {

	if r.ComputerName == "" {
		log.Fatal(color.HiGreenString("The '--computername, -c' flag cannot be empty.\nPlease check your credentials.\ne.g: iis-hero login cred"))
	}

	return r.DomainName != "" && r.Username != "" && r.Password != "" && r.ComputerName != ""
}

func (r RemoteComputer) isLocal() bool {

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Hostname couldn't be resolved.")
	}

	if hostname == r.ComputerName {

		if r.isCredentialRequired() {
			return false
		}
	}

	return hostname == r.ComputerName
}

func (r RemoteComputer) getUsername() string {

	var sb strings.Builder
	sb.WriteString(r.DomainName)
	sb.WriteString("\\")
	sb.WriteString(r.Username)
	return sb.String()

}
