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
		Short: "CLI tool for running subgraph deployments",
		Long:  "LayerG-Deploy is a command-line tool to run subgraph deployments.",
	}
	executeCmd = &cobra.Command{
		Use:   "run",
		Short: "Clone repo and run the crawler",
		Long:  "Clones the repository, and runs the crawler.",
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
	executeCmd.Flags().StringVar(&subgraphRepoUrl, "url", "https://github.com/u2u-labs/layerg-crawler-v2", "Subgraph repository URL")
	executeCmd.Flags().StringVar(&branch, "branch", "master", "Branch to fetch from the repository")

	executeCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable caching during execution")
	executeCmd.Flags().BoolVar(&devExperimental, "dev", false, "Enable experimental features for development")

	rootCmd.AddCommand(executeCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
