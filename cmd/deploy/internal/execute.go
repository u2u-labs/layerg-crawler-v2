package internal

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	subgraphRepoUrl = ""
	branch          = "master"
)

func init() {
}

// Configuration for deployment
type DeploymentConfig struct {
	RepoURL          string
	Branch           string
	BasePath         string
	ShortID          string
	DatabaseName     string
	DatabasePassword string
	RedisDBNumber    int
	RedisPassword    string
	CRDBPort         int
	RedisPort        int
	QueryPort        int
}

// Template for the modified docker-compose.yml
const dockerComposeTemplate = `version: "3.5"

networks:
  default:
    external:
      name: crawler-network

services:
  app:
    image: u2labs/layerg-crawler:latest
    container_name: {{.ShortID}}-crawler-app
    command: --config layerg-crawler.yaml
    volumes:
      - ./cfg_{{.ShortID}}/layerg-crawler.yaml:/go/bin/layerg-crawler.yaml
      - ./cfg_{{.ShortID}}/subgraph.yaml:/go/bin/subgraph.yaml
    environment:
      - COCKROACH_DB_DRIVER=postgres
      - COCKROACH_DB_URL=postgres://root@crdb:26257/{{.DatabaseName}}?sslmode=disable
      - REDIS_DB_URL=redis:6379
      - REDIS_DB={{.RedisDBNumber}}
      - REDIS_DB_PASSWORD={{.RedisPassword}}
    depends_on:
      migrate:
        condition: service_completed_successfully
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: 300m
        tag: "{{ "{{" }}.ImageName{{ "}}" }}|{{ "{{" }}.Name{{ "}}" }}|{{ "{{" }}.ImageFullID{{ "}}" }}|{{ "{{" }}.FullID{{ "}}" }}"

  query:
    image: u2labs/layerg-crawler:latest
    container_name: {{.ShortID}}-crawler-query
    command: query --config layerg-crawler.yaml
    volumes:
      - ./cfg_{{.ShortID}}/layerg-crawler.yaml:/go/bin/layerg-crawler.yaml
      - ./cfg_{{.ShortID}}/schema.graphql:/go/bin/schema.graphql
    environment:
      - COCKROACH_DB_DRIVER=postgres
      - COCKROACH_DB_URL=postgres://root@crdb:26257/{{.DatabaseName}}?sslmode=disable
      - REDIS_DB_URL=redis:6379
      - REDIS_DB={{.RedisDBNumber}}
      - REDIS_DB_PASSWORD={{.RedisPassword}}
    depends_on:
      migrate:
        condition: service_completed_successfully
    ports:
      - "{{.QueryPort}}:8084"
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: 300m
        tag: "{{ "{{" }}.ImageName{{ "}}" }}|{{ "{{" }}.Name{{ "}}" }}|{{ "{{" }}.ImageFullID{{ "}}" }}|{{ "{{" }}.FullID{{ "}}" }}"

  db-setup:
    image: cockroachdb/cockroach:v24.2.1
    container_name: {{.ShortID}}-crawler-dbsetup
    command: sql --insecure --host=crdb --port=26257 --execute='CREATE DATABASE IF NOT EXISTS {{.DatabaseName}};'

  migrate:
    build:
      dockerfile: migrate.Dockerfile
    container_name: {{.ShortID}}-crawler-migrate
    command: ["system-migrate-up", "generated-migrate-up"]
    environment:
      - GOOSE_DRIVER=postgres
      - GOOSE_DBSTRING=postgres://root@crdb:26257/{{.DatabaseName}}?sslmode=disable
    depends_on:
      - db-setup
`

// Safe characters for folder names: a-z, A-Z, 0-9
const safeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortID(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	// Generate random bytes
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Map each byte to a character from safeChars
	result := make([]byte, length)
	for i, b := range bytes {
		result[i] = safeChars[int(b)%len(safeChars)]
	}

	return string(result), nil
}

// expandPath handles environment variables and tilde expansion in paths
func expandPath(path string) (string, error) {
	// Handle $HOME or ~ at the beginning of the path
	if strings.HasPrefix(path, "$HOME") || strings.HasPrefix(path, "~") {
		currentUser, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("failed to get current user: %w", err)
		}

		if strings.HasPrefix(path, "$HOME") {
			path = strings.Replace(path, "$HOME", currentUser.HomeDir, 1)
		} else if strings.HasPrefix(path, "~") {
			path = strings.Replace(path, "~", currentUser.HomeDir, 1)
		}
	}

	// Handle other environment variables
	for strings.Contains(path, "$") {
		startIndex := strings.Index(path, "$")
		endIndex := strings.Index(path[startIndex:], "/")

		if endIndex == -1 {
			endIndex = len(path)
		} else {
			endIndex += startIndex
		}

		envVar := path[startIndex:endIndex]
		if strings.Contains(envVar, "/") {
			envVar = envVar[:strings.Index(envVar, "/")]
		}

		// Remove $ from the environment variable name
		envName := envVar[1:]
		envValue := os.Getenv(envName)

		if envValue == "" {
			logger.Warnf("Environment variable %s not found, keeping as is", envName)
			break
		}

		path = strings.Replace(path, envVar, envValue, 1)
	}

	return path, nil
}

