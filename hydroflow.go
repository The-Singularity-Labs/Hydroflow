package hydroflow

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
$(TARGET_DIR)/{{$rule.Target}}: $(TARGET_DIR) {{range $prereq := $rule.Prerequisites}}$(TARGET_DIR)/{{$prereq}} {{end}} {{range $command := $rule.EscapedRecipe}}
	{{$command}}{{end}}{{if not $rule.UpdatesTarget}}
	touch {{$.TargetDir}}/{{$rule.Target}}
	{{end}}
{{end}}
`

type Hydroflow struct {
	Name string `yaml:"name"`
	Author string `yaml:"author"`
	TargetDir string `yaml:"target_directory"`
	Rules []Rule  `yaml:"rules"`
} 

type Rule struct {
	Target string `yaml:"target"`
	Recipe []string `yaml:"recipe"`
	UpdatesTarget bool `yaml:"updates_target"`
	Prerequisites []string `yaml:"prerequisites"`
}

func (r Rule) EscapedRecipe() []string {
	results := []string{}
	for _, command := range r.Recipe {
		// make needs $ escaped with $$
		// https://til.hashrocket.com/posts/k3kjqxtppx-escape-dollar-sign-on-makefiles
		results = append(results, strings.ReplaceAll(command, "$", "$$"))
	}
	return results
}

func (h Hydroflow) Sinks() string {
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
	fmt.Println(dependencyCounts)

	results := []string{}
	for target, cnt := range dependencyCounts {
		if cnt == 0 {
			results = append(results, fmt.Sprintf("$(TARGET_DIR)/%s", target))
		}
	}
	return strings.Join(results, " ")
}

func (h Hydroflow) Validate() error {
	// TODO:
	// validate targets have no spaces or weird symbols
	// validate commands don't have unescaped newlines
	// validate no cycles (ie its actually a dag)
	return nil
}

func (h Hydroflow) GenerateMakefile() (string, error) {
    makefileTemplate, err := template.New("makefile_template").Parse(MakefileTemplate)
    if err != nil {
        return "", fmt.Errorf("error parsing template: %w", err)
    }
	var results bytes.Buffer
	err = makefileTemplate.Execute(&results, h)
	if err != nil {
		return "", fmt.Errorf("Error executing template %w: ", err)
	}
	return results.String(), nil
}


