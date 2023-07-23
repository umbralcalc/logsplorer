package main

import (
	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func main() {
	settingsFile := ""
	settings := simulator.NewLoadSettingsConfigFromYaml(settingsFile)
	config := &learning.LearningConfig{
		Streaming:  streamingConfigs,
		Objectives: objectives,
	}
	learning.RunFilterParamsLearning(config, settings)
}
