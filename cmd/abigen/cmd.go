package abigen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var GeneratorCmd = &cobra.Command{
	Use:   "abigen",
	Short: "Generate Go helper functions from ABI files",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")

		// Generate abi-var from input filename
		abiName := strings.TrimSuffix(filepath.Base(input), ".json")
		abiName = strings.ToUpper(abiName)

		functions, err := ParseABIFile(input)

		if err != nil {
			return fmt.Errorf("failed to parse ABI file: %v", err)
		}

		// Generate ABI string for view functions
		abiString, err := FilterViewFunctions(functions)
		if err != nil {
			return fmt.Errorf("failed to generate ABI string: %v", err)
		}

		var code strings.Builder
		code.WriteString("package " + filepath.Base(filepath.Dir(output)) + "\n\n")
		code.WriteString(`
// Code generated - DO NOT EDIT.
import (
	"context"
	"math/big"
	"strings"

	u2u "github.com/unicornultrafoundation/go-u2u"
	"github.com/unicornultrafoundation/go-u2u/accounts/abi"
	"github.com/unicornultrafoundation/go-u2u/common"
	"github.com/unicornultrafoundation/go-u2u/ethclient"
	"go.uber.org/zap"
)

`)

		// Write ABI string constant
		code.WriteString(fmt.Sprintf("const %s = `%s`\n\n", abiName, abiString))

		// Generate helper functions for view functions only
		for _, fn := range functions {
			if fn.Type != "function" || fn.StateMutability != "view" {
				continue
			}
			fnCode, err := GenerateFunction(fn, abiName)
			if err != nil {
				return fmt.Errorf("failed to generate function %s: %v", fn.Name, err)
			}
			code.WriteString(fnCode)
			code.WriteString("\n")
		}

		if err := os.WriteFile(output, []byte(code.String()), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %v", err)
		}

		fmt.Printf("Successfully generated helper functions and ABI string in %s\n", output)
		return nil
	},
}

func init() {
	GeneratorCmd.Flags().StringP("input", "i", "", "Input ABI file")
	GeneratorCmd.Flags().StringP("output", "o", "", "Output Go file")
	GeneratorCmd.MarkFlagRequired("input")
	GeneratorCmd.MarkFlagRequired("output")
}
