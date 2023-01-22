package main

import (
	"os"
	"fmt"

    "github.com/the-singularity-labs/hydroflow"
    
	"gopkg.in/yaml.v3"
	"github.com/hoenirvili/skapt"
	"github.com/hoenirvili/skapt/argument"
	"github.com/hoenirvili/skapt/flag"
)

func main() {
	app := skapt.Application{
		Name:        "Hydrowflow",
		Description: "Generate a makefile representing your DAG",
		Version:     "1.0.0",
		Handler: func(ctx *skapt.Context) error {
			var err error
			file := ctx.String("file")
			config := ctx.String("config")
			outfile := ctx.String("outfile")
			if file == "" && config == "" {
				return fmt.Errorf("Must provide either a --config or a --file argument")
			}

			if outfile == "" {
				outfile = ".build/Makefile"
			}
			
			configBytes := []byte(config)

			if file != "" {
				configBytes, err = os.ReadFile(file) 
				if err != nil {
					fmt.Errorf("error opening config file: %w", err)
				}
			}

			hydroflowConfig := hydroflow.Hydroflow{}
			if err = yaml.Unmarshal(configBytes, &hydroflowConfig); err != nil {
				return fmt.Errorf("error parsing config JSON: %w", err)
			}

			if err = hydroflowConfig.Validate(); err != nil {
				return fmt.Errorf("error validating config: %w", err)
			}

			makefileString, err := hydroflowConfig.GenerateMakefile()
			if err != nil {
				return fmt.Errorf("error generating Makefile contents: %w", err)
			}

		    f, err := os.Create(outfile)
			if err != nil {
				return fmt.Errorf("error creating Makefile on disk: %w", err)
			}
			defer f.Close()

			_, err = f.WriteString(makefileString)
			if err != nil {
				return fmt.Errorf("error saving Makefile to disk: %w", err)
			}

			return nil
		},
		Flags: flag.Flags{{
			Short: "f", Long: "file",
			Description: "Filepath to config",
			Type:        argument.String,
			Required:	 false,
		}, {
			Short: "c", Long: "config",
			Description: "String of JSON or YAML representing Hydroflow config",
			Type:        argument.String,
			Required:	 true,
		}, {
			Short: "o", Long: "outfile",
			Description: "Filepath where file will be written",
			Type:        argument.String,
			Required:	 true,
		}},
	}
	app.Exec(os.Args)
}
