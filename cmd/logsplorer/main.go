package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/umbralcalc/learnadex/pkg/learning"
)

type DataFilter struct {
	AllowedValues   []float64
	ValueRangeUpper float64
	ValueRangeLower float64
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
	if &d.ValueRangeUpper != nil {
		if value > d.ValueRangeUpper {
			ignore = true
			return ignore
		}
	}
	if &d.ValueRangeLower != nil {
		if value < d.ValueRangeLower {
			ignore = true
			return ignore
		}
	}
	return ignore
}

func readLogEntries(
	filename string,
	dataFilterByParam map[string]*DataFilter,
) ([]learning.JsonLogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	logEntries := make([]learning.JsonLogEntry, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var logEntry learning.JsonLogEntry
		line := scanner.Bytes()

		err := json.Unmarshal(line, &logEntry)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}

		include := true
		for param, filter := range dataFilterByParam {
			if param == "partition_index" {
				if filter.Ignore(float64(logEntry.PartitionIndex)) {
					include = false
					break
				}
			} else if param == "objective" {
				if filter.Ignore(logEntry.Objective) {
					include = false
					break
				}
			} else {
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

		logEntries = append(logEntries, logEntry)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return logEntries, nil
}

func main() {
	http.HandleFunc("/api/logsplorer", func(w http.ResponseWriter, r *http.Request) {
		logFilenamesGet := r.URL.Query().Get("filenames")
		filenames := strings.Split(logFilenamesGet, ",")
		allLogEntries := make([]learning.JsonLogEntry, 0)
		dataFilterByParam := make(map[string]*DataFilter)
		for key, values := range r.URL.Query() {
			if key != "filenames" {
				dataFilterByParam[key] = &DataFilter{}
				for _, value := range values {
					if strings.Contains(value, ">") {
						val, err := strconv.ParseFloat(strings.Trim(value, ">"), 64)
						if err != nil {
							http.Error(
								w,
								"Error converting string in query to float64",
								http.StatusInternalServerError,
							)
						}
						dataFilterByParam[key].ValueRangeLower = val
					} else if strings.Contains(value, "<") {
						val, err := strconv.ParseFloat(strings.Trim(value, "<"), 64)
						if err != nil {
							http.Error(
								w,
								"Error converting string in query to float64",
								http.StatusInternalServerError,
							)
						}
						dataFilterByParam[key].ValueRangeUpper = val
					} else {
						dataFilterByParam[key].AllowedValues = make([]float64, 0)
						for _, allowedValue := range strings.Split(value, ",") {
							val, err := strconv.ParseFloat(allowedValue, 64)
							if err != nil {
								http.Error(
									w,
									"Error converting string in query to float64",
									http.StatusInternalServerError,
								)
							}
							dataFilterByParam[key].AllowedValues = append(
								dataFilterByParam[key].AllowedValues,
								val,
							)
						}
					}
				}
			}
		}
		for _, filename := range filenames {
			logEntries, err := readLogEntries(filename, dataFilterByParam)
			if err != nil {
				http.Error(w, "Error reading log entries", http.StatusInternalServerError)
				return
			}
			allLogEntries = append(allLogEntries, logEntries...)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allLogEntries)
	})

	http.ListenAndServe(":8080", nil)
}
