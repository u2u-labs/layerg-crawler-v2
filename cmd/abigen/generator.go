package abigen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

// Template for generating Go helper functions
const functionTemplate = `
func {{.FunctionName}}(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address{{.InputParams}}) ({{.OutputType}}, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader({{.ABIVar}}))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return {{.ZeroValue}}, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("{{.Name}}"{{.CallParams}})
	if err != nil {
		sugar.Errorf("Failed to pack data for {{.Name}}: %v", err)
		return {{.ZeroValue}}, err
	}

	// Call the contract
	msg := u2u.CallMsg{
		To:   contractAddress,
		Data: data,
	}

	// Execute the call
	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		sugar.Errorf("Failed to call contract: %v", err)
		return {{.ZeroValue}}, err
	}

	// Unpack the result
	{{.UnpackCode}}
	return {{.ReturnValue}}, nil
}
`

type FunctionData struct {
	FunctionName string
	Name         string
	InputParams  string
	OutputType   string
	ABIVar       string
	CallParams   string
	ZeroValue    string
	UnpackCode   string
	ReturnValue  string
}

func GenerateFunction(fn ABIFunction, abiVar string) (string, error) {
	// Convert Solidity types to Go types
	inputParams := ""
	callParams := ""
	for _, input := range fn.Inputs {
		goType := solTypeToGoType(input.Type)
		if input.Name == "" {
			input.Name = "arg"
		}
		inputParams += fmt.Sprintf(", %s %s", input.Name, goType)
		callParams += fmt.Sprintf(", %s", input.Name)
	}

	// Handle output type
	outputType := "interface{}"
	zeroValue := "nil"
	if len(fn.Outputs) == 1 {
		outputType = solTypeToGoType(fn.Outputs[0].Type)
		zeroValue = zeroValueForType(outputType)
	}

	// Create unpack code
	unpackVar := "unpackVar"
	if len(fn.Outputs) == 1 {
		unpackVar = "var " + unpackVar + " " + outputType
	}
	unpackCode := fmt.Sprintf(`%s
	err = contractABI.UnpackIntoInterface(&%s, "%s", result)
	if err != nil {
		sugar.Errorf("Failed to unpack %s: %%v", err)
		return %s, err
	}`, unpackVar, "unpackVar", fn.Name, fn.Name, zeroValue)

	data := FunctionData{
		FunctionName: filepath.Base(abiVar) + strings.Title(fn.Name),
		Name:         fn.Name,
		InputParams:  inputParams,
		OutputType:   outputType,
		ABIVar:       abiVar,
		CallParams:   callParams,
		ZeroValue:    zeroValue,
		UnpackCode:   unpackCode,
		ReturnValue:  "unpackVar",
	}

	tmpl, err := template.New("function").Parse(functionTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func solTypeToGoType(solType string) string {
	switch {
	case strings.HasPrefix(solType, "uint"), strings.HasPrefix(solType, "int"):
		return "*big.Int"
	case solType == "address":
		return "common.Address"
	case solType == "bool":
		return "bool"
	case solType == "string":
		return "string"
	case strings.HasPrefix(solType, "bytes"):
		return "[]byte"
	default:
		return "interface{}"
	}
}

func zeroValueForType(goType string) string {
	switch goType {
	case "*big.Int":
		return "nil"
	case "common.Address":
		return "common.Address{}"
	case "bool":
		return "false"
	case "string":
		return `""`
	case "[]byte":
		return "nil"
	default:
		return "nil"
	}
}
