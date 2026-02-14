package servers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/internal/tests"
	"github.com/goccy/go-yaml"
)

type dockerCompose struct {
	Services map[string]struct {
		Image string `yaml:"image"`
	} `yaml:"services"`
}

func loadDockerComposeServiceImage(t *testing.T, service string) (string, string) {
	t.Helper()

	// Get the absolute path to the docker-compose.yml file
	dockerComposePath := filepath.Join(tests.Root, "docker-compose.dev.yaml")

	// Check if the file exists
	if _, err := os.Stat(dockerComposePath); errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("docker-compose.yml file not found at path: %s", dockerComposePath))
	}

	// Open the docker-compose file
	file, err := os.Open(dockerComposePath)
	if err != nil {
		panic(fmt.Errorf("failed to open docker-compose.dev.yaml file: %w", err))
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("failed to read docker-compose.dev.yaml file: %w", err))
	}

	// Unmarshal the YAML content into the struct
	var compose dockerCompose
	if err := yaml.Unmarshal(content, &compose); err != nil {
		panic(fmt.Errorf("failed to parse docker-compose.dev.yaml file: %w", err))
	}

	serviceData, ok := compose.Services[service]
	if !ok {
		panic(fmt.Errorf("service %s not found in docker-compose.dev.yaml file", service))
	}

	image, tag, ok := strings.Cut(serviceData.Image, ":")
	if !ok {
		panic(fmt.Errorf("invalid image format for service %s: %s", service, serviceData.Image))
	}

	return image, tag
}
