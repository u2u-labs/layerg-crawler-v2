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
	publishCmd = &cobra.Command{
		Use:   "publish",
		Short: "Publish a build to IPFS",
		Long:  "Uploads a specified directory or file to an IPFS server and returns the content identifier (CID).",
		Run:   publishFn,
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
	publishCmd.Flags().String("ipfs-url", "https://api.pinata.cloud/pinning/pinFileToIPFS", "IPFS server URL to upload files")
	executeCmd.Flags().StringVar(&subgraphRepoUrl, "url", "https://github.com/u2u-labs/layerg-crawler-v2", "Subgraph repository URL")

	publishCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable caching during deployment")
	executeCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable caching during execution")
	publishCmd.Flags().BoolVar(&devExperimental, "dev", false, "Enable experimental features for development")
	executeCmd.Flags().BoolVar(&devExperimental, "dev", false, "Enable experimental features for development")

	rootCmd.AddCommand(publishCmd, executeCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
