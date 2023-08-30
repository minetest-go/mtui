package depanalyzer

type DependsInfo struct {
	Depends         []string `json:"depends"`
	OptionalDepends []string `json:"optional_depends"`
}

type ModConfig struct {
	*DependsInfo
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Title       string `json:"title"`
}
