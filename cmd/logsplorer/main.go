package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/umbralcalc/learnadex/pkg/learning"
	"gopkg.in/yaml.v2"
)

// LogsplorerConfig defines the settings which can be used to configure the
// logs exploration and visualisation app.
type LogsplorerConfig struct {
	Address string `yaml:"address"`
	Handle  string `yaml:"handle"`
}

// QueryLogEntry is the output format from the logsplorer api.
type QueryLogEntry struct {
	PartitionIterations int         `json:"partition_iterations"`
	Entry               interface{} `json:"entry"`
}

// ValueLimit just represents a limit on the range of a given log value.
type ValueLimit struct {
	Upper bool
	Limit float64
}

// DataFilter is the struct containing the filtering logic to apply to the logs.
type DataFilter struct {
	AllowedValues []float64
	ValueLimits   []ValueLimit
}

func (d *DataFilter) Ignore(value float64) bool {
	ignore := false
	if d.AllowedValues != nil {
		ignore = true
		for _, allowed := range d.AllowedValues {
			if value == allowed {
				ignore = false
			}
		}
		if ignore {
			return ignore
		}
	}
	if d.ValueLimits != nil {
		for _, limit := range d.ValueLimits {
			if limit.Upper && value >= limit.Limit ||
				!limit.Upper && value <= limit.Limit {
				ignore = true
				return ignore
			}
		}
	}
	return ignore
}

func (d *DataFilter) SetValue(value string) error {
	if strings.Contains(value, ">") {
		if d.ValueLimits == nil {
			d.ValueLimits = make([]ValueLimit, 0)
		}
		val, err := strconv.ParseFloat(strings.Split(value, ">")[1], 64)
		if err != nil {
			return err
		}
		d.ValueLimits = append(d.ValueLimits, ValueLimit{Upper: false, Limit: val})
	} else if strings.Contains(value, "<") {
		if d.ValueLimits == nil {
			d.ValueLimits = make([]ValueLimit, 0)
		}
		val, err := strconv.ParseFloat(strings.Split(value, "<")[1], 64)
		if err != nil {
			return err
		}
		d.ValueLimits = append(d.ValueLimits, ValueLimit{Upper: true, Limit: val})
	} else {
		d.AllowedValues = make([]float64, 0)
		for _, allowedValue := range strings.Split(value, ",") {
			val, err := strconv.ParseFloat(allowedValue, 64)
			if err != nil {
				return err
			}
			d.AllowedValues = append(
				d.AllowedValues,
				val,
			)
		}
	}
	return nil
}

// readLogEntries reads a file while apply the filtering logic to its data line-by-line
// and then returns the corresponding log entry structs which pass through the filter.
func readLogEntries(
	filename string,
	dataFilterByParam map[string]*DataFilter,
) ([]QueryLogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	partitionIterations := make(map[int]int)
	queryLogEntries := make([]QueryLogEntry, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var logEntry learning.JsonLogEntry
		line := scanner.Bytes()

		err := json.Unmarshal(line, &logEntry)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}

		// keep a track of how many iterations each partition has been through
		_, ok := partitionIterations[logEntry.PartitionIndex]
		if !ok {
			partitionIterations[logEntry.PartitionIndex] = 0
		}
		partitionIterations[logEntry.PartitionIndex] += 1

		include := true
		for param, filter := range dataFilterByParam {
			switch param {
			case "partition_iterations":
				if filter.Ignore(float64(partitionIterations[logEntry.PartitionIndex])) {
					include = false
					break
				}
			case "partition_index":
				if filter.Ignore(float64(logEntry.PartitionIndex)) {
					include = false
					break
				}
			case "objective":
				if filter.Ignore(logEntry.Objective) {
					include = false
					break
				}
			default:
				_, ok := logEntry.FloatParams[param]
				if !ok {
					_, ok = logEntry.IntParams[param]
					if !ok {
						fmt.Println("API Error: param not available:", param)
						return nil, nil
					}
					for _, value := range logEntry.IntParams[param] {
						if filter.Ignore(float64(value)) {
							include = false
							break
						}
					}
				} else {
					for _, value := range logEntry.FloatParams[param] {
						if filter.Ignore(value) {
							include = false
							break
						}
					}
				}
			}
			if !include {
				break
			}
		}
		if !include {
			continue
		}

		queryLogEntries = append(
			queryLogEntries,
			QueryLogEntry{
				PartitionIterations: partitionIterations[logEntry.PartitionIndex],
				Entry:               logEntry,
			},
		)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return queryLogEntries, nil
}

// reorderKeyValuesSymbols is a small adjustment needed to key and values so that other
// symbols can be used in api query strings.
func reorderKeyValuesSymbols(key string, values []string) (string, []string) {
	if strings.Contains(key, ">") {
		values = []string{key}
		key = strings.Split(key, ">")[0]
	}
	if strings.Contains(key, "<") {
		values = []string{key}
		key = strings.Split(key, "<")[0]
	}
	return key, values
}

// LogsplorerArgParse builds the configs parsed as args to the logsplorer binary and
// also retrieves other args.
func LogsplorerArgParse() *LogsplorerConfig {
	parser := argparse.NewParser("logsplorer", "visualisation and exploration of learnadex logs")
	configFile := parser.String(
		"c",
		"config",
		&argparse.Options{Required: true, Help: "yaml config path"},
	)
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	if *configFile == "" {
		panic(fmt.Errorf("Parsed no config file"))
	}
	yamlFile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	var config LogsplorerConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return &config
}

func main() {
	config := LogsplorerArgParse()
	http.HandleFunc(config.Handle, func(w http.ResponseWriter, r *http.Request) {
		logFilenamesGet := r.URL.Query().Get("filenames")
		filenames := strings.Split(logFilenamesGet, ",")
		allQueryLogEntries := make([]QueryLogEntry, 0)
		dataFilterByParam := make(map[string]*DataFilter)
		for key, values := range r.URL.Query() {
			if key == "filenames" {
				continue
			}
			key, values = reorderKeyValuesSymbols(key, values)
			dataFilterByParam[key] = &DataFilter{}
			for _, value := range values {
				err := dataFilterByParam[key].SetValue(value)
				if err != nil {
					http.Error(
						w,
						"Error converting string in query to float64",
						http.StatusInternalServerError,
					)
				}
			}
		}
		for _, filename := range filenames {
			queryLogEntries, err := readLogEntries(filename, dataFilterByParam)
			if err != nil {
				http.Error(w, "Error reading log entries", http.StatusInternalServerError)
				return
			}
			allQueryLogEntries = append(allQueryLogEntries, queryLogEntries...)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allQueryLogEntries)
	})

	http.ListenAndServe(config.Address, nil)
}
