package model

type AppPool struct {
	Name                  string `json:"name"`
	State                 string `json:"state"`
	Sites                 string `json:"sites"`
	StartMode             string `json:"startMode"`
	ManagedPipelineMode   string `json:"managedPipelineMode"`
	ManagedRuntimeVersion string `json:"managedRuntimeVersion"`
	AutoStart             bool   `json:"autoStart"`
	IdleTimeoutAction     string `json:"idleTimeoutAction"`
	Idletimeout           int64  `json:"idletimeout"`
	PSComputerName        string `json:"PSComputerName"`
}

type WindowsService struct {
	Name        string `json:"Name"`
	DisplayName string `json:"DisplayName"`
	Description string `json:"Description"`
	State       string `json:"state"`
	StartMode   string `json:"startmode"`
	//Path           string `json:"Path"`
	ExecutablePath string `json:"ExecutablePath"`
	PSComputerName string `json:"PSComputerName"`
}

type FolderList []struct {
	Name           string `json:"Name"`
	Path           string `json:"Path"`
	Size           string `json:"Size"`
	LastAccessTime string `json:"LastAccessTime"`
	LastWriteTime  string `json:"LastWriteTime"`
	ComputerName   string `json:"ComputerName"`
}

type VirtualDir struct {
	VirtualDirectory string `json:"VirtualDirectory"`
	SiteName         string `json:"SiteName"`
	ApplicationName  string `json:"ApplicationName"`
	PhysicalPath     string `json:"PhysicalPath"`
	PSComputerName   string `json:"PSComputerName"`
}

type SitesIIS struct {
	Name            string `json:"name"`
	State           string `json:"state"`
	ApplicationPool string `json:"applicationPool"`
	Bindings        string `json:"bindings"`
	PhysicalPath    string `json:"physicalPath"`
	Protocol        string `json:"protocol"`
	PreloadEnabled  string `json:"PreloadEnabled"`
	//Applications       string `json:"Applications"`
	//VirtualDirectories string `json:"VirtualDirectories"`
	PSComputerName string `json:"PSComputerName"`
}

type SitesApplications struct {
	ApplicationName string `json:"ApplicationName"`
	Website         string `json:"Website"`
	ApplicationPool string `json:"ApplicationPool"`
	PhysicalPath    string `json:"physicalPath"`
	PreloadEnabled  bool   `json:"PreloadEnabled"`
	PSComputerName  string `json:"PSComputerName"`
}

//VirtualDirectories string `json:"VirtualDirectories"`

type VdirInfo struct {
	SiteName        string `json:"SiteName"`
	ApplicationName string `json:"ApplicationName"`
}
