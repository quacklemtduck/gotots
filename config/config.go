package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OutPath     string `yaml:"out-path"`
	UseComments bool   `yaml:"use-comments"`
	ConfigPath  string `yaml:"-"`
}

var defaultConfig = Config{
	OutPath:     "",
	UseComments: false,
}

func LoadConfig(startPath string) Config {
	path, ok := findGototsFile(startPath)
	if !ok {
		return defaultConfig
	}
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("Error reading config file:", err)
	}
	var conf Config = Config{}
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalln("Error parsing config file:", err)
	}
	conf.ConfigPath = path
	return conf
}

// findGototsFile searches for "gotots.yaml" in the given path and its parent directories.
// It returns the path to "gotots.yaml" if found, along with an `ok` boolean indicating if it was found.
func findGototsFile(startPath string) (string, bool) {
	currentPath := startPath

	for {
		// Check if "gotots.yaml" exists in the current directory.
		gototsPath := filepath.Join(currentPath, "gotots.yaml")
		if _, err := os.Stat(gototsPath); err == nil {
			// Return the path to "gotots.yaml" if found.
			return gototsPath, true
		}

		// Check if "gotots.yml" exists in the current directory.
		gototsPath = filepath.Join(currentPath, "gotots.yml")
		if _, err := os.Stat(gototsPath); err == nil {
			// Return the path to "gotots.yaml" if found.
			return gototsPath, true
		}

		// Check if "go.mod" exists in the current directory.
		goModPath := filepath.Join(currentPath, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// Stop the search if we encounter "go.mod".
			return "", false
		}

		// Get the parent directory.
		parentPath := filepath.Dir(currentPath)

		// If we've reached the root directory, stop the search.
		if parentPath == currentPath {
			return "", false
		}

		// Move up to the parent directory.
		currentPath = parentPath
	}
}
