package util

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/model"

	"log"
)

func JsonToStructAppPool(resp string) []model.AppPool {

	data := []byte(resp)

	var m []model.AppPool

	err := json.Unmarshal(data, &m)

	if err != nil {

		//Generating dummy json data

		dummyJson := `
		{
			"name":  "null",
			"state":  "null",
			"autoStart":  false,
			"managedRuntimeVersion":  "null",
			"managedPipelineMode":  "null",
			"startMode":  "null",
			"idletimeout":  0,
			"idleTimeoutAction":  "null"
		}
		`
		dummyData := []byte("[" + resp + "," + dummyJson + "]")

		err := json.Unmarshal(dummyData, &m)

		if err != nil {

			log.Fatal(err)
		}

		for i, s := range m {

			if s.Name == "null" || s.State == "null" {

				m = append(m[:i], m[i+1:]...)
			}
		}

	}

	for i := range m {

		if m[i].ManagedRuntimeVersion == "" {
			m[i].ManagedRuntimeVersion = "No Managed Code"
		}
		if m[i].State == "" {

			magenta := color.New(color.BgHiMagenta)
			coloredStr := magenta.Sprint("IIS is stopped\ne.g: iis-hero start")
			m[i].State = coloredStr
		}

		if m[i].PSComputerName == "" {

			m[i].PSComputerName, _ = os.Hostname()
		}

		if m[i].Sites == "" {

			var sb strings.Builder
			bg := color.New(color.FgHiRed)
			sb.WriteString(bg.Sprint("There are no websites "))
			sb.WriteString(bg.Sprint("running on this "))
			sb.WriteString(bg.Sprint("Application Pool "))
			m[i].Sites = sb.String()
		}

	}

	return m

}

func JsonToStructIISSite(resp string) []model.SitesIIS {

	data := []byte(resp)

	var m []model.SitesIIS

	err := json.Unmarshal(data, &m)

	if err != nil {

		dummyJson := `
		{
			"name":  "null",
			"applicationPool":  "null",
			"state":  "null",
			"physicalPath":  "null",
			"bindings":  "null",
			"protocol":  "null"	
		}
		`
		dummyData := []byte("[" + resp + "," + dummyJson + "]")

		err := json.Unmarshal(dummyData, &m)

		if err != nil {

			log.Fatal(err)
		}

		for i, s := range m {

			if s.Name == "null" || s.ApplicationPool == "null" {

				m = append(m[:i], m[i+1:]...)
			}
		}
	}

	return m

}

func JsonStructToWindowsService(resp string) model.WindowsService {

	data := []byte(resp)

	var m model.WindowsService

	err := json.Unmarshal(data, &m)

	if err != nil {

		log.Print(err)

	}

	if m.PSComputerName == "" {

		m.PSComputerName, _ = os.Hostname()

	}
	//path := strings.Split(m.ExecutablePath, " ")
	//m.ExecutablePath = path[0]

	return m

}

func JsonToStructFolderList(resp string) model.FolderList {

	data := []byte(resp)

	var m model.FolderList

	err := json.Unmarshal(data, &m)

	if err != nil {

		dummyJson := `
				{
					"Name":  "null",
					"Path":  "null",
					"SIZE":  "null",
					"LASTACCESSTIME":  "null"
				}
				`
		dummyData := []byte("[" + resp + "," + dummyJson + "]")

		err := json.Unmarshal(dummyData, &m)

		if err != nil {

			log.Fatal(err)
		}

		for i, s := range m {

			if s.Name == "null" || s.Path == "null" {

				m = append(m[:i], m[i+1:]...)
			}

		}

	}

	return m

}

func JsonToStructVirtualDir(resp string) []model.VirtualDir {

	data := []byte(resp)

	var m []model.VirtualDir

	err := json.Unmarshal(data, &m)

	if err != nil {

		//Generating dummy json data

		dummyJson := `
				{
					"Website":  "null",
					"VirtualDirectory":  "null",
					"PhysicalPath":  "null",
					"ApplicationPool":  "null"
				}
				`
		dummyData := []byte("[" + resp + "," + dummyJson + "]")

		err := json.Unmarshal(dummyData, &m)

		if err != nil {

			log.Fatal(err)
		}

		for i, s := range m {

			if s.VirtualDirectory == "null" || s.PhysicalPath == "null" {

				m = append(m[:i], m[i+1:]...)
			}

		}

	}

	for i, s := range m {

		if s.PSComputerName == "" {

			hostname, _ := os.Hostname()
			m[i].PSComputerName = hostname
		}
	}

	return m

}

func JsonToStructApp(resp string) []model.SitesApplications {

	data := []byte(resp)

	var m []model.SitesApplications

	err := json.Unmarshal(data, &m)

	if err != nil {

		//Generating dummy json data

		dummyJson := `
				{
					"Website":  "null",
					"VirtualDirectory":  "null",
					"PhysicalPath":  "null",
					"ApplicationPool":  "null",
					"ApplicationName": "null"
				}
				`
		dummyData := []byte("[" + resp + "," + dummyJson + "]")

		err := json.Unmarshal(dummyData, &m)

		if err != nil {

			log.Fatal(err)
		}

		for i, s := range m {

			if s.ApplicationName == "null" || s.PhysicalPath == "null" {

				m = append(m[:i], m[i+1:]...)
			}

		}

	}

	for i, s := range m {

		if s.PSComputerName == "" {

			hostname, _ := os.Hostname()
			m[i].PSComputerName = hostname
		}
	}

	return m

}

func JsonToVdirInfo(resp string) model.VdirInfo {

	data := []byte(resp)

	var m model.VdirInfo

	err := json.Unmarshal(data, &m)

	if err != nil {

		log.Fatal(color.HiRedString("Multiple websites were found with the same virtual directory.\nPlease specify the website and application name of the virtual directory."))
	}

	return m

}
