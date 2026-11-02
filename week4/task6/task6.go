package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r *Rectangle) Scale(k float64) {
	r.Width *= k
	r.Height *= k
}

type Circle struct {
	Radius float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func main() {
	// 4. Put Rectangle value and Circle pointer into []Shape
	rect := Rectangle{Width: 4, Height: 5}
	circ := &Circle{Radius: 3}

	shapes := []Shape{
		rect, // Rectangle satisfies Shape (value receiver methods)
		circ, // *Circle satisfies Shape (pointer receiver methods)
	}

	for i, s := range shapes {
		fmt.Printf("Shape %d - Area: %.2f, Perimeter: %.2f\n", i+1, s.Area(), s.Perimeter())
	}
}
