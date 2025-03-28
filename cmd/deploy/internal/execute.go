package internal

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Gateway URL
var (
	subgraphCid, configCid, executableCid, migrationCid string

	baseGateway string
)

// Deployment Metadata struct to track deployment state
type DeploymentMetadata struct {
	ID             string    `json:"id"`
	SubgraphCid    string    `json:"subgraph_cid"`
	ConfigCid      string    `json:"config_cid"`
	ExecutableCid  string    `json:"executable_cid"`
	MigrationCid   string    `json:"migration_cid"`
	SubgraphPath   string    `json:"subgraph_path"`
	ConfigPath     string    `json:"config_path"`
	ExecutablePath string    `json:"executable_path"`
	MigrationPath  string    `json:"migration_path"`
	DatabaseName   string    `json:"database_name"`
	DatabaseConn   string    `json:"database_conn"`
	LocalPath      string    `json:"local_path"`
	CreatedAt      time.Time `json:"created_at"`
	LastRunAt      time.Time `json:"last_run_at"`
	Status         string    `json:"status"` // pending, running, completed, failed
}

type DeploymentManager struct {
	metadataFile string
	metadata     DeploymentMetadata
	logger       *zap.SugaredLogger
}

func NewDeploymentManager(subgraphCid, configCid, executableCid, migrationCid string, cids map[string]string, logger *zap.SugaredLogger) *DeploymentManager {
	// Use a consistent metadata file path
	metadataPath := filepath.Join(os.TempDir(), fmt.Sprintf("deployment-manager_%s.json", executableCid))
	migrationPath := ""
	for key := range cids {
		if strings.HasPrefix(key, "build/migrations/") && strings.HasSuffix(key, ".sql") {
			migrationPath = strings.ReplaceAll(key, "build", "")
			break
		}
	}

	return &DeploymentManager{
		metadataFile: metadataPath,
		logger:       logger,
		metadata: DeploymentMetadata{
			ID:             uuid.New().String(),
			SubgraphCid:    subgraphCid,
			ConfigCid:      configCid,
			ExecutableCid:  executableCid,
			MigrationCid:   migrationCid,
			SubgraphPath:   "/subgraph.yaml",
			ConfigPath:     "/.layerg-crawler.yaml",
			ExecutablePath: "/layerg-crawler",
			MigrationPath:  migrationPath,
			CreatedAt:      time.Now(),
			Status:         "pending",
		},
	}
}

func (dm *DeploymentManager) LoadOrCreateMetadata() error {
	dm.logger.Info("Loading or creating deployment metadata",
		zap.String("metadataFile", dm.metadataFile),
		zap.String("ecid", dm.metadata.ExecutableCid))

	// Try to load existing metadata
	data, err := os.ReadFile(dm.metadataFile)
	if err == nil {
		// Metadata exists, try to parse
		var existingMetadata DeploymentMetadata
		if err := json.Unmarshal(data, &existingMetadata); err == nil {
			// Check if existing deployment is recent and valid
			if existingMetadata.ExecutableCid == dm.metadata.ExecutableCid &&
				time.Since(existingMetadata.CreatedAt) < 24*time.Hour {
				dm.logger.Info("Found existing valid metadata",
					zap.String("id", existingMetadata.ID),
					zap.Time("createdAt", existingMetadata.CreatedAt))
				dm.metadata = existingMetadata
				return nil
			}
		}
	}

	// If no valid existing metadata, create new
	dm.logger.Info("Creating new deployment metadata")
	return dm.SaveMetadata()
}

func (dm *DeploymentManager) SaveMetadata() error {
	data, err := json.Marshal(dm.metadata)
	if err != nil {
		dm.logger.Error("Failed to marshal metadata",
			zap.Error(err))
		return err
	}

	err = os.WriteFile(dm.metadataFile, data, 0644)
	if err != nil {
		dm.logger.Error("Failed to save metadata file",
			zap.String("path", dm.metadataFile),
			zap.Error(err))
		return err
	}

	return nil
}

