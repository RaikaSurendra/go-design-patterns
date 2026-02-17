// Package main demonstrates the Template pattern.
//
// Template defines the skeleton of an algorithm, deferring some steps
// to implementations. In Go we use interfaces + a template function.
package main

import (
	"fmt"
	"strings"
)

// Record represents a parsed data entry.
type Record struct {
	Fields map[string]string
}

// DataMiner defines the steps that vary between implementations.
type DataMiner interface {
	Extract() ([]string, error)
	Parse(raw []string) []Record
}

// Mine is the template method that defines the fixed algorithm skeleton.
func Mine(miner DataMiner) ([]Record, error) {
	raw, err := miner.Extract()
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}
	records := miner.Parse(raw)
	return records, nil
}

// --- CSV implementation ---

type CSVMiner struct {
	Data string
}

func (c *CSVMiner) Extract() ([]string, error) {
	return strings.Split(c.Data, "\n"), nil
}

func (c *CSVMiner) Parse(raw []string) []Record {
	var records []Record
	for _, line := range raw {
		if line == "" {
			continue
		}
		fields := strings.Split(line, ",")
		r := Record{Fields: map[string]string{"raw": strings.Join(fields, " | ")}}
		records = append(records, r)
	}
	return records
}

// --- JSON-lines implementation ---

type JSONLinesMiner struct {
	Data string
}

func (j *JSONLinesMiner) Extract() ([]string, error) {
	return strings.Split(j.Data, "\n"), nil
}

func (j *JSONLinesMiner) Parse(raw []string) []Record {
	var records []Record
	for _, line := range raw {
		if line == "" {
			continue
		}
		r := Record{Fields: map[string]string{"json": line}}
		records = append(records, r)
	}
	return records
}

func main() {
	csvData := "name,age,city\nalice,30,NYC\nbob,25,SF"
	csvMiner := &CSVMiner{Data: csvData}

	fmt.Println("--- CSV Mining ---")
	records, _ := Mine(csvMiner)
	for _, r := range records {
		fmt.Println(r.Fields["raw"])
	}

	jsonData := `{"name":"alice","age":30}` + "\n" + `{"name":"bob","age":25}`
	jsonMiner := &JSONLinesMiner{Data: jsonData}

	fmt.Println("\n--- JSON Lines Mining ---")
	records, _ = Mine(jsonMiner)
	for _, r := range records {
		fmt.Println(r.Fields["json"])
	}
}
