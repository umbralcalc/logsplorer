package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/akamensky/argparse"
	"gopkg.in/yaml.v2"
)

// StochadexImplementationStrings is the yaml-loadable config which consists of
// string type names from the stochadex to insert into templating.
type StochadexImplementationStrings struct {
	Iterations           []string `yaml:"iterations"`
	OutputCondition      string   `yaml:"output_condition"`
	OutputFunction       string   `yaml:"output_function"`
	TerminationCondition string   `yaml:"termination_condition"`
	TimestepFunction     string   `yaml:"timestep_function"`
}

// ImplementationStrings is the yaml-loadable config which consists of string type
// names to insert into templating.
type ImplementationStrings struct {
	Streaming             StochadexImplementationStrings `yaml:"streaming"`
	Objectives            []string                       `yaml:"objectives"`
	OptimisationAlgorithm string                         `yaml:"optimisation_algorithm"`
}

// LearnadexArgParse builds the configs parsed as args to the learnadex binary and
// also retrieves other args.
func LearnadexArgParse() (
	string,
	*ImplementationStrings,
) {
	parser := argparse.NewParser("learnadex", "inference and emulation of stochastic phenomena")
	settingsFile := parser.String(
		"s",
		"settings",
		&argparse.Options{Required: true, Help: "yaml config path for settings"},
	)
	implementationsFile := parser.String(
		"i",
		"implementations",
		&argparse.Options{
			Required: true,
			Help:     "yaml config path for string implementations",
		},
	)
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if *settingsFile == "" {
		panic(fmt.Errorf("Parsed no settings config file"))
	}
	if *implementationsFile == "" {
		panic(fmt.Errorf("Parsed no implementations config file"))
	}
	yamlFile, err := ioutil.ReadFile(*implementationsFile)
	if err != nil {
		panic(err)
	}
	var implementations ImplementationStrings
	err = yaml.Unmarshal(yamlFile, &implementations)
	if err != nil {
		panic(err)
	}
	return *settingsFile, &implementations
}

// writeMainProgram writes string representations of various types of data
// to a template tmp/main.go file ready for runtime execution in this main.go
func writeMainProgram() {
	fmt.Println("\nReading in args...")
	settingsFile, implementations := LearnadexArgParse()
	fmt.Println("\nParsed implementations:")
	fmt.Println(implementations)
	iterations := "[]simulator.Iteration{" +
		strings.Join(implementations.Streaming.Iterations, ", ") + "}"
	objectives := "[]learning.LogLikelihood{" +
		strings.Join(implementations.Objectives, ", ") + "}"
	codeTemplate := template.New("learnadexMain")
	template.Must(codeTemplate.Parse(`package main

import (
	"fmt"

	"github.com/umbralcalc/learnadex/pkg/filter"
	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/learnadex/pkg/models"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/optimize"
)

func main() {
	settings := simulator.NewLoadSettingsConfigFromYaml("{{.SettingsFile}}")
	iterations := {{.Iterations}}
	config := &learning.LearnadexConfig{
		Learning: &learning.LearningConfig{
			Streaming:  &simulator.LoadImplementationsConfig{
				Iterations:      iterations,
				OutputCondition: {{.OutputCondition}},
				OutputFunction:  {{.OutputFunction}},
				TerminationCondition: {{.TerminationCondition}},
				TimestepFunction: {{.TimestepFunction}},
			},
			Objectives: {{.Objectives}},
		},
		Optimiser: {{.Algorithm}},
	}
	params := learning.RunFilterParamsLearning(config, settings)
	for i, p := range params {
		fmt.Println("partition ", i)
		for k, v := range p.FloatParams {
			fmt.Println(k, v)
		}
		for k, v := range p.IntParams {
			fmt.Println(k, v)
		}
	}
}`))
	file, err := os.Create("tmp/main.go")
	if err != nil {
		err := os.Mkdir("tmp", 0755)
		if err != nil {
			panic(err)
		}
		file, err = os.Create("tmp/main.go")
	}
	err = codeTemplate.Execute(
		file,
		map[string]string{
			"SettingsFile":         settingsFile,
			"Iterations":           iterations,
			"OutputCondition":      implementations.Streaming.OutputCondition,
			"OutputFunction":       implementations.Streaming.OutputFunction,
			"TerminationCondition": implementations.Streaming.TerminationCondition,
			"TimestepFunction":     implementations.Streaming.TimestepFunction,
			"Objectives":           objectives,
			"Algorithm":            implementations.OptimisationAlgorithm,
		},
	)
	if err != nil {
		panic(err)
	}
	file.Close()
}

func main() {
	// hydrate the template code and write it to tmp/main.go
	writeMainProgram()

	// execute the code
	runCmd := exec.Command("go", "run", "tmp/main.go")
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	if err := runCmd.Run(); err != nil {
		panic(err)
	}
}