func (dm *DeploymentManager) FetchFromIPFS() (string, error) {
	dm.logger.Info("Fetching subgraph from IPFS",
		zap.String("ecid", dm.metadata.ExecutableCid))

	// If local path exists and is valid, return it
	if dm.metadata.LocalPath != "" && !noCache {
		if _, err := os.Stat(dm.metadata.LocalPath); err == nil {
			dm.logger.Info("Using existing local path",
				zap.String("path", dm.metadata.LocalPath))
			return dm.metadata.LocalPath, nil
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Use cautiously
			},
		},
	}

	// Create a temporary directory for download
	tempDir, err := os.MkdirTemp("", "subgraph-*")
	if err != nil {
		dm.logger.Error("Failed to create temp directory", zap.Error(err))
		return "", err
	}

	// Path to save binary and config
	subgraphPath := filepath.Join(tempDir, "subgraph.yaml")
	binaryPath := filepath.Join(tempDir, "crawler")
	configPath := filepath.Join(tempDir, "config.yaml")
	migrationDir := filepath.Join(tempDir, "migrations")
	migrationPath := filepath.Join(tempDir, dm.metadata.MigrationPath)

	// Ensure migration directory exists
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		dm.logger.Error("Failed to create migrations directory",
			zap.String("path", migrationDir),
			zap.Error(err))
		return "", err
	}

	// Prepare downloads
	var wg sync.WaitGroup
	var downloadErrs []error
	var mu sync.Mutex

	downloads := []struct {
		cid  string
		path string
		name string
	}{
		{dm.metadata.ExecutableCid + dm.metadata.ExecutablePath, binaryPath, "binary"},
		{dm.metadata.ConfigCid + dm.metadata.ConfigPath, configPath, "config"},
		{dm.metadata.MigrationCid + dm.metadata.MigrationPath, migrationPath, "migrations"},
		{dm.metadata.SubgraphCid + dm.metadata.SubgraphPath, subgraphPath, "subgraph"},
	}

	// Concurrent downloads
	for _, download := range downloads {
		wg.Add(1)
		go func(d struct {
			cid  string
			path string
			name string
		}) {
			defer wg.Done()

			err := downloadFileFromGateway(client, baseGateway, d.cid, d.path, dm.logger)
			if err != nil {
				mu.Lock()
				downloadErrs = append(downloadErrs,
					fmt.Errorf("failed to fetch %s: %w", d.name, err))
				mu.Unlock()
				dm.logger.Error("Download failed",
					zap.String("component", d.name),
					zap.Error(err))
			}
		}(download)
	}

	// Wait for all downloads to complete
	wg.Wait()

	// Check if any download failed
	if len(downloadErrs) > 0 {
		return "", fmt.Errorf("multiple download errors: %v", downloadErrs)
	}

	// Make binary executable
	err = os.Chmod(binaryPath, 0755)
	if err != nil {
		dm.logger.Error("Failed to make binary executable",
			zap.String("path", binaryPath),
			zap.Error(err))
		return "", err
	}

	// Update metadata
	dm.metadata.LocalPath = tempDir
	dm.SaveMetadata()

	dm.logger.Info("Successfully fetched subgraph from IPFS",
		zap.String("localPath", tempDir))

	return tempDir, nil
}

func downloadFileFromGateway(client *http.Client, baseGateway, cid, localPath string, logger *zap.SugaredLogger) error {
	// Construct full URL
	fileUrl := fmt.Sprintf("%s/%s", baseGateway, cid)
	fmt.Println(fileUrl)

	// Create request
	req, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		logger.Error("Failed to create request",
			zap.String("fileUrl", fileUrl),
			zap.Error(err))
		return err
	}

	// Add user agent to improve chances of successful download
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to download file",
			zap.String("fileUrl", fileUrl),
			zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create local file
	out, err := os.Create(localPath)
	if err != nil {
		logger.Error("Failed to create local file",
			zap.String("path", localPath),
			zap.Error(err))
		return err
	}
	defer out.Close()

	// Copy content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logger.Error("Failed to write file",
			zap.String("path", localPath),
			zap.Error(err))
		return err
	}

	return nil
}

func (dm *DeploymentManager) CreateIsolatedDatabase(baseConnString string) (*sql.DB, string, error) {
	dm.logger.Info("Creating isolated database")

	// If database already exists, return existing connection
	if dm.metadata.DatabaseName != "" && dm.metadata.DatabaseConn != "" {
		db, err := sql.Open("postgres", dm.metadata.DatabaseConn)
		if err == nil {
			dm.logger.Info("Reusing existing database",
				zap.String("dbName", dm.metadata.DatabaseName))
			return db, dm.metadata.DatabaseName, nil
		}
	}

	// Generate unique database name
	uniqueDBName := fmt.Sprintf("subgraph_%s", strings.ReplaceAll(dm.metadata.ID, "-", "_"))

	// Connect to base PostgreSQL instance
	baseDB, err := sql.Open("postgres", baseConnString)
	if err != nil {
		dm.logger.Error("Failed to connect to base database",
			zap.String("connString", baseConnString),
			zap.Error(err))
		return nil, "", err
	}
	defer baseDB.Close()

	// Create new database
	_, err = baseDB.Exec(fmt.Sprintf("CREATE DATABASE %s", uniqueDBName))
	if err != nil {
		dm.logger.Error("Failed to create database",
			zap.String("dbName", uniqueDBName),
			zap.Error(err))
		return nil, "", err
	}

	// Parse base connection string
	u, err := url.Parse(baseConnString)
	if err != nil {
		dm.logger.Error("Failed to parse connection string",
			zap.String("connString", baseConnString),
			zap.Error(err))
		return nil, "", err
	}

	// Replace the database name in the path
	u.Path = "/" + uniqueDBName

	// Construct new connection string
	isolatedConnString := u.String()

	// Connect to the new database
	isolatedDB, err := sql.Open("postgres", isolatedConnString)
	if err != nil {
		dm.logger.Error("Failed to connect to isolated database",
			zap.String("connString", isolatedConnString),
			zap.Error(err))
		return nil, "", err
	}

	// Update metadata
	dm.metadata.DatabaseName = uniqueDBName
	dm.metadata.DatabaseConn = isolatedConnString
	dm.SaveMetadata()

	dm.logger.Info("Successfully created isolated database",
		zap.String("dbName", uniqueDBName))

	return isolatedDB, uniqueDBName, nil
}

