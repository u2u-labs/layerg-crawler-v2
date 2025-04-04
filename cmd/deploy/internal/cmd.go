package internal

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	rootCmd = &cobra.Command{
		Use:   "layerg-deploy",
		Short: "CLI tool for deploying and executing files on IPFS",
		Long:  "LayerG-Deploy is a command-line tool to publish builds to an IPFS server and execute downloaded files.",
	}
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a build to IPFS",
		Long:  "Uploads a specified directory or file to an IPFS server and returns the content identifier (CID).",
		Run:   deployFn,
	}
	executeCmd = &cobra.Command{
		Use:   "run",
		Short: "Download and execute files from IPFS",
		Long:  "Fetches files from an IPFS server using the provided CID and runs them.",
		Run:   executeFn,
	}
	logger          *zap.SugaredLogger
	noCache         bool
	devExperimental bool
)

func init() {
	// Initialize logger
	l, err := zap.NewDevelopment(zap.AddStacktrace(zapcore.InvalidLevel))
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	defer l.Sync()
	logger = l.Sugar()

	_ = godotenv.Load()
	deployCmd.Flags().String("ipfs-url", "https://api.pinata.cloud/pinning/pinFileToIPFS", "IPFS server URL to upload files")
	executeCmd.Flags().StringVar(&subgraphCid, "scid", "", "Subgraph CID")
	executeCmd.Flags().StringVar(&configCid, "ccid", "", "Config CID")
	executeCmd.Flags().StringVar(&executableCid, "ecid", "", "Executable CID")
	executeCmd.Flags().StringVar(&migrationCid, "mcid", "", "Migration CID")
	executeCmd.Flags().StringVar(&baseGateway, "gw", "https://gateway.pinata.cloud/ipfs", "IPFS gateway URL")

	deployCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable caching during deployment")
	executeCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable caching during execution")
	deployCmd.Flags().BoolVar(&devExperimental, "dev", false, "Enable experimental features for development")
	executeCmd.Flags().BoolVar(&devExperimental, "dev", false, "Enable experimental features for development")

	rootCmd.AddCommand(deployCmd, executeCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
