package main

import "fmt"

type Student struct {
	Name    string
	Age     int
	Courses []string
}

func NewStudent(name string, age int, courses ...string) *Student {
	return &Student{Name: name, Age: age, Courses: courses}
}

func main() {

	var s1 Student
	s2 := Student{Name: "Alice", Age: 20}
	s3 := Student{"Bob", 22, []string{"Math", "Science"}}
	s4 := &Student{Name: "Charlie"}
	s5 := NewStudent("David", 21)

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)
	fmt.Println(*s4)
	fmt.Println(s5)

	addRealCourse(s4, "History")
	fmt.Println("Appended course to s4:", *s4)

	addRealCourse(s5, "Art")
	fmt.Println("Appended course to s5:", *s5)

}

func addCourse(s Student, course string) {
	s.Courses = append(s.Courses, course)
}

func addRealCourse(s *Student, course string) {
	s.Courses = append(s.Courses, course)
}
