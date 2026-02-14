# Facade Pattern <span class="gp-difficulty gp-difficulty--easy">Easy</span>

The facade pattern provides a simplified interface to a complex subsystem. It
wraps multiple components behind a single, easy-to-use API so that clients don't
need to interact with each subsystem directly.

## Implementation

```go
package computer

import "fmt"

// --- Subsystem components ---

type CPU struct{}

func (c *CPU) Freeze() { fmt.Println("CPU: freeze") }
func (c *CPU) Jump(position int64) { fmt.Printf("CPU: jump to 0x%x\n", position) }
func (c *CPU) Execute() { fmt.Println("CPU: executing") }

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

func New() *Computer {
	return &Computer{}
}

// Start hides the complex boot sequence behind a single method.
func (c *Computer) Start() {
	c.cpu.Freeze()
	c.memory.Load(bootAddress, c.hardDrive.Read(bootSector, sectorSize))
	c.cpu.Jump(bootAddress)
	c.cpu.Execute()
}
```

## Usage

```go
// The client interacts with one simple method instead of
// coordinating CPU, Memory, and HardDrive directly.
pc := computer.New()
pc.Start()
// CPU: freeze
// HardDrive: reading 512 bytes from sector 0
// Memory: loading 512 bytes at 0x7c00
// CPU: jump to 0x7c00
// CPU: executing
```

## Rules of Thumb

- Facade does not prevent clients from accessing subsystems directly if they need to — it simply provides a convenient default path.
- Use facade when you want to layer a subsystem: the facade defines the entry point for each level.
- Facade and mediator are similar in that they abstract existing classes. Facade defines a simpler interface, while mediator introduces new behavior by coordinating between components.
