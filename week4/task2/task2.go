package main

import "fmt"

type Person struct {
	Name  string
	Email string
}

func (p *Person) Contact() string {
	return "Contacts of " + p.Name + ": " + p.Email
}

type Student struct {
	Person
	GPA float64
}

func (s Student) Contact() string {
	return "Contacts of " + s.Name + ": " + s.Email + ", GPA: " + fmt.Sprintf("%.2f", s.GPA)
}

type Teacher struct {
	Person
	Department string
}

func (t *Teacher) Contact() string {
	return "Contacts of " + t.Name + ": " + t.Email + ", Department: " + t.Department
}

func main() {

	s := Student{
		Person: Person{
			Name:  "Alice",
			Email: "alice@example.com",
		},
		GPA: 3.8,
	}

	t := Teacher{
		Person: Person{
			Name:  "Bob",
			Email: "bob@example.com",
		},
		Department: "Computer Science",
	}

	fmt.Println(s.Contact())
	fmt.Println(t.Contact())
}
