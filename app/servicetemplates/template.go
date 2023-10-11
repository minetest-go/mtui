package servicetemplates

import (
	"embed"
	"fmt"
	"os"
	"path"
	"text/template"
)

//go:embed *
var Files embed.FS

func CreateTemplate(filename string, destination_path string, model any) error {
	destination_file := path.Join(destination_path, filename)
	fi, _ := os.Stat(destination_file)
	if fi != nil {
		// file already exists, do nothing
		return nil
	}

	data, err := Files.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not open file '%s': %v", filename, err)
	}

	t, err := template.New("tmpl").Parse(string(data))
	if err != nil {
		return fmt.Errorf("error templating %s: %v", filename, err)
	}

	f, err := os.OpenFile(destination_file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open file '%s': %v", destination_file, err)
	}
	defer f.Close()

	return t.Execute(f, model)
}

func CreateMatterbridgeTemplate(destination_path string, model *MatterbridgeModel) error {
	return CreateTemplate("matterbridge.toml", destination_path, model)
}

func CreateMapserverTemplate(destination_path string, model *MapserverModel) error {
	return CreateTemplate("mapserver.json", destination_path, model)
}
