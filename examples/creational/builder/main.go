// Package main demonstrates the Builder pattern.
//
// Builder separates the construction of a complex object from its
// representation so that the same construction process can create
// different configurations.
package main

import "fmt"

type Speed float64

const (
	MPH Speed = 1
	KPH Speed = 1.60934
)

type Color string

const (
	Blue  Color = "blue"
	Green Color = "green"
	Red   Color = "red"
)

type Wheels string

const (
	SportsWheels Wheels = "sports"
	SteelWheels  Wheels = "steel"
)

type Car struct {
	color    Color
	wheels   Wheels
	topSpeed Speed
}

func (c *Car) Drive() {
	fmt.Printf("Driving a %s car with %s wheels (top speed: %.0f MPH)\n",
		c.color, c.wheels, c.topSpeed/MPH)
}

type CarBuilder struct {
	color    Color
	wheels   Wheels
	topSpeed Speed
}

func NewBuilder() *CarBuilder {
	return &CarBuilder{
		color:    Blue,
		wheels:   SteelWheels,
		topSpeed: 100 * MPH,
	}
}

func (b *CarBuilder) Paint(c Color) *CarBuilder {
	b.color = c
	return b
}

func (b *CarBuilder) SetWheels(w Wheels) *CarBuilder {
	b.wheels = w
	return b
}

func (b *CarBuilder) TopSpeed(s Speed) *CarBuilder {
	b.topSpeed = s
	return b
}

func (b *CarBuilder) Build() *Car {
	return &Car{
		color:    b.color,
		wheels:   b.wheels,
		topSpeed: b.topSpeed,
	}
}

func main() {
	familyCar := NewBuilder().
		Paint(Red).
		SetWheels(SteelWheels).
		TopSpeed(50 * MPH).
		Build()

	sportsCar := NewBuilder().
		Paint(Green).
		SetWheels(SportsWheels).
		TopSpeed(150 * MPH).
		Build()

	fmt.Println("Family car:")
	familyCar.Drive()

	fmt.Println("\nSports car:")
	sportsCar.Drive()
}
