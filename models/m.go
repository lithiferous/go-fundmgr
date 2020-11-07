package models

// human
type Payer struct {
	Person *Person
	Name   string
}

// client
type Person struct {
	Name  string
	Value float64
}

// transaction
type Delta struct {
	Name  string
	Value float64
}
