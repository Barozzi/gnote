package config

import (
	"fmt"
	"os"
	"os/user"

	"gopkg.in/yaml.v3"
)

type Config struct {
	VaultPath    string `yaml:"vault_path"`
	DayPath      string `yaml:"day_subpath"`
	ProjectsPath string `yaml:"projects_subpath"`
	AreasPath    string `yaml:"areas_subpath"`
	ArchivesPath string `yaml:"archives_subpath"`
}

func ReadConfig() (*Config, error) {
	usr, err := user.Current()
	file, err := os.Open(fmt.Sprintf("%s/.config/gnote/gnote.yaml", usr.HomeDir))
	if err != nil {
		fmt.Println("ReadConfig: failed to open file.", usr.HomeDir)
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("ReadConfig: failed to parse yaml")
		return nil, err
	}

	return &config, nil
}
