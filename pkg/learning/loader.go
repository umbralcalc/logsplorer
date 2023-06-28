package learning

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/akamensky/argparse"
	"gopkg.in/yaml.v2"
)

// ExtraLoadSettingsConfig is the yaml-loadable config extends the settings available
// the stochasdex simulator.LoadSettingsConfig to include settings that are only
// necessary in the learnadex package.
type ExtraLoadSettingsConfig struct {
	BurnInSteps int `yaml:"burn_in_steps"`
}

// NewExtraLoadSettingsConfigFromYaml creates a new ExtraLoadSettingsConfig from
// a provided yaml path.
func NewExtraLoadSettingsConfigFromYaml(path string) *ExtraLoadSettingsConfig {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var settings ExtraLoadSettingsConfig
	err = yaml.Unmarshal(yamlFile, &settings)
	if err != nil {
		panic(err)
	}
	return &settings
}

// NewExtraLoadSettingsConfigFromYaml calls NewExtraLoadSettingsConfigFromYaml
// with a yaml path string provided by argparse.
func NewExtraLoadSettingsConfigFromArgParsedYaml() *ExtraLoadSettingsConfig {
	parser := argparse.NewParser(
		"learnadex",
		"learn a variety of stochastic phenomena directly from input data",
	)
	s := parser.String(
		"s",
		"string",
		&argparse.Options{Required: true, Help: "yaml config path for settings"},
	)
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	return NewExtraLoadSettingsConfigFromYaml(*s)
}