func (dm *DeploymentManager) RunMigrations(db *sql.DB, migrationDir string) error {
	dm.logger.Info("Running migrations",
		zap.String("migrationDir", migrationDir))

	db, err := sql.Open("postgres", dm.metadata.DatabaseConn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Set the migration directory
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %v", err)
	}

	// Run migrations
	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}

	dm.logger.Info("Migrations completed successfully")

	return nil
}

func readCidsFromFile(subgraphCid, configCid, executableCid, migrationCid *string) (map[string]string, error) {
	// Get the system's temp directory
	tempDir := os.TempDir()

	// Define the full file path
	filePath := filepath.Join(tempDir, "subgraph_deploy_cids.json")

	// Open the file
	file, err := os.Open(filePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// If the file doesn't exist, return empty strings and no error
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode JSON into a map
	var cids map[string]string
	err = json.NewDecoder(file).Decode(&cids)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	*subgraphCid = cids["build/subgraph.yaml"]
	*configCid = cids["build/.layerg-crawler.yaml"]
	*executableCid = cids["build/layerg-crawler"]

	for key, value := range cids {
		if strings.HasPrefix(key, "build/migrations/") && strings.HasSuffix(key, ".sql") {
			*migrationCid = value
			break // Take the first match
		}
	}

	if *migrationCid == "" {
		return nil, fmt.Errorf("no migration file found in JSON")
	}

	return cids, nil
}

func executeFn(cmd *cobra.Command, args []string) {
	// Read CIDs from file (if available)
	var err error
	cids, err := readCidsFromFile(&subgraphCid, &configCid, &executableCid, &migrationCid)
	if err != nil {
		logger.Fatal("Failed to read subgraph")
		return
	}

	// Get flag values from Cobra (defaulting to file values if they exist)
	subgraphCid, _ = cmd.Flags().GetString("scid")
	configCid, _ = cmd.Flags().GetString("ccid")
	executableCid, _ = cmd.Flags().GetString("ecid")
	migrationCid, _ = cmd.Flags().GetString("mcid")
	baseGateway, _ = cmd.Flags().GetString("gw")

	// Load environment variables
	err = godotenv.Load()
	if err != nil {
		logger.Warn("Error loading .env file", zap.Error(err))
	}

	// Create deployment manager
	deploymentManager := NewDeploymentManager(subgraphCid, configCid, executableCid, migrationCid, cids, logger)

	// Log start of deployment
	logger.Info("Starting subgraph deployment",
		zap.String("ecid", executableCid))

	// Load or create metadata
	err = deploymentManager.LoadOrCreateMetadata()
	if err != nil {
		logger.Fatal("Failed to load or create metadata", zap.Error(err))
	}

	// Fetch from IPFS
	tempDir, err := deploymentManager.FetchFromIPFS()
	if err != nil {
		logger.Fatal("Failed to fetch from IPFS", zap.Error(err))
	}
	defer os.RemoveAll(tempDir)

	// Base connection string
	baseConnString := os.Getenv("COCKROACH_DB_URL")

	// Create isolated database
	db, _, err := deploymentManager.CreateIsolatedDatabase(baseConnString)
	if err != nil {
		logger.Fatal("Failed to create isolated database", zap.Error(err))
	}
	defer db.Close()

	// Run migrations
	migrationDir := filepath.Join(tempDir, "migrations")
	err = deploymentManager.RunMigrations(db, migrationDir)
	if err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}

	// Execute crawler with isolated config
	configPath := filepath.Join(tempDir, "config.yaml")
	binaryPath := filepath.Join(tempDir, "crawler")

	binaryCmd := exec.Command(binaryPath, "--config", configPath)
	binaryCmd.Stdout = os.Stdout
	binaryCmd.Stderr = os.Stderr
	binaryCmd.Env = append(os.Environ(),
		fmt.Sprintf("COCKROACH_DB_URL=%s", deploymentManager.metadata.DatabaseConn),
	)

	err = binaryCmd.Run()
	if err != nil {
		logger.Fatal("Crawler execution failed", zap.Error(err))
	}

	// Update metadata status
	deploymentManager.metadata.Status = "completed"
	deploymentManager.SaveMetadata()

	logger.Info("Subgraph deployment completed successfully")
}
