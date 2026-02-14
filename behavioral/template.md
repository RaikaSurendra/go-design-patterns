# Template Pattern

The template pattern defines the skeleton of an algorithm in a base operation,
deferring some steps to subclasses. It lets subclasses redefine certain steps of
an algorithm without changing the algorithm's structure.

In Go, since there is no classical inheritance, we use an interface for the
varying steps and a function for the fixed skeleton.

## Implementation

```go
package template

import "fmt"

// DataMiner defines the steps that vary between implementations.
type DataMiner interface {
	Open(path string) error
	Extract() ([]string, error)
	Parse(raw []string) []Record
	Close()
}

// Record represents a parsed data entry.
type Record struct {
	Fields map[string]string
}

// Mine is the template method. It defines the fixed algorithm skeleton and
// delegates the varying steps to the DataMiner interface.
func Mine(miner DataMiner, path string) ([]Record, error) {
	if err := miner.Open(path); err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer miner.Close()

	raw, err := miner.Extract()
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	records := miner.Parse(raw)
	return records, nil
}
```

A concrete implementation for CSV files:

```go
package template

import (
	"encoding/csv"
	"os"
	"strings"
)

type CSVMiner struct {
	file *os.File
}

func (c *CSVMiner) Open(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	c.file = f
	return nil
}

func (c *CSVMiner) Extract() ([]string, error) {
	reader := csv.NewReader(c.file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, row := range records {
		lines = append(lines, strings.Join(row, ","))
	}
	return lines, nil
}

func (c *CSVMiner) Parse(raw []string) []Record {
	var records []Record
	for _, line := range raw {
		fields := strings.Split(line, ",")
		r := Record{Fields: map[string]string{"raw": strings.Join(fields, " | ")}}
		records = append(records, r)
	}
	return records
}

func (c *CSVMiner) Close() {
	if c.file != nil {
		c.file.Close()
	}
}
```

## Usage

```go
miner := &template.CSVMiner{}

records, err := template.Mine(miner, "data.csv")
if err != nil {
	log.Fatal(err)
}

for _, r := range records {
	fmt.Println(r.Fields["raw"])
}
```

## Rules of Thumb

- Use template when you have several classes that contain nearly identical algorithms with some minor differences.
- In Go, the template function accepts an interface rather than relying on inheritance — this is idiomatic composition over inheritance.
- Hook methods (optional steps with default behavior) can be expressed as interfaces with default implementations using embedding.
