package hydrowflow

import (
	"fmt"
	"bytes"
	"strings"
	"text/template"
)

const MakefileTemplate string = `
#
# {{.Name}}
# Author: {{.Author}}
#
TARGET_DIR?={{.TargetDir}}

ifeq ($(OS),Windows_NT)
	uname_S := Windows
else
	uname_S := $(shell uname -s)
endif

all: {{.Sinks}}

{{.TargetDir}}:
	mkdir -p {{.TargetDir}}

{{range $rule := .Rules}}
{{.TargetDir}}/{{$rule}}: {{.TargetDir}} {{range $prereq := $rule.Prerequisites}}{{.TargetDir}}/{{$prereq}}{{.end}}
	{{range $command := $rule.Recipe}}
	{{$command}}
	{{end}}
	{{if not $rule.UpdatesTarget}}
	touch {{.TargetDir}}/{{$rule}}
	{{end}}
{{end}}
`

type Hydrowflow struct {
	Name string `json:"name"`
	Author string `json:"author"`
	TargetDir string `json:"target_directory"`
	Rules []Rule  `json:"rules"`
} 

type Rule struct {
	Target string `json:"target"`
	Recipe []string `json:"recipe"`
	UpdatesTarget bool `json:"updates_target"`
	Prerequisites []string `json:"prerequisites"`
}

func (h Hydrowflow) Sinks() string {
	dependencyCounts := map[string]int{}
	for _, rule := range h.Rules {
		if _, exists := dependencyCounts[rule.Target]; !exists {
			dependencyCounts[rule.Target] = 0
		}
		for _, prereq := range rule.Prerequisites {
			if existing, exists := dependencyCounts[prereq]; exists {
				dependencyCounts[prereq] = existing + 1
			} else {
				dependencyCounts[prereq] = 1
			}		
		}

	}

	results := []string{}
	for target, cnt := range dependencyCounts {
		if cnt == 0 {
			results = append(results, target)
		}
	}
	return strings.Join(results, " ")
}

func (h Hydrowflow) Validate() error {
	// TODO:
	// validate targets have no spaces or weird symbols
	// validate commands don't have unescaped newlines
	// validate no cycles (ie its actually a dag)
	return nil
}

func (h Hydrowflow) GenerateMakefile() (string, error) {
    makefileTemplate, err := template.New("makefile_template").Parse(MakefileTemplate)
    if err != nil {
        return "", fmt.Errorf("error parsing template: %w", err)
    }
	var results bytes.Buffer
	makefileTemplate.Execute(&results, h)
	return results.String(), nil
}


