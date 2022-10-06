package modmanager

type ModStatus struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}
