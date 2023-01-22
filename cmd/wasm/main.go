package main

import (
	"fmt"
	"syscall/js"
	
	"github.com/the-singularity-labs/hydroflow"

	"gopkg.in/yaml.v3"
)



func main() {
	js.Global().Set("GenerateMakefile", hydrowflowWrapper())
	<-make(chan bool)
}

func hydrowflowWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return wrap("", fmt.Errorf("Not enough arguments"))
		}
		configBytes := []byte(args[0].String())

		hydroflowConfig := hydroflow.Hydroflow{}
		if err := yaml.Unmarshal(configBytes, &hydroflowConfig); err != nil {
			return wrap("", fmt.Errorf("error parsing config JSON: %w", err))
		}

		if err := hydroflowConfig.Validate(); err != nil {
			return wrap("", fmt.Errorf("error validating config: %w", err))
		}

		makefileString, err := hydroflowConfig.GenerateMakefile()
		if err != nil {
			return wrap("", fmt.Errorf("error generating Makefile contents: %w", err))
		}

		return wrap(makefileString, nil)
	})
}

func wrap(result string, err error) map[string]interface{} {
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	return map[string]interface{}{
		"error":   errString,
		"data": result,
	}
}