package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// deploy cmd to deploy binary and config file to ipfs server

var (
	ipfsServerURL, binaryPath, configPath, generatedMigrationPath string
)

func main() {
	_ = godotenv.Load()
	flag.StringVar(&ipfsServerURL, "ipfs-url", "https://api.pinata.cloud/pinning/pinFileToIPFS", "IPFS server URL to upload files")
	flag.StringVar(&binaryPath, "e", "layerg-crawler", "Path to the binary file to deploy")
	flag.StringVar(&configPath, "c", ".layerg-crawler.yaml", "Path to the config file to deploy")
	flag.StringVar(&generatedMigrationPath, "m", "generated/migrations", "Path to the generated migration files")
	flag.Parse()

	err := deploy()
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
}

func deploy() error {
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	// Upload binary file
	binaryCID, err := uploadToPinata(binaryPath, apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("failed to upload binary: %v", err)
	}
	fmt.Println("Binary uploaded to IPFS:", binaryCID)

	//// Upload config file
	//configCID, err := uploadToPinata(configPath, apiKey, apiSecret)
	//if err != nil {
	//	return fmt.Errorf("failed to upload config file: %v", err)
	//}
	//fmt.Println("Config file uploaded to IPFS:", configCID)

	// Upload migration files
	uploadMigrationFiles(generatedMigrationPath, apiKey, apiSecret)

	return nil
}

func uploadToPinata(filePath, apiKey, apiSecret string) (string, error) {
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

	fmt.Println("Uploading file...", filePath)
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

func uploadMigrationFiles(migrationPath, apiKey, apiSecret string) {
	// List files in the migration directory
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration directory: %v", err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		wg.Add(1)
		filePath := fmt.Sprintf("%s/%s", migrationPath, file.Name())
		go func(filePath string) {
			defer wg.Done()
			cid, err := uploadToPinata(filePath, apiKey, apiSecret)
			if err != nil {
				log.Fatalf("Failed to upload migration file: %v", err)
			}
			fmt.Printf("Migration file %s uploaded to IPFS with CID: %s\n", file.Name(), cid)
		}(filePath)
	}
	wg.Wait()
}
