package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/akamensky/argparse"
	"gopkg.in/yaml.v2"
)

// ImplementationStrings is the yaml-loadable config which consists of string type
// names to insert into templating.
type ImplementationStrings struct {
	DataStreamers         []string `yaml:"data_streamers"`
	Objectives            []string `yaml:"objectives"`
	OptimisationAlgorithm string   `yaml:"optimisation_algorithm"`
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
	streamingConfigs := "[]*learning.DataStreamingConfig{" +
		strings.Join(implementations.DataStreamers, ", ") + "}"
	objectives := "[]learning.LogLikelihood{" +
		strings.Join(implementations.Objectives, ", ") + "}"
	codeTemplate := template.New("learnadexMain")
	template.Must(codeTemplate.Parse(`package main

func main() {
	settings := simulator.NewLoadSettingsConfigFromYaml({{.SettingsFile}})
	yamlFile, err = ioutil.ReadFile({{.SettingsFile}})
	if err != nil {
		panic(err)
	}
	var extraSettings learning.ExtraLoadSettings
	err = yaml.Unmarshal(yamlFile, &extraSettings)
	if err != nil {
		panic(err)
	}
	config := &learning.LearnadexConfig{
		Learning: &learning.LearningConfig{
			Streaming:  {{.StreamingConfigs}},
			Objectives: {{.Objectives}},
		},
		Optimiser: &learning.OptimiserConfig{
			Algorithm: {{.Algorithm}},
			ParamsToOptimise: extraSettings.ParamsToOptimise,
		},
	}
	learning.RunFilterParamsLearning(config, settings)
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
			"SettingsFile":     settingsFile,
			"StreamingConfigs": streamingConfigs,
			"Objectives":       objectives,
			"Algorithm":        implementations.OptimisationAlgorithm,
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
}
