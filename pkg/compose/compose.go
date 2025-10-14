// Definitions of blocks found in Docker Compose files
// and utility methods to parse them.
package compose

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// Docker Compose configuration, maps 1:1 with a docker-compose.yaml
type Config struct {
	Services map[string]Service `yaml:"services"`
}

// A Service is a container running one or more application
type Service struct {
	Image string `yaml:"image"`
}

// Parse a Docker Compose configuration file into a workable struct
func Parse(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error opening file %s: %s", path, err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error reading file %s: %s", path, err)
	}

	return config
}
