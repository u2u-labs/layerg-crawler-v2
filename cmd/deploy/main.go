package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	ipfsServerURL string
	logger        *zap.SugaredLogger
)

func main() {
	// Initialize logger
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	defer l.Sync()
	logger = l.Sugar()

	_ = godotenv.Load()
	flag.StringVar(&ipfsServerURL, "ipfs-url", "https://api.pinata.cloud/pinning/pinFileToIPFS", "IPFS server URL to upload files")
	flag.Parse()

	err = moveFilesToBuild()
	if err != nil {
		logger.Fatalf("Failed to move files: %v", err)
	}

	err = deploy()
	if err != nil {
		logger.Fatalf("Failed to deploy: %v", err)
	}
}

func deploy() error {
	apiKey := os.Getenv("PINATA_API_KEY")
	apiSecret := os.Getenv("PINATA_API_SECRET")

	// upload build dir
	buildCids, err := uploadFolderToPinata("build", apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("failed to upload build folder: %v", err)
	}
	logger.Infof("Build folder uploaded to IPFS: %v", buildCids)

	// write cids to temp file
	tempFilePath, err := writeCIDsToTempFile(buildCids)
	if err != nil {
		logger.Errorf("Error writing CIDs: %v", err)
		return err
	}
	logger.Infof("CIDs written to: %s", tempFilePath)

	return nil
}

func uploadFolderToPinata(folderPath, apiKey, apiSecret string) (map[string]string, error) {
	// Collect files in the folder
	files := make(map[string]string)
	var wg sync.WaitGroup
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			// Upload each file separately
			cid, err := uploadFileToPinata(path, apiKey, apiSecret)
			if err != nil {
				logger.Errorf("Failed to upload %s: %v", path, err)
				return
			}

			files[path] = cid
			logger.Infof("Uploaded %s -> %s", path, cid)
		}(path)
		return nil
	})

	wg.Wait()
	if err != nil {
		return nil, fmt.Errorf("error walking through folder: %v", err)
	}

	return files, nil
}

func uploadFileToPinata(filePath, apiKey, apiSecret string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a buffer and multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Attach file
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	// Close writer
	writer.Close()

	// Create HTTP request
	req, err := http.NewRequest("POST", ipfsServerURL, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("pinata_api_key", apiKey)
	req.Header.Set("pinata_secret_api_key", apiSecret)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	logger.Infof("Uploading file... %s", filePath)
	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON response
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract CID
	if ipfsHash, ok := result["IpfsHash"].(string); ok {
		return ipfsHash, nil
	}
	return "", fmt.Errorf("unexpected response format: %s", respBody)
}

func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func moveFilesToBuild() error {
	// Create the build directory if it doesn't exist
	buildDir := "build"
	// remove old build directory
	err := os.RemoveAll(buildDir)
	if err != nil {
		logger.Errorf("Failed to remove old build directory: %v", err)
		return err
	}

	err = os.MkdirAll(buildDir+"/migrations", 0755)
	if err != nil {
		logger.Errorf("Failed to create build directory: %v", err)
		return err
	}

	// List of files to copy
	files := []string{"layerg-crawler", ".layerg-crawler.yaml", "subgraph.yaml"}

	// Copy individual files
	for _, file := range files {
		dstPath := filepath.Join(buildDir, filepath.Base(file))
		err = copyFile(file, dstPath)
		if err != nil {
			logger.Errorf("Failed to copy %s: %v", file, err)
			return err
		}
		logger.Infof("Copied %s to %s", file, dstPath)
	}

	// Copy SQL files from generated/migrations/
	migrationSrcDir := "generated/migrations"
	migrationFiles, err := filepath.Glob(filepath.Join(migrationSrcDir, "*.sql"))
	if err != nil {
		logger.Errorf("Failed to list migration files: %v", err)
		return err
	}

	for _, srcPath := range migrationFiles {
		dstPath := filepath.Join(buildDir, "migrations", filepath.Base(srcPath))
		err = copyFile(srcPath, dstPath)
		if err != nil {
			logger.Errorf("Failed to copy %s: %v", srcPath, err)
			return err
		}
		logger.Infof("Copied %s to %s", srcPath, dstPath)
	}
	return nil
}

func writeCIDsToTempFile(data map[string]string) (string, error) {
	// Get the system's temp directory
	tempDir := os.TempDir()

	// Define the full file path
	filePath := filepath.Join(tempDir, "subgraph_deploy_cids.json")

	// Open the file for writing (create if not exists, truncate if exists)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer file.Close()

	// Convert map to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Write JSON to file
	_, err = file.Write(jsonData)
	if err != nil {
		return "", fmt.Errorf("failed to write JSON to file: %v", err)
	}

	return filePath, nil
}