// createDatabaseInCockroach creates a new database in CockroachDB
func createDatabaseInCockroach(config DeploymentConfig) error {
	cmd := exec.Command("docker", "exec", "crawler-db",
		"cockroach", "sql", "--insecure",
		"--execute", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.DatabaseName))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create database: %s\nOutput: %s", err, output)
	}

	logger.Infof("Created database %s in CockroachDB", config.DatabaseName)
	return nil
}

// updateConfigFiles updates configuration files with the new database settings
func updateConfigFiles(config DeploymentConfig, deployDir string) error {
	// Update layerg-crawler.yaml (assuming it needs to be updated with DB connection info)
	configFilePath := filepath.Join(deployDir, "layerg-crawler.yaml")

	if _, err := os.Stat(configFilePath); err == nil {
		content, err := os.ReadFile(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		updatedContent := strings.Replace(
			string(content),
			"database: layerg",
			fmt.Sprintf("database: %s", config.DatabaseName),
			-1,
		)

		updatedContent = strings.Replace(
			updatedContent,
			"db: 0",
			fmt.Sprintf("db: %d", config.RedisDBNumber),
			-1,
		)

		err = os.WriteFile(configFilePath, []byte(updatedContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write updated config file: %w", err)
		}
	}

	return nil
}

// createDockerComposeFile creates a docker-compose.yml file from the template
func createDockerComposeFile(config DeploymentConfig, deployDir string) error {
	tmpl, err := template.New("docker-compose").Parse(dockerComposeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	filePath := filepath.Join(deployDir, "docker-compose.yaml")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create docker-compose file: %w", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, config)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// deployGraph clones the repo and starts the services
func deployGraph(config DeploymentConfig) error {
	// Create deployment directory with short ID
	deployDir := filepath.Join(config.BasePath, fmt.Sprintf("layerg-crawler-%s", config.ShortID))

	logger.Infof("Creating deployment directory: %s", deployDir)
	if err := os.MkdirAll(deployDir, 0755); err != nil {
		return fmt.Errorf("failed to create deployment directory: %w", err)
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", config.Branch, config.RepoURL, deployDir)
	logger.Infof("Cloning repository from %s to %s, branch %s: %s", config.RepoURL, deployDir, config.Branch, cmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone repository: %s\nOutput: %s", err, output)
	}

	// Create docker-compose.yml file
	if err := createDockerComposeFile(config, deployDir); err != nil {
		return err
	}

	// Create database in CockroachDB
	if err := createDatabaseInCockroach(config); err != nil {
		return err
	}

	// Update configuration files
	if err := updateConfigFiles(config, deployDir); err != nil {
		return err
	}

	// Start services using docker-compose
	logger.Infof("Starting services with docker-compose in %s", deployDir)
	cmd = exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = deployDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start services: %s\n", err)
	}

	logger.Infof("Deployment successful! ShortID: %s, Database: %s, Redis DB: %d, Query Port: %d",
		config.ShortID, config.DatabaseName, config.RedisDBNumber, config.QueryPort)

	return nil
}

func executeFn(cmd *cobra.Command, args []string) {
	shortID, err := generateShortID(6)
	if err != nil {
		logger.Fatalf("Failed to generate short ID: %v", err)
	}

	// Expand the base path to handle environment variables
	basePath, err := expandPath("$HOME/.layerg/crawler/deployments")
	if err != nil {
		logger.Fatalf("Failed to expand base path: %v", err)
	}

	// Configure deployment
	config := DeploymentConfig{
		RepoURL:          subgraphRepoUrl,
		Branch:           branch,
		BasePath:         basePath,
		ShortID:          shortID,
		DatabaseName:     fmt.Sprintf("layerg_%s", shortID),
		DatabasePassword: os.Getenv("COCKROACH_PASSWORD"),
		RedisDBNumber:    1, // Increment this for each new deployment
		CRDBPort:         26257,
		RedisPort:        6379,
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		QueryPort:        8084 + 1, // Use a different port for each deployment
	}

	// Deploy the graph node
	if err := deployGraph(config); err != nil {
		logger.Fatalf("Deployment failed: %v", err)
	}

	logger.Info("Subgraph deployment completed successfully")
}
