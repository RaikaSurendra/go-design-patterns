// Package main demonstrates the Facade pattern.
//
// Facade provides a simplified interface to a complex subsystem,
// wrapping multiple components behind a single easy-to-use API.
package main

import "fmt"

// --- Subsystem components ---

type CPU struct{}

func (c *CPU) Freeze()              { fmt.Println("CPU: freeze") }
func (c *CPU) Jump(position int64)  { fmt.Printf("CPU: jump to 0x%x\n", position) }
func (c *CPU) Execute()             { fmt.Println("CPU: executing") }

type Memory struct{}

func (m *Memory) Load(position int64, data []byte) {
	fmt.Printf("Memory: loading %d bytes at 0x%x\n", len(data), position)
}

type HardDrive struct{}

func (h *HardDrive) Read(lba int64, size int) []byte {
	fmt.Printf("HardDrive: reading %d bytes from sector %d\n", size, lba)
	return make([]byte, size)
}

// --- Facade ---

const bootAddress int64 = 0x7C00
const bootSector int64 = 0
const sectorSize int = 512

type Computer struct {
	cpu       CPU
	memory    Memory
	hardDrive HardDrive
}

func NewComputer() *Computer {
	return &Computer{}
}

// Start hides the complex boot sequence behind a single method.
func (c *Computer) Start() {
	c.cpu.Freeze()
	c.memory.Load(bootAddress, c.hardDrive.Read(bootSector, sectorSize))
	c.cpu.Jump(bootAddress)
	c.cpu.Execute()
}

func main() {
	// The client interacts with one simple method instead of
	// coordinating CPU, Memory, and HardDrive directly.
	pc := NewComputer()
	pc.Start()
}
